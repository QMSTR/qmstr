package database

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"encoding/json"

	"github.com/QMSTR/qmstr/pkg/service"
	client "github.com/dgraph-io/dgo"
	"github.com/dgraph-io/dgo/protos/api"

	"google.golang.org/grpc"
)

const schema = `
path: string @index(trigram) .
hash: string @index(exact) .
type: string @index(hash) .
name: string @index(hash) .
nodeType: int @index(int) .
dataNodes: uid @reverse .
data: string @index(hash) .
`

type DataBase struct {
	client      *client.Dgraph
	insertQueue chan *service.FileNode
	insertMutex *sync.Mutex
	pending     uint64
}

// Setup connects to dgraph and returns the instance
func Setup(dbAddr string, queueWorkers int) (*DataBase, error) {
	log.Println("Setting up database connection")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	conn, err := grpc.DialContext(ctx, dbAddr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithDefaultCallOptions(grpc.MaxCallSendMsgSize(100<<20),
			grpc.MaxCallRecvMsgSize(100<<20)))
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Failed to connect to the dgraph server: %v", err)
	}

	db := &DataBase{
		client:      client.NewDgraphClient(api.NewDgraphClient(conn)),
		insertQueue: make(chan *service.FileNode, 1000),
		insertMutex: &sync.Mutex{},
	}

	for {
		err = db.client.Alter(ctx, &api.Operation{
			Schema: schema,
		})
		if err == nil {
			cancel()
			break
		}
	}

	for i := 0; i < queueWorkers; i++ {
		go queueWorker(db)
	}

	return db, nil
}

func (db *DataBase) AwaitBuildComplete() {
	// TODO replace busy waiting with proper signaling
	log.Println("Waiting for inserts")
	for {
		pendingInserts := atomic.LoadUint64(&db.pending)
		if pendingInserts == 0 {
			break
		}
		log.Printf("Pending inserts %d", pendingInserts)
		time.Sleep(2 * time.Second)
	}
	close(db.insertQueue)
}

// AddFileNode adds a node to the insert queue
func (db *DataBase) AddFileNode(node *service.FileNode) {
	atomic.AddUint64(&db.pending, 1)
	for _, dep := range node.DerivedFrom {
		db.AddFileNode(dep)
	}
	db.insertQueue <- node
}

// the queueWorker runs in a go routine and inserts the nodes from the insert queue into the database
func queueWorker(db *DataBase) {
	for {
		node := <-db.insertQueue
		if node == nil {
			return
		}
		ready := true
		for idx, dep := range node.DerivedFrom {
			if dep.Uid == "" {
				// missing dep
				ready = false
				// look up dep in db
				uid, err := db.GetFileNodeUid(dep.Hash)
				if err != nil {
					panic(err)
				}
				// found uid
				if uid != "" {
					node.DerivedFrom[idx].Uid = uid
				}
			}
		}

		if !ready {
			// put node back to queue
			go func() { db.insertQueue <- node }()
			continue
		}

		// we are ready to insert the node
		db.insertMutex.Lock()
		uid, err := db.GetFileNodeUid(node.Hash)
		if err != nil {
			panic(err)
		}
		if uid != "" {
			node.Uid = uid
		}
		service.SanitizeFileNode(node)
		uid, err = dbInsert(db.client, node)
		if err != nil {
			panic(err)
		}
		atomic.AddUint64(&db.pending, ^uint64(0))
		db.insertMutex.Unlock()
	}
}

func (db *DataBase) AlterFileNode(node *service.FileNode) (string, error) {
	db.insertMutex.Lock()
	service.SanitizeFileNode(node)
	uid, err := dbInsert(db.client, node)
	db.insertMutex.Unlock()
	return uid, err
}

func (db *DataBase) AlterPackageNode(pkgNode *service.PackageNode) (string, error) {
	db.insertMutex.Lock()
	// Get the package uid from db and pass it to the altered package node
	pkg, err := db.GetPackageNode(pkgNode.Session)
	service.SanitizePackageNode(pkgNode, pkg)
	uid, err := dbInsert(db.client, pkgNode)
	db.insertMutex.Unlock()
	return uid, err
}

func (db *DataBase) AddInfoNodes(nodeID string, infonodes ...*service.InfoNode) error {
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()

	const q = `
	query Node($id: string){
		node(func: uid($id)) @recurse(loop: false) {
			uid
			nodeType
			additionalInfo
		}
	}
	`
	vars := map[string]string{"$id": nodeID}
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return err
	}

	type GenericNode struct {
		Uid            string
		NodeType       int32
		AdditionalInfo []*service.InfoNode
	}

	var receiverNode GenericNode
	err = json.Unmarshal(resp.Json, &receiverNode)
	if err != nil {
		log.Fatal(err)
	}

	if receiverNode.NodeType != service.NodeTypeFileNode && receiverNode.NodeType != service.NodeTypePackageNode {
		return errors.New("can not attach infonode, receiver is neither file nor package node")
	}

	receiverNode.AdditionalInfo = append(receiverNode.AdditionalInfo, infonodes...)

	_, err = dbInsert(db.client, receiverNode)
	if err != nil {
		return err
	}

	return nil
}

// GetFileNodeUid returns the UID of the node if exists otherwise ""
func (db *DataBase) GetFileNodeUid(hash string) (string, error) {

	ret := map[string][]*service.FileNode{}

	q := `query Node($Hash: string){
		  hasNode(func: eq(hash, $Hash)) {
			uid
		  }}`

	vars := map[string]string{"$Hash": hash}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return "", err
	}

	// no node with such hash
	if len(ret["hasNode"]) == 0 {
		return "", nil
	}
	return ret["hasNode"][0].Uid, nil
}

func (db *DataBase) AddPackageNode(pNode *service.PackageNode) (string, error) {
	pNode.NodeType = service.NodeTypePackageNode
	return dbInsert(db.client, pNode)
}

// the data should be JSON marshalable
func dbInsert(c *client.Dgraph, data interface{}) (string, error) {
	txn := c.NewTxn()
	defer txn.Discard(context.Background())

	jNode, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	ret, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: jNode})
	if err != nil {
		return "", err
	}

	uid := ret.Uids["blank-0"]
	return uid, nil
}

func getVarName(index int) string {
	var result string
	for index > 25 {
		result = result + string(rune(65))
		index = index - 26
	}
	result = result + string(rune(65+index))
	return result
}

func (db *DataBase) GetInfoNodeByDataNode(infonodetype string, datanodes ...*service.InfoNode_DataNode) (*service.InfoNode, error) {

	var retInfoNode *service.InfoNode

	runeDataNodeMap := map[string]*service.InfoNode_DataNode{}

	for idx, datanode := range datanodes {
		service.SanitizeDataNode(datanodes[idx])
		runeDataNodeMap[getVarName(idx)] = datanode
	}

	ret := map[string][]*service.InfoNode{}

	q := `query InfoNodeByDataNode($InfoType: string) {
				{{range $var, $data := .}}
				var(func:eq(nodeType, 3)) @filter(eq(type, "{{$data.Type}}") AND eq(data, "{{$data.Data}}")) {
					{{$var}} as ~dataNodes
				}
				{{end}}
		
				getInfoByData(func:eq(nodeType, 2)) @filter(eq(type, $InfoType) {{range $var, $data := .}} AND uid({{$var}}) {{end}}) {
					uid
					type
					dataNodes {
						type
						data
					}
				}
	  		}`

	queryTmpl, err := template.New("infobydata").Parse(q)

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, runeDataNodeMap)
	if err != nil {
		panic(err)
	}

	queryString := b.String()

	vars := map[string]string{"$InfoType": infonodetype}

	err = db.queryInfoNodes(queryString, vars, &ret)
	if err != nil {
		return nil, err
	}

	if len(ret["getInfoByData"]) > 0 {
		retInfoNode = ret["getInfoByData"][0]
	}

	if retInfoNode == nil {
		infoNode := &service.InfoNode{Type: infonodetype, NodeType: service.NodeTypeInfoNode}
		infoNode.DataNodes = datanodes
		uid, err := dbInsert(db.client, infoNode)
		if err != nil {
			return nil, err
		}
		infoNode.Uid = uid
		retInfoNode = infoNode
	}

	return retInfoNode, nil
}

func (db *DataBase) GetPackageNode(session string) (*service.PackageNode, error) {
	ret := map[string][]*service.PackageNode{}

	q := `query PackageNode($Session: string) {
		getPackageNode(func: eq(nodeType, 5)) @recurse(loop: false) {
			uid
			session
			nodeType
			hash
			name
			type
			targets
			derivedFrom
			path
			additionalInfo
			dataNodes
			data
		  }}`

	vars := map[string]string{"$Session": session}

	err := db.queryPackage(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	if len(ret["getPackageNode"]) < 1 {
		return nil, errors.New("No package node found")
	}

	return ret["getPackageNode"][0], nil

}

func (db *DataBase) GetFileNodesByType(filetype string, recursive bool) ([]*service.FileNode, error) {
	ret := map[string][]*service.FileNode{}

	q := `query FileNodeByType($FileType: string){
		  getFileNodeByType(func: eq(nodeType, 1)) {{.Filter}} {{.Recurse}}{
			uid
			hash
			path
			derivedFrom
		  }}`

	queryTmpl, err := template.New("filenodesbytype").Parse(q)

	type QueryParams struct {
		Recurse string
		Filter  string
	}

	qp := QueryParams{}
	if recursive {
		qp.Recurse = "@recurse(loop: false)"
	}
	if filetype != "" {
		qp.Filter = "@filter(eq(type, $FileType))"
	}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		panic(err)
	}

	vars := map[string]string{"$FileType": filetype}

	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getFileNodeByType"], nil
}

func (db *DataBase) GetFileNodesByFileNode(filenode *service.FileNode, recursive bool) ([]*service.FileNode, error) {
	ret := map[string][]*service.FileNode{}

	q := `query FileNodeByFileNode($Filter: string){
		getFileNodeByFileNode(func: eq(nodeType, 1)) {{.Type}} {{.Recurse}}{
		  uid
		  hash
		  path
		  derivedFrom
		}}`

	queryTmpl, err := template.New("filenodesbyfilenode").Parse(q)

	type QueryParams struct {
		Recurse string
		Type    string
		Filter  string
	}

	qp := QueryParams{}

	if recursive {
		qp.Recurse = "@recurse(loop: false)"
	}
	if filenode.Type != "" {
		qp.Type = "@filter(eq(type, $Filter))"
		qp.Filter = filenode.Type
	}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		panic(err)
	}

	vars := map[string]string{"$Filter": qp.Filter}

	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getFileNodeByFileNode"], nil

}

func (db *DataBase) GetAnalyzerByName(name string) (*service.Analyzer, error) {
	ret := map[string][]*service.Analyzer{}

	q := `query AnalyzerByName($AnaName: string){
		  getAnalyzerByType(func: eq(nodeType, 4)) @filter(eq(name, $AnaName)) {
			uid
			hash
			path
			derivedFrom
		  }}`

	vars := map[string]string{"$AnaName": name}

	err := db.queryAnalyzer(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	if len(ret["getAnalyzerByName"]) < 1 {
		// No such analyzer
		analyzer := &service.Analyzer{Name: name, NodeType: service.NodeTypeAnalyzerNode}
		uid, err := dbInsert(db.client, analyzer)
		if err != nil {
			return nil, err
		}
		analyzer.Uid = uid
		return analyzer, nil
	}

	return ret["getAnalyzerByName"][0], nil
}

func (db *DataBase) GetNodesByType(valuetype string, recursive bool, namefilter string) ([]*service.FileNode, error) {

	ret := map[string][]*service.FileNode{}

	q := `query NodeByType($Type: string, $Name: string){
		  getNodeByType(func: eq(type, $Type)) {{.Filter}} {{.Recurse}}{
			uid
			hash
			path
			derivedFrom
		  }}`

	queryTmpl, err := template.New("nodesbytype").Parse(q)

	type QueryParams struct {
		Recurse string
		Filter  string
	}

	qp := QueryParams{}
	if recursive {
		qp.Recurse = "@recurse(loop: false)"
	}
	if namefilter != "" {
		qp.Filter = "@filter(eq(name, $Name))"
	}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		panic(err)
	}

	vars := map[string]string{"$Type": valuetype, "$Name": namefilter}

	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getNodeByType"], nil
}

func (db *DataBase) GetNodeByHash(hash string, recursive bool) (*service.FileNode, error) {

	ret := map[string][]*service.FileNode{}

	q := `query NodeByHash($Hash: string){
		  getNodeByHash(func: eq(hash, $Hash)) {
			uid
			hash
			path
		  }}`

	if recursive {
		q = `query NodeByHash($Hash: string){
		  getNodeByHash(func: eq(hash, $Hash)) @recurse(loop: false) {
			uid
			hash
			path
			derivedFrom
		  }}`
	}

	vars := map[string]string{"$Hash": hash}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getNodeByHash"][0], nil
}

func (db *DataBase) queryPackage(query string, queryVars map[string]string, resultMap *map[string][]*service.PackageNode) error {
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), query, queryVars)
	if err != nil {
		return fmt.Errorf("Could not query for package node with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", query, queryVars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal query response: %v", err)
	}
	return nil
}

func (db *DataBase) queryNodes(query string, queryVars map[string]string, resultMap *map[string][]*service.FileNode) error {
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), query, queryVars)
	if err != nil {
		return fmt.Errorf("Could not query for filenode with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", query, queryVars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal query response: %v", err)
	}
	return nil
}

func (db *DataBase) queryAnalyzer(query string, queryVars map[string]string, resultMap *map[string][]*service.Analyzer) error {
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), query, queryVars)
	if err != nil {
		return fmt.Errorf("Could not query for analyzer with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", query, queryVars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal query response: %v", err)
	}
	return nil
}
func (db *DataBase) queryInfoNodes(query string, queryVars map[string]string, resultMap *map[string][]*service.InfoNode) error {
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), query, queryVars)
	if err != nil {
		return fmt.Errorf("Could not query for info node with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", query, queryVars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal query response: %v", err)
	}
	return nil
}

package database

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"encoding/json"

	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"

	"google.golang.org/grpc"
)

const schema = `
path: string @index(trigram) .
hash: string @index(exact) .
type: string @index(hash) .
spdxIdentifier: string @index(hash) .
name: string @index(hash) .
licenseKey: string @index(hash) .
copyrightHolderName: string @index(hash) .
`

const (
	ArtifactTypeLink string = "linkedtarget"
	ArtifactTypeSrc  string = "sourcecode"
	ArtifactTypeObj  string = "object"
)

type Node struct {
	Uid             string             `json:"uid,omitempty"`
	Hash            string             `json:"hash,omitempty"`
	Type            string             `json:"type,omitempty"`
	Path            string             `json:"path,omitempty"`
	Name            string             `json:"name,omitempty"`
	DerivedFrom     []*Node            `json:"derivedFrom,omitempty"`
	License         []*License         `json:"license,omitempty"`
	CopyrightHolder []*CopyrightHolder `json:"copyrightHolder,omitempty"`
}

type CopyrightHolder struct {
	Uid  string `json:"uid,omitempty"`
	Name string `json:"copyrightHolderName"`
}

type License struct {
	Uid            string `json:"uid,omitempty"`
	Key            string `json:"licenseKey"`
	Name           string `json:"licenseName,omitempty"`
	SpdxIdentifier string `json:"spdxIdentifier,omitempty"`
}

var UnknownLicense *License = &License{Key: "UNKNOWN"}

type DataBase struct {
	client      *client.Dgraph
	insertQueue chan *Node
	insertMutex *sync.Mutex
	pending     uint64
}

func NewNode(path string, hash string) Node {
	node := Node{
		Path: path,
		Hash: hash,
		Name: filepath.Base(path),
	}
	return node
}

// Setup connects to dgraph and returns the instance
func Setup(dbAddr string, queueWorkers int) (*DataBase, error) {
	log.Println("Setting up database connection")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	conn, err := grpc.DialContext(ctx, dbAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Failed to connect to the dgraph server: %v", err)
	}

	db := &DataBase{
		client:      client.NewDgraphClient(api.NewDgraphClient(conn)),
		insertQueue: make(chan *Node, 1000),
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

	db.insertLicense(UnknownLicense)

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
		fmt.Printf("Pending inserts %d", pendingInserts)
		time.Sleep(2 * time.Second)
	}
}

// AddNode adds a node to the insert queue
func (db *DataBase) AddNode(node *Node) {
	atomic.AddUint64(&db.pending, 1)
	for _, dep := range node.DerivedFrom {
		db.AddNode(dep)
	}
	db.insertQueue <- node
}

func queueWorker(db *DataBase) {
	for node := range db.insertQueue {
		ready := true
		for idx, dep := range node.DerivedFrom {
			if dep.Uid == "" {
				// missing dep
				ready = false
				// look up dep in db
				uid, err := db.HasNode(dep.Hash)
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
			db.insertQueue <- node
			continue
		}

		// we are ready to insert the node
		db.insertMutex.Lock()
		uid, err := db.HasNode(node.Hash)
		if err != nil {
			panic(err)
		}
		if uid != "" {
			node.Uid = uid
		}
		uid, err = dbInsert(db.client, node)
		if err != nil {
			panic(err)
		}
		atomic.AddUint64(&db.pending, ^uint64(0))
		db.insertMutex.Unlock()
	}
}

func (db *DataBase) AlterNode(node *Node) (string, error) {
	db.insertMutex.Lock()
	uid, err := dbInsert(db.client, node)
	db.insertMutex.Unlock()
	return uid, err
}

// HasNode returns the UID of the node if exists otherwise ""
func (db *DataBase) HasNode(hash string) (string, error) {

	ret := map[string][]Node{}

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

func (db *DataBase) GetNodesByType(nodetype string, recursive bool, namefilter string) ([]Node, error) {

	ret := map[string][]Node{}

	q := `query NodeByType($Type: string, $Name: string){
		  getNodeByType(func: eq(type, $Type)) {{.Filter}} {{.Recurse}}{
			uid
			hash
			path
			derivedFrom
			license
			spdxIdentifier
			licenseKey
			copyrightHolder
			copyrightHolderName
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

	vars := map[string]string{"$Type": nodetype, "$Name": namefilter}

	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getNodeByType"], nil
}

func (db *DataBase) GetNodeByHash(hash string, recursive bool) (Node, error) {

	ret := map[string][]Node{}

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
			license
			spdxIdentifier
			copyrightHolder
			copyrightHolderName
		  }}`
	}

	vars := map[string]string{"$Hash": hash}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return Node{}, err
	}

	return ret["getNodeByHash"][0], nil
}

func (db *DataBase) queryNodes(query string, queryVars map[string]string, resultMap *map[string][]Node) error {
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), query, queryVars)
	if err != nil {
		return fmt.Errorf("Could not query with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", query, queryVars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal query response: %v", err)
	}
	return nil
}

// Get License will return the license uid for the given license key
func (db *DataBase) GetLicenseUid(license *License) (string, error) {
	var uid string
	ret := map[string][]License{}

	err := db.queryLicense(license.Key, &ret)
	if err != nil {
		return "", err
	}

	// license not found
	if len(ret["getLicenseByKey"]) == 0 {
		uid, err = db.insertLicense(license)
		if err != nil {
			return "", err
		}
	} else {
		uid = ret["getLicenseByKey"][0].Uid
		license.Uid = uid
	}
	return uid, nil
}

func (db *DataBase) GetCopyrightHolderUid(copyrightHolder *CopyrightHolder) (string, error) {
	var uid string
	ret := map[string][]CopyrightHolder{}

	err := db.querycopyrightHolder(copyrightHolder.Name, &ret)
	if err != nil {
		return "", err
	}

	//copyrightHolder not found
	if len(ret["getcopyrightHolderByName"]) == 0 {
		uid, err = db.insertcopyrightHolder(copyrightHolder)
		if err != nil {
			return "", err
		}
	} else {
		uid = ret["getcopyrightHolderByName"][0].Uid
		copyrightHolder.Uid = uid
	}
	return uid, nil
}

func (db *DataBase) insertLicense(license *License) (string, error) {
	var uid string
	ret := map[string][]License{}
	db.insertMutex.Lock()
	err := db.queryLicense(license.Key, &ret)
	if err != nil {
		db.insertMutex.Unlock()
		return "", err
	}
	if len(ret["getLicenseByKey"]) == 0 {
		uid, err = dbInsert(db.client, license)
		license.Uid = uid
		if err != nil {
			db.insertMutex.Unlock()
			return "", err
		}
	} else {
		uid = ret["getLicenseByKey"][0].Uid
		license.Uid = uid
	}
	db.insertMutex.Unlock()
	return uid, nil
}

func (db *DataBase) insertcopyrightHolder(copyrightHolder *CopyrightHolder) (string, error) {
	var uid string
	ret := map[string][]CopyrightHolder{}
	db.insertMutex.Lock()
	err := db.querycopyrightHolder(copyrightHolder.Name, &ret)
	if err != nil {
		log.Printf("Error querying copyrightHolder")
		db.insertMutex.Unlock()
		return "", err
	}
	if len(ret["getcopyrightHolderByName"]) == 0 {
		uid, err = dbInsert(db.client, copyrightHolder)
		copyrightHolder.Uid = uid
		if err != nil {
			log.Printf("Error inserting copyrightHolder in db")
			db.insertMutex.Unlock()
			return "", err
		}
	} else {
		uid = ret["getcopyrightHolderByName"][0].Uid
		copyrightHolder.Uid = uid
	}
	db.insertMutex.Unlock()
	return uid, nil
}

func (db *DataBase) queryLicense(licenseKey string, resultMap *map[string][]License) error {
	q := `query LicenseByKey($Key: string){
		  getLicenseByKey(func: eq(licenseKey, $Key)) {
			uid
		  }}`

	vars := map[string]string{"$Key": licenseKey}
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return fmt.Errorf("Could not query with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", q, vars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal license query response: %v", err)
	}
	return nil
}

func (db *DataBase) querycopyrightHolder(copyrightHolderName string, resultMap *map[string][]CopyrightHolder) error {
	q := `query CopyrightHolderByName($Name: string){
		getcopyrightHolderByName(func: eq(copyrightHolderName, $Name)) {
			uid
			}}`
	vars := map[string]string{"$Name": copyrightHolderName}
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return fmt.Errorf("Could not query with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", q, vars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal copyrightHolder query response: %v", err)
	}
	return nil
}

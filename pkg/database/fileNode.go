package database

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync/atomic"
	"text/template"

	"github.com/QMSTR/qmstr/pkg/service"
)

// AddBuildFileNode adds a node to the insert queue in build phase
func (db *DataBase) AddFileNode(node *service.FileNode) {
	atomic.AddUint64(&db.pending, 1)
	for _, file := range node.DerivedFrom {
		db.AddFileNode(file)
	}
	for _, dep := range node.Dependencies {
		db.AddFileNode(dep)
	}
	db.insertQueue <- node
}

// AddFileNodes stores the given FileNodes in the FileNode identified by the nodeID
func (db *DataBase) AddFileNodes(nodeID string, filenodes ...*service.FileNode) error {
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()

	const q = `
	query Node($id: string){
		node(func: uid($id)) @filter(has(fileNodeType)) @recurse(loop: false) {
			uid
			targets
			derivedFrom
		}
	}
	`
	vars := map[string]string{"$id": nodeID}
	fileNodes := map[string][]*service.FileNode{}

	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp.Json, &fileNodes)
	if err != nil {
		log.Fatal(err)
	}

	if len(fileNodes["node"]) < 1 {
		return errors.New("No such file node in graph")
	}

	fileNode := fileNodes["node"][0]
	fileNode.DerivedFrom = append(fileNode.DerivedFrom, filenodes...)

	_, err = dbInsert(db.client, fileNode)
	if err != nil {
		return err
	}

	return nil
}

// GetFileNodeUid returns the UID of the node if exists otherwise ""
func (db *DataBase) GetFileNodeUid(hash string) (string, error) {

	var ret map[string][]*service.FileNode

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

// GetFileNodesByFileNode queries filenodes on a specific attribute of a provided filenode.
// For instance, you can provide a filenode with a certain filetype and get all the filenodes
// with this filetype.
// You can query for just one attribute. For instance, if you set filetype and hash, only the
// hash will be used in the query.
func (db *DataBase) GetFileNodesByFileNode(filenode *service.FileNode, recursive bool) ([]*service.FileNode, error) {
	var ret map[string][]*service.FileNode

	q := `query FileNodeByFileNode($Filter: string, $TypeFilter: int){
			getFileNodeByFileNode(func: has(fileNodeType)) {{.Query}} {{.Recurse}}{
			  uid
			  hash
			  name
			  path
			  type
			  derivedFrom
			  dependencies
			  additionalInfo
			  confidenceScore
			  analyzer
			  dataNodes
			  data
			}}`

	queryTmpl, err := template.New("filenodesbyfilenode").Parse(q)

	type QueryParams struct {
		Recurse    string
		Query      string
		Filter     string
		TypeFilter int
	}

	qp := QueryParams{}
	vars := map[string]string{}

	if recursive {
		qp.Recurse = "@recurse(loop: false)"
	}
	if filenode.FileType != 0 {
		//get the int value from the enumeration
		t := service.FileNode_Type_value[filenode.FileType.String()]
		nt := int(t)
		qp.TypeFilter = nt
		qp.Query = "@filter(eq(fileType, $TypeFilter))"
		//convert it to string to query it
		vars["$TypeFilter"] = strconv.Itoa(nt)
	}
	if filenode.Hash != "" {
		qp.Filter = filenode.Hash
		qp.Query = "@filter(eq(hash, $Filter))"
		vars["$Filter"] = qp.Filter
	}
	if filenode.Name != "" {
		qp.Filter = filenode.Name
		qp.Query = "@filter(eq(name, $Filter))"
		vars["$Filter"] = qp.Filter
	}
	if filenode.Path != "" {
		qp.Filter = filenode.Path
		qp.Query = "@filter(eq(path, $Filter))"
		vars["$Filter"] = qp.Filter
	}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		panic(err)
	}

	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	fileNodes := ret["getFileNodeByFileNode"]
	if len(fileNodes) < 1 {
		return nil, fmt.Errorf("No file node %v found in the database", filenode)
	}
	return ret["getFileNodeByFileNode"], nil
}

func (db *DataBase) GetFileNodeHashByPath(path string) (string, error) {

	var ret map[string][]*service.FileNode

	q := `query Node($Path: string){
		  hasNode(func: eq(path, $Path)) {
			hash
		  }}`

	vars := map[string]string{"$Path": path}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return "", err
	}

	// no node with such path
	if len(ret["hasNode"]) == 0 {
		return "", errors.New("No node with such path")
	}
	return ret["hasNode"][0].Hash, nil
}

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

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

// AddFileNode adds a node to the insert queue in build phase
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
func (db *DataBase) GetFileNodeUid(path string, hash string) (string, error) {
	var ret map[string][]*service.FileNode

	q := `query Node($Path: string, $Hash: string){
		  hasNode(func: eq(path, $Path)) @cascade{
			uid
			fileData @filter(eq(hash, $Hash))
		  }}`

	vars := map[string]string{"$Path": path, "$Hash": hash}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return "", err
	}

	// no node with such path and hash
	if len(ret["hasNode"]) == 0 {
		return "", nil
	}
	return ret["hasNode"][0].Uid, nil
}

// GetFileNodesWithUIDByHash returns the UID of the fileNode
func (db *DataBase) GetFileNodesWithUIDByHash(hash string) ([]*service.FileNode, error) {
	var ret map[string][]*service.FileNode

	q := `query Node($Hash: string){
		  hasNode(func: has(fileNodeType)) @cascade{
			uid
			fileData @filter(eq(hash, $Hash))
		  }}`

	vars := map[string]string{"$Hash": hash}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	// no node with such hash
	if len(ret["hasNode"]) == 0 {
		return nil, fmt.Errorf("no file node with such hash in the database")
	}
	return ret["hasNode"], nil
}

// GetFileDataUID returns the UID of the fileData node if exists in the db
func (db *DataBase) GetFileDataUID(hash string) (string, error) {
	var ret map[string][]*service.FileNode_FileDataNode

	q := `query FileData($Hash: string){
		  hasFileData(func: eq(hash, $Hash)){
			uid
		  }}`

	vars := map[string]string{"$Hash": hash}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return "", err
	}

	// no fileData node with such hash
	if len(ret["hasFileData"]) == 0 {
		return "", nil
	}
	return ret["hasFileData"][0].Uid, nil
}

// GetFileNodesByFileNode queries filenodes on a specific attribute of a provided filenode.
// For instance, you can provide a filenode with a certain name and get all the filenodes
// with this name.
// You can query for just one attribute. For instance, if you set path and hash, only the
// hash will be used in the query.
func (db *DataBase) GetFileNodesByFileNode(in *service.GetFileNodeMessage, offset int, first int) ([]*service.FileNode, error) {
	filenode := in.FileNode
	var ret map[string][]*service.FileNode

	q := `query FileNodeByFileNode($Filter: string, $Offset: int, $First: int){
			getFileNodeByFileNode(func: has(fileNodeType), {{.Pagination}}) {{.Query}} @recurse(loop: true, depth:2){
			  uid
			  fileNodeType
			  path
			  name
			  hash
			  fileDataNodeType
			  fileData
			  timestamp
			}}`

	queryTmpl, err := template.New("filenodesbyfilenode").Parse(q)

	type QueryParams struct {
		Recurse    string
		Query      string
		Filter     string
		Offset     int
		First      int
		Pagination string
	}

	qp := QueryParams{}
	vars := map[string]string{}

	if filenode.Uid != "" {
		qp.Filter = filenode.Uid
		qp.Query = "@filter(uid($Filter))"
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
	if filenode.FileData != nil {
		filesWithUID, err := db.GetFileNodesWithUIDByHash(filenode.FileData.GetHash())
		if err != nil {
			return nil, err
		}
		if len(filesWithUID) > 1 && in.UniqueNode {
			return nil, fmt.Errorf("more than one FileNode match %v. Please provide a better identifier", filenode)
		}
		var filesWithAllData []*service.FileNode
		for _, file := range filesWithUID {
			nodeFiles, err := db.GetFileNodesByFileNode(&service.GetFileNodeMessage{FileNode: file}, 0, 500)
			if err != nil {
				return nil, err
			}
			filesWithAllData = append(filesWithAllData, nodeFiles...)
		}
		return filesWithAllData, nil
	}
	// Pagination
	qp.First, qp.Offset = first, offset
	qp.Pagination = "offset: $Offset, first: $First"
	vars["$First"] = strconv.Itoa(first)
	vars["$Offset"] = strconv.Itoa(offset)

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		return nil, err
	}

	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	fileNodes := ret["getFileNodeByFileNode"]
	if len(fileNodes) < 1 && offset == 0 {
		return nil, fmt.Errorf("No file node %v found in the database", filenode)
	}

	return fileNodes, nil
}

func (db *DataBase) GetFileNodeHashByPath(path string) (string, error) {
	var ret map[string][]*service.FileNode

	q := `query Node($Path: string){
		  hasNode(func: eq(path, $Path)) {
			fileData{
				hash
			}
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
	return ret["hasNode"][0].FileData.GetHash(), nil
}

// GetFileNodeLeaves returns the filenodes that are not derived from other filenodes
func (db *DataBase) GetFileNodeLeaves(offset int, first int) ([]*service.FileNode, error) {
	var ret map[string][]*service.FileNode

	q := `query Leaves($Offset: int, $First: int){
		  filenodeleaves(func: has(fileNodeType), first: $First, offset: $Offset) @filter(NOT has(derivedFrom)) {
			uid
			path
			fileData{
				uid
				hash
			}
		  }}`

	vars := map[string]string{
		"$First":  strconv.Itoa(first),
		"$Offset": strconv.Itoa(offset),
	}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	// no leaf file nodes
	if len(ret["filenodeleaves"]) == 0 {
		return nil, fmt.Errorf("no leaf file nodes found")
	}
	return ret["filenodeleaves"], nil
}

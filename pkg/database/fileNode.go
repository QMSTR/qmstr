package database

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync/atomic"
	"text/template"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

// AddBuildFileNode adds a node to the insert queue in build phase
func (db *DataBase) AddBuildFileNode(node *service.FileNode) {
	atomic.AddUint64(&db.pending, 1)
	for _, dep := range node.DerivedFrom {
		db.AddBuildFileNode(dep)
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

func (db *DataBase) GetFileNodesByFileNode(filenode *service.FileNode, recursive bool) ([]*service.FileNode, error) {
	ret := map[string][]*service.FileNode{}

	q := `query FileNodeByFileNode($Filter: string){
		getFileNodeByFileNode(func: has(fileNodeType)) {{.Type}} {{.Recurse}}{
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

// GetFileNodeByHash returns the file node for the file with the provided checksum
func (db *DataBase) GetFileNodeByHash(hash string, recursive bool) (*service.FileNode, error) {

	//ret := map[string][]*service.FileNode{}
	var ret map[string][]interface{}

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

	err := db.queryAnyNodes(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getNodeByHash"][0].(*service.FileNode), nil
}

func (db *DataBase) GetFileNodesByType(filetype string, recursive bool) ([]*service.FileNode, error) {
	ret := map[string][]*service.FileNode{}

	q := `query FileNodeByType($FileType: string){
		  getFileNodeByType(func: has(fileNodeType)) {{.Filter}} {{.Recurse}}{
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

package database

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"sync/atomic"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

// AddBuildFileNode adds a node to the insert queue in build phase
func (db *DataBase) AddFileNode(node *service.FileNode) {
	atomic.AddUint64(&db.pending, 1)
	for _, dep := range node.DerivedFrom {
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

func (db *DataBase) GetFileNodesByFileNode(filenode *service.FileNode) ([]*service.FileNode, error) {
	var ret map[string][]*service.FileNode

	q := `query FileNodeByFileNode($Type: int){
		getFileNodeByFileNode(func: eq(fileType, $Type)) @recurse(loop: false){
		  uid
		  hash
		  path
		  derivedFrom
		}}`

	//get the int value from the enumeration
	t := service.FileNode_Type_value[filenode.FileType.String()]
	nt := int(t)
	//convert it to string to query it
	vars := map[string]string{"$Type": strconv.Itoa(nt)}
	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getFileNodeByFileNode"], nil
}

// GetFileNodeByHash returns the file node for the file with the provided checksum
func (db *DataBase) GetFileNodeByHash(hash string, recursive bool) (*service.FileNode, error) {

	var ret map[string][]*service.FileNode

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

//
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

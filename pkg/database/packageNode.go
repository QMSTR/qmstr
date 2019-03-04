package database

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync/atomic"

	"github.com/QMSTR/qmstr/pkg/service"
)

var ErrNoSuchPackage = errors.New("no such package")

// AddPackageNode adds a node to the insert queue
func (db *DataBase) AddPackageNode(node *service.PackageNode) {
	atomic.AddUint64(&db.pending, 1)
	for _, dep := range node.Targets {
		db.AddFileNode(dep)
	}
	db.insertQueue <- node
}

func (db *DataBase) GetPackageNodeByName(name string) (*service.PackageNode, error) {
	var ret map[string][]*service.PackageNode

	q := `query GetPackageNodeByName($Name: string){
		  getPackageNodeByName(func: has(packageNodeType)) @filter(eq(name, $Name)) {
			uid
			name
		  }}`

	vars := map[string]string{"$Name": name}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	// no such node
	if len(ret["getPackageNodeByName"]) == 0 {
		return nil, ErrNoSuchPackage
	}
	return ret["getPackageNodeByName"][0], nil
}

func (db *DataBase) GetPackageNode() (*service.PackageNode, error) {
	var ret map[string][]*service.PackageNode

	q := `{
		getPackageNode(func: has(packageNodeType)) @recurse(loop: false) {
			uid
			buildConfig
			hash
			name
			type
			targets
			derivedFrom
			dependencies
			fileType
			path
			additionalInfo
			confidenceScore
			analyzer
			dataNodes
			data
		  }
		}`

	err := db.queryNodesSimple(q, &ret)
	if err != nil {
		return nil, err
	}

	pkgNodes := ret["getPackageNode"]
	if len(pkgNodes) < 1 {
		return nil, errors.New("No package node found")
	}

	return pkgNodes[0], nil

}

// AddPackageFileNodes stores the given FileNodes in a PackageNode identified by the nodeID
func (db *DataBase) AddPackageFileNodes(nodeID string, filenodes ...*service.FileNode) error {
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()

	const q = `
	query Node($id: string){
		node(func: uid($id)) @filter(has(packageNodeType)) @recurse(loop: false) {
			uid
			targets
		}
	}
	`
	vars := map[string]string{"$id": nodeID}
	packageNodes := map[string][]*service.PackageNode{}

	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp.Json, &packageNodes)
	if err != nil {
		log.Fatal(err)
	}

	if len(packageNodes["node"]) < 1 {
		return errors.New("No such package node in graph")
	}

	packageNode := packageNodes["node"][0]

	packageNode.Targets = append(packageNode.Targets, filenodes...)

	_, err = dbInsert(db.client, packageNode)
	if err != nil {
		return err
	}

	return nil
}

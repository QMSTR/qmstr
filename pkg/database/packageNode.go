package database

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/QMSTR/qmstr/pkg/service"
)

func (db *DataBase) AddPackageNode(pNode *service.PackageNode) (string, error) {
	log.Printf("Adding package [%v] node to the graph", pNode)
	return dbInsert(db.client, pNode)
}
func (db *DataBase) GetPackageNode(session string) (*service.PackageNode, error) {
	ret := map[string][]*service.PackageNode{}

	q := `query PackageNode($Session: string) {
		getPackageNode(func: has(packageNodeType)) @recurse(loop: false) {
			uid
			session
			buildConfig
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

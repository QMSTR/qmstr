package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/QMSTR/go-qmstr/service"
)

// AddDiagnosticNodes stores the given DiagnosticNodes in a PackageNode or FileNode identified by the nodeID
func (db *DataBase) AddDiagnosticNodes(nodeID string, diagnosticnodes ...*service.DiagnosticNode) error {
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()

	const q = `
	query Node($id: string){
		node(func: uid($id)) @filter(has(packageNodeType) or has(fileNodeType)) @recurse(loop: false) {
			uid
			diagnosticInfo
			packageNodeType
			fileNodeType
		}
	}
	`
	vars := map[string]string{"$id": nodeID}
	var result map[string][]interface{}

	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return err
	}

	err = json.Unmarshal(resp.Json, &result)
	if err != nil {
		log.Fatal(err)
	}

	if len(result["node"]) < 1 {
		return fmt.Errorf("No package or file node with uid %s found", nodeID)
	}

	receiverNode := result["node"][0].(map[string]interface{})
	var diagnosticInfo []*service.DiagnosticNode
	if diagnosticInfoInter, ok := receiverNode["diagnosticInfo"]; ok {
		diagnosticInfo = diagnosticInfoInter.([]*service.DiagnosticNode)
	}
	diagnosticInfo = append(diagnosticInfo, diagnosticnodes...)

	if _, ok := receiverNode["packageNodeType"]; ok {
		packageNode := service.PackageNode{}
		packageNode.Uid = nodeID
		packageNode.DiagnosticInfo = diagnosticInfo
		_, err = dbInsert(db.client, &packageNode)
		if err != nil {
			return err
		}
	} else {
		fileNode := service.FileNode{}
		fileNode.Uid = nodeID
		fileNode.DiagnosticInfo = diagnosticInfo
		_, err = dbInsert(db.client, &fileNode)
		if err != nil {
			return err
		}
	}
	return nil
}

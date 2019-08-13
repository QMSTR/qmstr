package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"text/template"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
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
	} else if _, ok := receiverNode["fileNodeType"]; ok {
		fileNode := service.FileNode{}
		fileNode.Uid = nodeID
		fileNode.FileData.DiagnosticInfo = diagnosticInfo
		_, err = dbInsert(db.client, &fileNode)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("wrong type of node. Can't connect diagnostic nodes to it")
	}
	return nil
}

//GetDiagnosticNodeBySeverity queries diagnostic nodes on a specific severity
func (db *DataBase) GetDiagnosticNodeBySeverity(diagnosticNode *service.DiagnosticNode) ([]*service.DiagnosticNode, error) {
	var ret map[string][]*service.DiagnosticNode

	const q = `query DiagnosticData($Severity: int){
		getDiagnosticData(func: has(diagnosticNodeType)) @filter(eq(severity, $Severity)) {
			diagnosticInfo
			message
		}}`

	queryTmpl, err := template.New("diagnosticnodebyseverity").Parse(q)

	type QueryParams struct {
		Severity int
	}

	qp := QueryParams{}
	//get the int value from the enumeration
	t := service.DiagnosticNode_Severity_value[diagnosticNode.Severity.String()]
	nt := int(t)
	qp.Severity = nt

	//convert it to string to query it
	vars := map[string]string{"$Severity": strconv.Itoa(nt)}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		panic(err)
	}
	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	messages := ret["getDiagnosticData"]
	if len(messages) < 1 {
		return nil, fmt.Errorf("No diagnostic node %v found in the database", strconv.Itoa(nt))
	}
	return messages, nil
}

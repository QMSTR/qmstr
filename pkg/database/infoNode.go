package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

// AddInfoNodes stores the given InfoNodes in a PackageNode or FileNode identified by the nodeID
func (db *DataBase) AddInfoNodes(nodeID string, infonodes ...*service.InfoNode) error {
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()

	const q = `
	query Node($id: string){
		node(func: uid($id)) @filter(has(packageNodeType) or has(fileNodeType)) @recurse(loop: false) {
			uid
			additionalInfo
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
	var additionalInfo []*service.InfoNode
	if additionalInfoInter, ok := receiverNode["AdditionalInfo"]; ok {
		additionalInfo = additionalInfoInter.([]*service.InfoNode)
	}
	additionalInfo = append(additionalInfo, infonodes...)

	if _, ok := receiverNode["packageNodeType"]; ok {
		packageNode := service.PackageNode{}
		packageNode.Uid = nodeID
		packageNode.AdditionalInfo = additionalInfo
		_, err = dbInsert(db.client, &packageNode)
		if err != nil {
			return err
		}
	} else {
		fileNode := service.FileNode{}
		fileNode.Uid = nodeID
		fileNode.AdditionalInfo = additionalInfo
		_, err = dbInsert(db.client, &fileNode)
		if err != nil {
			return err
		}
	}
	return nil
}

func (db *DataBase) GetInfoData(rootNodeID string, infotype string, datatype string) ([]string, error) {
	const q = `
	query InfoData($id: string, $itype: string, $dtype: string){
		var(func: uid($id)) @recurse(loop: false) {
			uid
			name:type
			analyzer
			trustLevel
			A as dataNodes @filter(eq(type, $dtype))
			data
			confidenceScore
			additionalInfo @filter(eq(type, $itype)) (orderdesc: confidenceScore, first: 1)
			derivedFrom
		}

		infodata(func: uid(A)) {
			data
		}
	}
	`

	vars := map[string]string{"$id": rootNodeID, "$itype": infotype, "$dtype": datatype}

	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return nil, err
	}

	type Data struct {
		Data string
	}

	type InfoData struct {
		Infodata []Data
	}

	var r InfoData

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Fatal(err)
	}

	ret := []string{}
	for _, data := range r.Infodata {
		ret = append(ret, data.Data)
	}
	return ret, nil
}

func (db *DataBase) GetInfoNodeByDataNode(infonodetype string, datanodes ...*service.InfoNode_DataNode) (*service.InfoNode, error) {

	var retInfoNode *service.InfoNode

	runeDataNodeMap := map[string]*service.InfoNode_DataNode{}

	for idx, datanode := range datanodes {
		runeDataNodeMap[getVarName(idx)] = datanode
	}

	var ret map[string]interface{}

	q := `query InfoNodeByDataNode($InfoType: string) {
				{{range $var, $data := .}}
				var(func: has(infoNodeType)) @filter(eq(type, "{{$data.Type}}") AND eq(data, "{{$data.Data}}")) {
					{{$var}} as ~dataNodes
				}
				{{end}}
		
				getInfoByData(func: has(infoNodeType)) @filter(eq(type, $InfoType) {{range $var, $data := .}} AND uid({{$var}}) {{end}}) {
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

	err = db.queryNodes(queryString, vars, &ret)
	if err != nil {
		return nil, err
	}

	infoNodes := ret["getInfoByData"].([]*service.InfoNode)

	if len(infoNodes) > 0 {
		retInfoNode = infoNodes[0]
	}

	if retInfoNode == nil {
		infoNode := &service.InfoNode{Type: infonodetype}
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

func (db *DataBase) GetAllInfoData(infotype string) ([]string, error) {
	const q = `
	query InfoData($itype: string){
		getInfoData(func: has(infoNodeType))  @filter(eq(type, $itype)) {
			A as dataNodes
		}

		infodata(func: uid(A)) {
			data
		}
	}
	`
	vars := map[string]string{"$itype": infotype}

	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return nil, err
	}

	type Data struct {
		Data string
	}

	type InfoData struct {
		Infodata []Data
	}

	var r InfoData

	err = json.Unmarshal(resp.Json, &r)
	if err != nil {
		log.Fatal(err)
	}

	ret := []string{}
	for _, data := range r.Infodata {
		ret = append(ret, data.Data)
	}
	return ret, nil
}

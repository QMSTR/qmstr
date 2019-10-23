package database

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"text/template"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

// AddInfoNodes stores the given InfoNodes in a PackageNode or FileNode identified by the nodeID
func (db *DataBase) AddInfoNodes(nodeID string, infonodes []*service.InfoNode) error {
	db.insertMutex.Lock()
	defer db.insertMutex.Unlock()

	const q = `
	query Node($id: string){
		node(func: uid($id)) @filter(has(projectNodeType) or has(packageNodeType) or has(fileDataNodeType)) @recurse(loop: false) {
			uid
			projectNodeType
			packageNodeType
			fileDataNodeType
			additionalInfo
			analyzer
			name
			type
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

	if additionalInfoInter, ok := receiverNode["additionalInfo"]; ok {
		// check if these info nodes have already been inserted in the db
		for _, infoNode := range additionalInfoInter.([]interface{}) {
			for attrName, attrValue := range infoNode.(map[string]interface{}) {
				if attrName == "analyzer" {
					for _, analyzer := range attrValue.([]interface{}) {
						for name, value := range analyzer.(map[string]interface{}) {
							if name == "name" {
								if value.(string) == infonodes[0].Analyzer[0].Name {
									log.Printf("Analyzer %v already created info nodes for file %s, skipping insert..", infonodes[0].Analyzer[0].Name, nodeID)
									return nil

								}
							}
						}
					}
				}
			}
		}
	}
	if _, ok := receiverNode["projectNodeType"]; ok {
		projectNode := service.ProjectNode{}
		projectNode.Uid = nodeID
		projectNode.AdditionalInfo = infonodes
		_, err = dbInsert(db.client, &projectNode)
		if err != nil {
			return err
		}
	} else if _, ok := receiverNode["packageNodeType"]; ok {
		packageNode := service.PackageNode{}
		packageNode.Uid = nodeID
		packageNode.AdditionalInfo = infonodes
		_, err = dbInsert(db.client, &packageNode)
		if err != nil {
			return err
		}
	} else if _, ok := receiverNode["fileDataNodeType"]; ok {
		fileDataNode := service.FileNode_FileDataNode{}
		fileDataNode.Uid = nodeID
		fileDataNode.AdditionalInfo = infonodes
		_, err = dbInsert(db.client, &fileDataNode)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("wrong type of node. Can't connect info nodes to it")
	}
	return nil
}

// GetInfoDataByTrustLevel returns infonodes containing the datanodes detected from the most trusted analyzer
func (db *DataBase) GetInfoDataByTrustLevel(fileID string, infotype string, datatype string) ([]string, error) {
	var ret map[string][]*service.InfoNode

	const q = `query InfoData($ID: string, $Itype: string, $Dtype: string){
		var(func: uid($ID)) @recurse(loop: false) {
			name
			fileData
			targets
			derivedFrom
			T as additionalInfo @filter(eq(type, $Itype))(orderdesc: confidenceScore, first: 1)
		}
		var(func: uid(T)){
			type
			analyzer{
				tr as trustLevel
				name
			}
		}
		var(){
			A as trustLevel: max(val(tr))
		}
		getInfoData(func: uid(T)) @recurse(loop: false) {
			name
			type
			analyzer @filter(eq(trustLevel, val(A)))
			dataNodes @filter(eq(type, $Dtype))
			data
		}
	  }`

	queryTmpl, err := template.New("infodatabytrustlevel").Parse(q)

	type QueryParams struct {
		ID    string
		Itype string
		Dtype string
	}
	qp := QueryParams{}

	vars := map[string]string{"$ID": fileID, "$Itype": infotype, "$Dtype": datatype}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		return nil, err
	}
	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	infoData := ret["getInfoData"]
	if len(infoData) < 1 {
		return nil, nil
	}

	realData := []string{}

	for _, info := range infoData {
		// infoData contains all the infodata attached to the filenode (with the declared info type)
		// but the query returns only the most trusted analyzer connected to the infonodes
		// So only info nodes with an analyzer attached are the trustworthy data
		if len(info.Analyzer) > 0 {
			for _, data := range info.DataNodes {
				realData = append(realData, data.Data)
			}
		}
	}

	return realData, nil
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

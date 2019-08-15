package database

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
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

// GetFileNodeUidByHash returns the UID of the fileNode
func (db *DataBase) GetFileNodeUidByHash(hash string) (string, error) {
	var ret map[string][]*service.FileNode

	q := `query Node($Hash: string){
		  hasNode(func: has(fileNodeType)) @cascade{
			uid
			fileData @filter(eq(hash, $Hash))
		  }}`

	vars := map[string]string{"$Hash": hash}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return "", err
	}

	// no node with such hash
	if len(ret["hasNode"]) == 0 {
		return "", fmt.Errorf("no file node with such hash in the database")
	}
	return ret["hasNode"][0].Uid, nil
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
// For instance, you can provide a filenode with a certain filetype and get all the filenodes
// with this filetype.
// You can query for just one attribute. For instance, if you set filetype and hash, only the
// hash will be used in the query.
func (db *DataBase) GetFileNodesByFileNode(filenode *service.FileNode, recursive bool) ([]*service.FileNode, error) {
	var ret map[string][]interface{}

	q := `query FileNodeByFileNode($Filter: string, $TypeFilter: int){
			getFileNodeByFileNode(func: has(fileNodeType)) {{.Query}} {{.Recurse}}{
			  uid
			  fileNodeType
			  path
			  name
			  hash
			  fileDataNodeType
			  fileData
			  type
			  timestamp
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
		nodeUID, err := db.GetFileNodeUidByHash(filenode.FileData.Hash)
		if err != nil {
			return nil, err
		}
		qp.Filter = nodeUID
		qp.Query = "@filter(uid($Filter))"
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

	fileNodesInterface := ret["getFileNodeByFileNode"]
	if len(fileNodesInterface) < 1 {
		return nil, fmt.Errorf("No file node %v found in the database", filenode)
	}

	var fileNodes []*service.FileNode
	for _, node := range fileNodesInterface {
		nodeMap := node.(map[string]interface{})
		result := &service.FileNode{}
		if err := decodeToNodeStruct(result, nodeMap); err != nil {
			return nil, err
		}
		fileNodes = append(fileNodes, result)
	}
	return fileNodes, nil
}

// TODO: Use dgraph v1.1.0 when it's released
// and remove the decodeToNodeStruct convenience function
func decodeToNodeStruct(s interface{}, nodeMap map[string]interface{}) error {
	structValue := reflect.ValueOf(s).Elem()
	for name, value := range nodeMap {
		field := reflect.ValueOf(value)
		switch field.Kind() {
		case reflect.Slice, reflect.Array:
			if name == "fileData" {
				fileData := &service.FileNode_FileDataNode{}
				for _, dataNode := range value.([]interface{}) {
					if err := decodeToNodeStruct(fileData, dataNode.(map[string]interface{})); err != nil {
						return err
					}
				}
				s.(*service.FileNode).FileData = fileData
				continue
			} else if name == "derivedFrom" || name == "dependencies" || name == "targets" {
				deps, err := getDependenciesStruct(value.([]interface{}))
				if err != nil {
					return err
				}
				value = deps
			} else if name == "additionalInfo" {
				var infos []*service.InfoNode
				for _, infoNode := range value.([]interface{}) {
					info := &service.InfoNode{}
					if err := decodeToNodeStruct(info, infoNode.(map[string]interface{})); err != nil {
						return err
					}
					infos = append(infos, info)
				}
				value = infos
			} else if name == "diagnosticInfo" {
				var infos []*service.DiagnosticNode
				for _, infoNode := range value.([]interface{}) {
					info := &service.DiagnosticNode{}
					if err := decodeToNodeStruct(info, infoNode.(map[string]interface{})); err != nil {
						return err
					}
					infos = append(infos, info)
				}
				value = infos
			} else if name == "dataNodes" {
				var dataNodes []*service.InfoNode_DataNode
				for _, dataNode := range value.([]interface{}) {
					data := &service.InfoNode_DataNode{}
					if err := decodeToNodeStruct(data, dataNode.(map[string]interface{})); err != nil {
						return err
					}
					dataNodes = append(dataNodes, data)
				}
				value = dataNodes
			} else if name == "analyzer" {
				var analyzers []*service.Analyzer
				for _, analyzer := range value.([]interface{}) {
					analyzerData := &service.Analyzer{}
					if err := decodeToNodeStruct(analyzerData, analyzer.(map[string]interface{})); err != nil {
						return err
					}
					analyzers = append(analyzers, analyzerData)
				}
				value = analyzers
			}
		}
		structFieldValue := structValue.FieldByName(strings.Title(name))
		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in node", name)
		}
		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", name)
		}
		if name == "fileType" {
			value = service.FileNode_Type(value.(float64))
		}
		val := reflect.ValueOf(value)
		if structFieldValue.Type() != val.Type() {
			return fmt.Errorf("Provided value type: %v: %v, didn't match node field type: %v", val.Type(), val, structFieldValue.Type())
		}
		structFieldValue.Set(val)
	}
	return nil
}

func getDependenciesStruct(nodes []interface{}) ([]*service.FileNode, error) {
	var deps []*service.FileNode
	for _, depNode := range nodes {
		dep := &service.FileNode{}
		if err := decodeToNodeStruct(dep, depNode.(map[string]interface{})); err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}
	return deps, nil
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
	return ret["hasNode"][0].FileData.Hash, nil
}

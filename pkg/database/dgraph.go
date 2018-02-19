package database

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"log"
	"path/filepath"
	"sync"
	"time"

	"encoding/json"

	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"

	"google.golang.org/grpc"
)

const schema = `
path: string @index(trigram) .
hash: string @index(exact) .
type: string @index(hash) .
spdxIdentifier: string @index(hash) .
name: string @index(hash) .
`

const (
	ArtifactTypeLink string = "linkedtarget"
	ArtifactTypeSrc  string = "sourcecode"
	ArtifactTypeObj  string = "object"
)

type Node struct {
	Uid         string     `json:"uid,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	Type        string     `json:"type,omitempty"`
	Path        string     `json:"path,omitempty"`
	Name        string     `json:"name,omitempty"`
	DerivedFrom []*Node    `json:"derivedFrom,omitempty"`
	License     []*License `json:"license,omitempty"`
}

type License struct {
	Uid            string `json:"uid,omitempty"`
	SpdxIdentifier string `json:"spdxIdentifier,omitempty"`
}

type DataBase struct {
	client      *client.Dgraph
	insertQueue chan *Node
	insertMutex *sync.Mutex
}

func NewNode(path string, hash string) Node {
	node := Node{
		Path: path,
		Hash: hash,
		Name: filepath.Base(path),
	}
	return node
}

// Setup connects to dgraph and returns the instance
func Setup(dbAddr string) (*DataBase, error) {
	log.Println("Setting up database connection")
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	conn, err := grpc.DialContext(ctx, dbAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		cancel()
		return nil, fmt.Errorf("Failed to connect to the dgraph server: %v", err)
	}

	db := &DataBase{
		client:      client.NewDgraphClient(api.NewDgraphClient(conn)),
		insertQueue: make(chan *Node, 1000),
		insertMutex: &sync.Mutex{},
	}

	for {
		err = db.client.Alter(ctx, &api.Operation{
			Schema: schema,
		})
		if err == nil {
			cancel()
			break
		}
	}

	go queueWorker(db)
	go queueWorker(db)
	go queueWorker(db)

	return db, nil
}

// AddNode adds a node to the insert queue
func (db *DataBase) AddNode(node *Node) {
	for _, dep := range node.DerivedFrom {
		db.AddNode(dep)
	}
	db.insertQueue <- node
}

func queueWorker(db *DataBase) {
	for node := range db.insertQueue {
		ready := true
		for idx, dep := range node.DerivedFrom {
			if dep.Uid == "" {
				// missing dep
				ready = false
				// look up dep in db
				uid, err := db.HasNode(dep.Hash)
				if err != nil {
					panic(err)
				}
				// found uid
				if uid != "" {
					node.DerivedFrom[idx].Uid = uid
				}
			}
		}

		if !ready {
			// put node back to queue
			db.insertQueue <- node
			continue
		}

		// we are ready to insert the node
		db.insertMutex.Lock()
		uid, err := db.HasNode(node.Hash)
		if err != nil {
			panic(err)
		}
		if uid != "" {
			node.Uid = uid
		}
		uid, err = dbInsert(db.client, node)
		if err != nil {
			panic(err)
		}
		db.insertMutex.Unlock()
	}
}

func (db *DataBase) AlterNode(node *Node) (string, error) {
	db.insertMutex.Lock()
	uid, err := dbInsert(db.client, node)
	db.insertMutex.Unlock()
	return uid, err
}

// HasNode returns the UID of the node if exists otherwise ""
func (db *DataBase) HasNode(hash string) (string, error) {

	ret := map[string][]Node{}

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

// the data should be JSON marshalable
func dbInsert(c *client.Dgraph, data interface{}) (string, error) {
	txn := c.NewTxn()
	defer txn.Discard(context.Background())

	jNode, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	ret, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: true, SetJson: jNode})
	if err != nil {
		return "", err
	}

	uid := ret.Uids["blank-0"]
	return uid, nil
}

func (db *DataBase) GetNodesByType(nodetype string, recursive bool, namefilter string) ([]Node, error) {

	ret := map[string][]Node{}

	q := `query NodeByType($Type: string, $Name: string){
		  getNodeByType(func: eq(type, $Type)) {{.Filter}} {{.Recurse}}{
			uid
			hash
			path
			derivedFrom
			license
			spdxIdentifier
		  }}`

	queryTmpl, err := template.New("nodesbytype").Parse(q)

	type QueryParams struct {
		Recurse string
		Filter  string
	}

	qp := QueryParams{}
	if recursive {
		qp.Recurse = "@recurse(loop: false)"
	}
	if namefilter != "" {
		qp.Filter = "@filter(eq(name, $Name))"
	}

	var b bytes.Buffer
	err = queryTmpl.Execute(&b, qp)
	if err != nil {
		panic(err)
	}

	vars := map[string]string{"$Type": nodetype, "$Name": namefilter}

	err = db.queryNodes(b.String(), vars, &ret)
	if err != nil {
		return nil, err
	}

	return ret["getNodeByType"], nil
}

func (db *DataBase) GetNodeByHash(hash string, recursive bool) (Node, error) {

	ret := map[string][]Node{}

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
			license
			spdxIdentifier
		  }}`
	}

	vars := map[string]string{"$Hash": hash}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return Node{}, err
	}

	return ret["getNodeByHash"][0], nil
}

func (db *DataBase) queryNodes(query string, queryVars map[string]string, resultMap *map[string][]Node) error {
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), query, queryVars)
	if err != nil {
		return fmt.Errorf("Could not query with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", query, queryVars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal query response: %v", err)
	}
	return nil
}

// Get License will return the license uid for the given spdx license identifier
func (db *DataBase) GetLicenseUid(license string) (string, error) {

	var uid string
	ret := map[string][]License{}

	err := db.queryLicense(license, &ret)
	if err != nil {
		return "", err
	}

	// license not found
	if len(ret["getLicenseByName"]) == 0 {
		uid, err = db.insertLicense(license)
		if err != nil {
			return "", err
		}
	} else {
		uid = ret["getLicenseByName"][0].Uid
	}
	return uid, nil
}

func (db *DataBase) insertLicense(license string) (string, error) {
	var uid string
	ret := map[string][]License{}
	db.insertMutex.Lock()
	err := db.queryLicense(license, &ret)
	if err != nil {
		return "", err
	}
	if len(ret["getLicenseByName"]) == 0 {
		uid, err = dbInsert(db.client, License{SpdxIdentifier: license})
		if err != nil {
			return "", err
		}
	} else {
		uid = ret["getLicenseByName"][0].Uid
	}
	db.insertMutex.Unlock()
	return uid, nil
}

func (db *DataBase) queryLicense(license string, resultMap *map[string][]License) error {
	q := `query LicenseByName($Name: string){
		  getLicenseByName(func: eq(spdxIdentifier, $Name)) {
			uid
		  }}`

	vars := map[string]string{"$Name": license}
	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return fmt.Errorf("Could not query with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", q, vars, err)
	}

	if err = json.Unmarshal(resp.Json, resultMap); err != nil {
		return fmt.Errorf("Could not unmashal license query response: %v", err)
	}
	return nil
}

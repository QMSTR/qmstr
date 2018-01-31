package database

import (
	"context"
	"fmt"

	"encoding/json"

	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"

	"google.golang.org/grpc"
)

const schema = `
path: string @index(trigram) .
hash: string @index(exact) .
type: string @index(hash) .
`

const (
	ArtifactTypeLink string = "linkedtarget"
	ArtifactTypeSrc  string = "sourcecode"
	ArtifactTypeObj  string = "object"
)

type Node struct {
	Uid         string  `json:"uid,omitempty"`
	Hash        string  `json:"hash,omitempty"`
	Type        string  `json:"type,omitempty"`
	Path        string  `json:"path,omitempty"`
	DerivedFrom []Node  `json:"derivedFrom,omitempty"`
	License     License `json:"license,omitempty"`
}

type License struct {
	Uid  string `json:"uid,omitempty"`
	Name string `json:"name,omitempty"`
}

type DataBase struct {
	client *client.Dgraph
}

// Setup connects to dgraph and returns the instance
func Setup(dbAddr string) (*DataBase, error) {

	conn, err := grpc.Dial(dbAddr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the dgraph server: %v", err)
	}

	db := &DataBase{
		client: client.NewDgraphClient(api.NewDgraphClient(conn)),
	}

	err = db.client.Alter(context.Background(), &api.Operation{
		Schema: schema,
	})
	if err != nil {
		return nil, fmt.Errorf("Fail to set schema and indices: %v", err)
	}
	return db, nil
}

// AddNode adds a node to the DB and returns it's UID
func (db *DataBase) AddNode(node *Node) (string, error) {
	return dbInser(db.client, node)
}

// HasNode returns the UID of the node if exists otherwise ""
func (db *DataBase) HasNode(hash string) (string, error) {

	ret := map[string][]Node{}

	q := `query Node($Hash: string){
		  hasNode(func: eq(hash, $Hash)) {
			uid
		  }}`

	vars := map[string]string{"$Hash": hash}

	resp, err := db.client.NewTxn().QueryWithVars(context.Background(), q, vars)
	if err != nil {
		return "", fmt.Errorf("Could not query with: \n\n%s\n\nVars:\n\n%v\n\nError: %v", q, vars, err)
	}

	if err = json.Unmarshal(resp.Json, &ret); err != nil {
		return "", fmt.Errorf("Could not unmashal `hasNode` response: %v", err)
	}

	// no node with such hash
	if len(ret["hasNode"]) == 0 {
		return "", nil
	}
	return ret["hasNode"][0].Uid, nil
}

// the data should be JSON marshalable
func dbInser(c *client.Dgraph, data interface{}) (string, error) {
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

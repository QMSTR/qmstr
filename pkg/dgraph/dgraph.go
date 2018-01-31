package dgraph

import (
	"context"
	"log"

	"encoding/json"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"
)

const schema = `
path: string @index(trigram) .
hash: string @index(exact) .
type: string @index(hash) .
`

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

// Setup Dgraph schema
func Setup(c *client.Dgraph) {
	err := c.Alter(context.Background(), &api.Operation{
		Schema: schema,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func AddData(c *client.Dgraph, buildMessage *pb.BuildMessage) {
	txn := c.NewTxn()
	defer txn.Discard(context.Background())

	q := `query Node($Hash: string){
		hasNode(func: eq(hash, $Hash)) {
			uid
			hash
		}}`

	for _, comp := range buildMessage.GetCompilations() {
		variables := map[string]string{"$Hash": comp.Target.GetHash()}

		resp, err := txn.QueryWithVars(context.Background(), q, variables)
		if err != nil {
			log.Fatal(err)
		}

		// Alter the existing Node

		type NodeC struct {
			Node []Node `json:"hasNode"`
		}

		var t NodeC
		err = json.Unmarshal(resp.Json, &t)
		if err != nil {
			log.Fatal(err)
		}
		// if no target in DB
		if len(t.Node) == 0 {
			variables["$Hash"] = comp.Source.GetHash()
			resp, err = txn.QueryWithVars(context.Background(), q, variables)
			if err != nil {
				log.Fatal(err)
			}

			var s NodeC
			err = json.Unmarshal(resp.Json, &s)
			if err != nil {
				log.Fatal(err)
			}
			var src Node
			// no source in DB
			if len(s.Node) == 0 {
				src = Node{
					Hash: comp.Source.GetHash(),
					Path: comp.Source.GetPath(),
				}
			} else {
				src = s.Node[0]
			}
			targ := Node{
				Hash:        comp.Target.GetHash(),
				Path:        comp.Target.GetPath(),
				DerivedFrom: []Node{src},
			}

			uid, err := addToDB(txn, &targ)
			if err != nil {
				log.Fatalf("Adding node to DB failed: %v\n", err)
			}
			log.Printf("UID: %s\n", uid)
		}
	}

	if len(buildMessage.Binary) > 0 {
		bin := buildMessage.Binary[0]
		variables := map[string]string{"$Hash": bin.Target.GetHash()}

		resp, err := txn.QueryWithVars(context.Background(), q, variables)
		if err != nil {
			log.Fatal(err)
		}

		// Alter the existing Node

		type NodeC struct {
			Node []Node `json:"hasNode"`
		}

		var t NodeC
		err = json.Unmarshal(resp.Json, &t)
		if err != nil {
			log.Fatal(err)
		}

		deps := []Node{}
		// if no target in DB
		if len(t.Node) == 0 {
			for _, dep := range bin.Input {
				variables["$Hash"] = dep.GetHash()
				resp, err = txn.QueryWithVars(context.Background(), q, variables)
				if err != nil {
					log.Fatal(err)
				}

				var s NodeC
				err = json.Unmarshal(resp.Json, &s)
				if err != nil {
					log.Fatal(err)
				}
				var src Node
				// no source in DB
				if len(s.Node) == 0 {
					src = Node{
						Hash: dep.GetHash(),
						Path: dep.GetPath(),
					}
				} else {
					src = s.Node[0]
				}
				deps = append(deps, src)
			}

			targ := Node{
				Hash:        bin.Target.GetHash(),
				Path:        bin.Target.GetPath(),
				DerivedFrom: deps,
			}

			uid, err := addToDB(txn, &targ)
			if err != nil {
				log.Fatalf("Adding node to DB failed: %v\n", err)
			}
			log.Printf("UID: %s\n", uid)

		}

	}

	txn.Commit(context.Background())
}

func addToDB(txn *client.Txn, node *Node) (string, error) {
	jNode, err := json.Marshal(node)
	if err != nil {
		return "", err
	}

	ret, err := txn.Mutate(context.Background(), &api.Mutation{CommitNow: false, SetJson: jNode})
	if err != nil {
		return "", err
	}

	uid := ret.Uids["blank-0"]
	return uid, nil
}

package dgraph

import (
	"context"
	"log"

	"encoding/json"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"
)

const Sch = "path: string @index(exact) ."

type Target struct {
	Uid          string   `json:"uid,omitempty"`
	Path         string   `json:"path,omitempty"`
	Source       string   `json:"source,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// Setup Dgraph schema
func Setup(c *client.Dgraph) {
	err := c.Alter(context.Background(), &api.Operation{
		Schema: Sch,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func AddData(c *client.Dgraph, buildMessage *pb.BuildMessage) {
	txn := c.NewTxn()
	defer txn.Discard(context.Background())
	// Loop through the targets and the sources
	for _, comp := range buildMessage.GetCompilations() {
		a := Target{
			Path:   comp.Target.GetPath(),
			Source: comp.Source.GetPath(),
		}
		out, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Failed to marsal the json encoding: %v", err)
		}
		_, err = txn.Mutate(context.Background(), &api.Mutation{CommitNow: false, SetJson: out})
		if err != nil {
			log.Fatalf("Failed to mutate the data: %v", err)
		}
	}
	// Loop through the targets and the dependencies
	for _, link := range buildMessage.GetBinary() {
		deps := []string{}
		for _, dep := range link.Input {
			deps = append(deps, dep.Path)
		}
		a := Target{
			Path:         link.Target.GetPath(),
			Dependencies: deps,
		}
		out, err := json.Marshal(a)
		if err != nil {
			log.Fatalf("Failed to marsal the json encoding: %v", err)
		}
		_, err = txn.Mutate(context.Background(), &api.Mutation{CommitNow: false, SetJson: out})
		if err != nil {
			log.Fatalf("Failed to mutate the data: %v", err)
		}
	}
	txn.Commit(context.Background())
}

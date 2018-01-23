package dgraph

import (
	"context"
	"log"

	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"
)

type Target struct {
	Uid          string   `json:"uid,omitempty"`
	Path         string   `json:"path,omitempty"`
	Source       string   `json:"source,omitempty"`
	Dependencies []string `json:"dependencies,omitempty"`
}

// Setup Dgraph schema
func Setup(c *client.Dgraph) {
	err := c.Alter(context.Background(), &api.Operation{
		Schema: `
                path: string @index(exact) . 
		`,
	})
	if err != nil {
		log.Fatal(err)
	}
}

package cli

import (
	"log"
	"reflect"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new node",
	Long:  "create a new node described by an node identifier",
	Run: func(cmd *cobra.Command, args []string) {
		setUpServer()
		createNode(args)
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createNode(args []string) {
	for _, nodeident := range args {
		node, err := ParseNodeID(nodeident)
		if err != nil {
			log.Fatalf("%v", err)
		}
		switch reflect.TypeOf(node) {
		case reflect.TypeOf((*service.FileNode)(nil)):
			log.Printf("Got node %v", node.(*service.FileNode))
		}
	}
}

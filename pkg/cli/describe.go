package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	packageName, target string
)

var describeCmd = &cobra.Command{
	Use:   "describe [type_of_node:node]",
	Short: "Print description of the node",
	Long: `Print description of the node and traverse the tree 
to print the description of the nodes connected to it.

input: [type_of_node:node], where type_of_node can be:
	- package
	- target 
	- info
and node, can be:
	- node name
	- node path 
	- node type`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpServer()
		if err := describeNode(args); err != nil {
			Log.Fatalf("Describe failed: %v", err)
		}
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(describeCmd)
}

func describeNode(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("Please provide the node you want to get description for: [type_of_node:node]")
	}

	input := strings.Split(args[0], ":")
	if len(input) != 2 {
		return fmt.Errorf("Please provide input in the format: [type_of_node:node]")
	}
	nodeType := input[0]
	node := input[1]

	queryNode := &service.FileNode{}

	switch nodeType {
	case "package":
		//query package name
		fmt.Printf("Type %s is not yet implemented", nodeType)
	case "target":
		path := strings.Contains(node, "/")
		if path {
			queryNode = &service.FileNode{Path: node}
		} else {
			queryNode = &service.FileNode{Name: node}
		}

		stream, err := controlServiceClient.GetFileNodeDescription(context.Background(), queryNode)
		if err != nil {
			log.Printf("Could not get file node %v", err)
			return err
		}

		for {
			fileNode, err := stream.Recv()
			if err == io.EOF {
				break
			}

			json, err := json.MarshalIndent(fileNode, "", "   ")
			fmt.Printf("%v \n", string(json))
		}
	case "info":
		//TODO: query for infonode
		fmt.Printf("Type %s is not yet implemented", nodeType)
	}

	return nil
}

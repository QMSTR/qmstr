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
	less                bool
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
	describeCmd.Flags().BoolVar(&less, "less", false, "show less information-info nodes are not traversed")
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

	switch nodeType {
	case "package":
		//query package name
		fmt.Printf("Type %s is not yet implemented", nodeType)
	case "target":
		queryNode := &service.FileNode{}
		path := strings.Contains(node, "/")
		if path {
			queryNode = &service.FileNode{Path: node}
		} else {
			queryNode = &service.FileNode{Name: node}
		}
		descriptionRequest := &service.FileDescriptionRequest{File: queryNode, LessInfo: less}
		stream, err := controlServiceClient.GetFileNodeDescription(context.Background(), descriptionRequest)
		if err != nil {
			log.Printf("Could not get file node %v", err)
			return err
		}

		for {
			fileNode, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				if err.Error() == "rpc error: code = Unknown desc = No file node found" {
					return fmt.Errorf("No file node %s found in the database", node)
				}
				return err
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

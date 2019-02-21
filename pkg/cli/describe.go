package cli

import (
	"fmt"
	"io"
	"strings"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

const (
	PKGNODE  = "package"
	FILENODE = "target"
)

var less bool

var describeCmd = &cobra.Command{
	Use:   "describe [type_of_node:node]",
	Short: "Print description of the node",
	Long: `Print description of the node and traverse the tree 
to print the description of the nodes connected to it.

input: [type_of_node:node], where type_of_node can be:
	- package (no need to provide node)
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
	nodeType := input[0]

	switch nodeType {
	case PKGNODE:
		pkgNode, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageRequest{})
		if err != nil {
			return err
		}
		fmt.Println(pkgNode.Describe(less))
	case FILENODE:
		if len(input) != 2 {
			return fmt.Errorf("Please provide input in the format: [type_of_node:node]")
		}
		node := input[1]
		if node == "" {
			return fmt.Errorf("Please provide input in the format: [type_of_node:node]")
		}
		queryNode := &service.FileNode{}
		path := strings.Contains(node, "/")
		if path {
			queryNode = &service.FileNode{Path: node}
		} else {
			queryNode = &service.FileNode{Name: node}
		}
		fNodes, err := controlServiceClient.GetFileNode(context.Background(), queryNode)
		if err != nil {
			return err
		}

		for {
			fileNode, err := fNodes.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return err
			}
			fmt.Println(fileNode.Describe(less, ""))
		}
		//TODO: case INFONODE
	default:
		return fmt.Errorf("Wrong input: %s. No such type of node", nodeType)
	}
	return nil
}

package cli

import (
	"fmt"
	"io"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var less bool

var describeCmd = &cobra.Command{
	Use:   "describe [type_of_node:attribute:value]",
	Short: "Print description of the node",
	Long: `Print description of the node and traverse the tree 
to print the description of the nodes connected to it.

input: [type_of_node:attribute:value], where type_of_node can be:
	- package
	- target 
attribute can be:
	- name
	- path 
	- type
	- hash
and value, the value of the attribute.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpBuildService()
		setUpControlService()
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
		return fmt.Errorf("Please provide the node you want to get description for: [type_of_node:attribute:value]")
	}
	n, err := ParseNodeID(args[0])
	if err != nil {
		return err
	}

	switch node := n.(type) {
	case *service.ProjectNode:
		pNode, err := buildServiceClient.GetProjectNode(context.Background(), &service.ProjectNode{})
		if err != nil {
			return err
		}
		fmt.Println(pNode.Describe(less))
	case *service.PackageNode:
		pkgNode, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageNode{})
		if err != nil {
			return err
		}
		fmt.Println(pkgNode.Describe(less, ""))
	case *service.FileNode:
		fNodes, err := controlServiceClient.GetFileNode(context.Background(), &service.GetFileNodeMessage{FileNode: node})
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
	default:
		return fmt.Errorf("Unsupported node: %s", node)
	}
	return nil
}

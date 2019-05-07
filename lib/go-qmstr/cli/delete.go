package cli

import (
	"errors"
	"fmt"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [type_of_node:attribute:value]",
	Short: "Delete node from the database",
	Long:  `Delete the provided node from the database.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		setUpBuildService()
		setUpControlService()
		awaitServer()
		if err := deleteNode(args); err != nil {
			Log.Fatalf("Delete failed: %v", err)
		}
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteNode(args []string) error {
	deleteNodeMsg := []*service.DeleteMessage{}
	// loop through the nodes we are going to delete
	for _, arg := range args {
		node, err := ParseNodeID(arg)
		if err != nil {
			return err
		}
		switch currentNode := node.(type) {
		case *service.ProjectNode:
			projectNode, err := buildServiceClient.GetProjectNode(context.Background(), &service.ProjectNode{Name: currentNode.Name})
			if err != nil {
				return err
			}
			fmt.Printf("Deleting project node: %v\n", projectNode.Name)
			deleteNodeMsg = append(deleteNodeMsg, &service.DeleteMessage{Uid: projectNode.Uid})
		case *service.PackageNode:
			pkgNode, err := getUniquePackageNode(currentNode)
			if err != nil {
				return err
			}
			deleteNodeMsg = append(deleteNodeMsg, &service.DeleteMessage{Uid: pkgNode.Uid})
		case *service.FileNode:
			fNode, err := getUniqueFileNode(currentNode)
			if err != nil {
				return err
			}
			deleteNodeMsg = append(deleteNodeMsg, &service.DeleteMessage{Uid: fNode.Uid})
		}
	}
	deleteStream, err := buildServiceClient.DeleteNode(context.Background())
	if err != nil {
		return err
	}

	for _, dltMsg := range deleteNodeMsg {
		if verbose {
			fmt.Printf("Deleting node: %v\n", dltMsg.Uid)
		}
		deleteStream.Send(dltMsg)
	}
	reply, err := deleteStream.CloseAndRecv()
	if err != nil {
		return err
	}
	if !reply.Success {
		return errors.New("failed deleting nodes")
	}
	return nil
}

// qmstrctl file:name:curl

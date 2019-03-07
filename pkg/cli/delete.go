package cli

import (
	"errors"
	"fmt"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var deleteCmd = &cobra.Command{
	Use:   "delete [type_of_node:attribute:value]",
	Short: "Delete node from the database",
	Long:  `Delete the provided node from the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpBuildService()
		setUpControlService()
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
	if len(args) < 1 {
		return fmt.Errorf("Please provide the node to be deleted: [type_of_node:attribute:value]")
	}
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
			pkgNode, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageNode{Name: currentNode.Name})
			if err != nil {
				return err
			}
			fmt.Printf("Deleting package node: %v\n", pkgNode.Name)
			deleteNodeMsg = append(deleteNodeMsg, &service.DeleteMessage{Uid: pkgNode.Uid})
		case *service.FileNode:
			fNode, err := getUniqueFileNode(currentNode)
			if err != nil {
				return fmt.Errorf("get unique file node fail. please use better matching params: %v", err)
			}
			fmt.Printf("Deleting file node: %v\n", fNode.Path)
			deleteNodeMsg = append(deleteNodeMsg, &service.DeleteMessage{Uid: fNode.Uid})
		}
	}
	deleteStream, err := buildServiceClient.DeleteNode(context.Background())
	if err != nil {
		return err
	}

	for _, dltMsg := range deleteNodeMsg {
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

package cli

import (
	"errors"
	"fmt"
	"io"

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
		switch node.(type) {
		case *service.PackageNode:
			pkgNode, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageRequest{})
			if err != nil {
				return err
			}
			fmt.Println(pkgNode)
		case *service.FileNode:
			fNodes, err := controlServiceClient.GetFileNode(context.Background(), node.(*service.FileNode))
			if err != nil {
				return err
			}
			fileN := &service.FileNode{}
			var count int
			for {
				fileNode, err := fNodes.Recv()
				if err != nil {
					if err == io.EOF {
						break
					}
					return err
				}
				count++
				fileN = fileNode
			}
			if fileN.Uid == "" {
				return fmt.Errorf("No file %v found in the database", node)
			}
			if count > 1 {
				return fmt.Errorf("Found more than one %v in database\n Please provide a better identifier", node)
			}
			deleteNodeMsg = append(deleteNodeMsg, &service.DeleteMessage{Uid: fileN.Uid})
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

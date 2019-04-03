package cli

import (
	"fmt"
	"reflect"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var disconnectCmdFlags = struct {
	edge string
}{}

var disconnectCmd = &cobra.Command{
	Use:   "disconnect [type_of_node:attribute:value] [type_of_node:attribute:value]",
	Short: "Disconnect nodes with specific edges",
	Long: `Usage: qmstrctl disconnect <that> <this>...
Disconnect from Node <that> Node(s) <this>.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpControlService()
		setUpBuildService()
		if err := disconnectNodes(cmd, args); err != nil {
			Log.Fatalf("Disonnect failed: %v", err)
		}
		tearDownServer()
	},
	Args: cobra.MinimumNArgs(2),
}

func init() {
	rootCmd.AddCommand(disconnectCmd)
	disconnectCmd.Flags().StringVar(&disconnectCmdFlags.edge, "edge", "", "Edge to use when disconnecting nodes")
}

func disconnectNodes(cmd *cobra.Command, args []string) error {
	thatID, err := ParseNodeID(args[0])
	if err != nil {
		return fmt.Errorf("ParseNodeID fail for %q: %v", args[1], err)
	}
	switch thatVal := thatID.(type) {
	case *service.FileNode:
		that, err := getUniqueFileNode(thatVal)
		if err != nil {
			return fmt.Errorf("get unique \"that\" node fail: %v", err)
		}
		err = disconnectFromFileNode(that, args[1:])
		if err != nil {
			return fmt.Errorf("disconnect file nodes fail: %v", err)
		}
	case *service.PackageNode:
		_, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageNode{})
		if err != nil {
			return fmt.Errorf("get package node fail: %v", err)
		}
	default:
		return fmt.Errorf("unsuported node type %T", thatVal)
	}
	return nil
}

func disconnectFromFileNode(node *service.FileNode, args []string) error {
	var deleteNodeMsg *service.DeleteMessage
	for _, nID := range args {
		thisID, err := ParseNodeID(nID)
		if err != nil {
			return fmt.Errorf("ParseNodeID fail for %q: %v", args[0], err)
		}
		switch thisVal := thisID.(type) {
		// FileNode -> FileNode
		case *service.FileNode:
			this, err := getUniqueFileNode(thisVal)
			if err != nil {
				return fmt.Errorf("get unique file node fail. please use better matching params: %v", err)
			}
			// default edge
			if disconnectCmdFlags.edge == "" {
				disconnectCmdFlags.edge = "derivedFrom"
			}
			// Which edge
			switch disconnectCmdFlags.edge {
			case "derivedFrom":
				deleteNodeMsg = &service.DeleteMessage{Uid: node.Uid, Edge: "derivedFrom"}
				for i, dr := range node.DerivedFrom {
					if reflect.DeepEqual(this, dr) {
						node.DerivedFrom = append(node.DerivedFrom[:i], node.DerivedFrom[i+1:]...)
					}
				}
			case "dependencies":
				deleteNodeMsg = &service.DeleteMessage{Uid: node.Uid, Edge: "dependencies"}
				for i, dep := range node.Dependencies {
					if dep == this {
						node.Dependencies = append(node.Dependencies[:i], node.Dependencies[i+1:]...)
					}
				}
			default:
				return fmt.Errorf("unknown edge %q for FileNode -> FileNode. Valid values %v", disconnectCmdFlags.edge, validFileToFileEdges)
			}
		default:
			return fmt.Errorf("cannot disconnect %T from FileNode", thisVal)
		}
	}
	// delete edge
	buildServiceClient.DeleteEdge(context.Background(), deleteNodeMsg)

	// ship node back with the modified edge
	err := sendFileNode(node)
	if err != nil {
		return fmt.Errorf("sending FileNode fail: %v", err)
	}
	return nil
}

func disconnectFromPackageNode(node *service.PackageNode, args []string) {}

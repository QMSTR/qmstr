package cli

import (
	"fmt"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
)

var disconnectCmdFlags = struct {
	fileNodeToFileNodeEdge string
}{}

func init() {
	rootCmd.AddCommand(disconnectCmd)
	disconnectCmd.Flags().StringVar(&disconnectCmdFlags.fileNodeToFileNodeEdge, "fileNodeToFileNodeEdge",
		"derivedFrom", fmt.Sprintf("Edge to use when disconnecting FileNode to FileNode. One of %v", validFileNodeToFileNodeEdges))
}

var disconnectCmd = &cobra.Command{
	Use:   "disconnect",
	Short: "qmstrctl disconnect <this> <that>",
	Long:  "Disconnect Node <this> from Node <that>",
	Run: func(cmd *cobra.Command, args []string) {
		setUpControlService()
		setUpBuildService()
		if err := disconnectCmdRun(cmd, args); err != nil {
			Log.Fatalf("Disconnect failed: %v", err)
		}
		tearDownServer()
	},
	Args: cobra.ExactArgs(2),
}

func disconnectCmdRun(cmd *cobra.Command, args []string) error {
	thisID, err := ParseNodeID(args[0])
	if err != nil {
		return fmt.Errorf("ParseNodeID fail for %q: %v", args[0], err)
	}
	thatID, err := ParseNodeID(args[1])
	if err != nil {
		return fmt.Errorf("ParseNodeID fail for %q: %v", args[1], err)
	}

	switch thatVal := thatID.(type) {
	case *service.FileNode:
		that, err := getUniqueFileNode(thatVal)
		if err != nil {
			return fmt.Errorf("get unique \"that\" node fail: %v", err)
		}
		switch thisVal := thisID.(type) {
		// FileNode -> FileNode
		case *service.FileNode:
			this, err := getUniqueFileNode(thisVal)
			if err != nil {
				return fmt.Errorf("get unique \"this\" node fail: %v", err)
			}
			// Which edge
			switch disconnectCmdFlags.fileNodeToFileNodeEdge {
			case "derivedFrom":
				that.DerivedFrom, err = removeFileNodeFromList(that.DerivedFrom, this)
			case "dependencies":
				that.Dependencies, err = removeFileNodeFromList(that.Dependencies, this)
			default:
				return fmt.Errorf("unknown edge for FileNode -> FileNode: %s", disconnectCmdFlags.fileNodeToFileNodeEdge)
			}
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("cannot disconnect %T to FileNode", thisVal)
		}
		// ship it back
		err = sendFileNode(that)
		if err != nil {
			return fmt.Errorf("sending FileNode fail: %v", err)
		}
	default:
		return fmt.Errorf("unsuported node type %T", thatVal)
	}
	return nil
}

func removeFileNodeFromList(list []*service.FileNode, node *service.FileNode) ([]*service.FileNode, error) {
	for i, n := range list {
		if n.Uid == node.Uid {
			if i != len(list)-1 {
				return append(list[:i], list[i+1:]...), nil
			}
			return list[:i], nil
		}
	}
	return nil, fmt.Errorf("nodes not connected")
}

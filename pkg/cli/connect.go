package cli

import (
	"fmt"
	"io"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var validFileNodeToFileNodeEdges = []string{
	"derivedFrom",
	"dependencies",
}

var connectCmdFlags = struct {
	fileNodeToFileNodeEdge string
}{}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVar(&connectCmdFlags.fileNodeToFileNodeEdge, "fileNodeToFileNodeEdge",
		"derivedFrom", fmt.Sprintf("Edge to use when connecting FileNode to FileNode. One of %v", validFileNodeToFileNodeEdges))
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "qmstrctl connect <this> <that>",
	Long:  "Connect Node <this> to Node <that>",
	Run: func(cmd *cobra.Command, args []string) {
		setUpControlService()
		setUpBuildService()
		if err := connectCmdRun(cmd, args); err != nil {
			Log.Fatalf("Connect failed: %v", err)
		}
		tearDownServer()
	},
	Args: cobra.ExactArgs(2),
}

func connectCmdRun(cmd *cobra.Command, args []string) error {
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
			switch connectCmdFlags.fileNodeToFileNodeEdge {
			case "derivedFrom":
				that.DerivedFrom = append(that.DerivedFrom, this)
			case "dependencies":
				that.Dependencies = append(that.Dependencies, this)
			default:
				return fmt.Errorf("unknown edge for FileNode -> FileNode: %s", connectCmdFlags.fileNodeToFileNodeEdge)
			}
		default:
			return fmt.Errorf("cannot connect %T to FileNode", thisVal)
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

func getUniqueFileNode(queryNode *service.FileNode) (*service.FileNode, error) {
	stream, err := controlServiceClient.GetFileNode(context.Background(), queryNode)
	node, err := stream.Recv()
	if err != nil {
		return nil, fmt.Errorf("recv FileNode fail for %v: %v", queryNode, err)
	}
	_, err = stream.Recv()
	if err == nil {
		return nil, fmt.Errorf("more than one FileNode match %v", queryNode)
	}
	if err != io.EOF {
		return nil, fmt.Errorf("probbing for more nodes fail: %v", err)
	}
	return node, nil
}

func sendFileNode(node *service.FileNode) error {
	stream, err := buildServiceClient.Build(context.Background())
	if err != nil {
		return fmt.Errorf("getting stream for build service fail: %v", err)
	}
	return stream.Send(node)
}

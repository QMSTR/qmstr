package cli

import (
	"fmt"
	"io"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var validFileToFileEdges = []string{
	"derivedFrom",
	"dependencies",
}

var validFileToPackageEdges = []string{
	"targets",
}

var connectCmdFlags = struct {
	edge string
}{}

func init() {
	rootCmd.AddCommand(connectCmd)
	connectCmd.Flags().StringVar(&connectCmdFlags.edge, "edge", "", "Edge to use when connecting nodes")
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect nodes with specific edges",
	Long: `Usage: qmstrctl connect <that> <this>...
Connect to Node <that> Node(s) <this>. In case of multiple edges for the specified types you can use --<type>To<type>Edge flag to specify the edge you want.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpControlService()
		setUpBuildService()
		if err := connectCmdRun(cmd, args); err != nil {
			Log.Fatalf("Connect failed: %v", err)
		}
		tearDownServer()
	},
	Args: cobra.MinimumNArgs(2),
}

func connectCmdRun(cmd *cobra.Command, args []string) error {
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
		err = connectToFileNode(that, args[1:])
		if err != nil {
			return fmt.Errorf("connectToFileNode fail: %v", err)
		}
	case *service.PackageNode:
		that, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageRequest{})
		if err != nil {
			return fmt.Errorf("get package node fail: %v", err)
		}
		err = connectToPackageNode(that, args[1:])
		if err != nil {
			return fmt.Errorf("connectToPackageNode fail: %v", err)
		}
	default:
		return fmt.Errorf("unsuported node type %T", thatVal)
	}
	return nil
}

func connectToFileNode(node *service.FileNode, args []string) error {
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
			if connectCmdFlags.edge == "" {
				connectCmdFlags.edge = "derivedFrom"
			}
			// Which edge
			switch connectCmdFlags.edge {
			case "derivedFrom":
				node.DerivedFrom = append(node.DerivedFrom, this)
			case "dependencies":
				node.Dependencies = append(node.Dependencies, this)
			default:
				return fmt.Errorf("unknown edge %q for FileNode -> FileNode. Valid values %v", connectCmdFlags.edge, validFileToFileEdges)
			}
		default:
			return fmt.Errorf("cannot connect %T to FileNode", thisVal)
		}
	}
	// ship it back
	err := sendFileNode(node)
	if err != nil {
		return fmt.Errorf("sending FileNode fail: %v", err)
	}
	return nil
}

func connectToPackageNode(node *service.PackageNode, args []string) error {
	stream, err := buildServiceClient.Package(context.Background())
	if err != nil {
		return err
	}
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
			if connectCmdFlags.edge == "" {
				connectCmdFlags.edge = "targets"
			}
			// Which edge
			switch connectCmdFlags.edge {
			case "targets":
				err = stream.Send(this)
				if err != nil {
					return fmt.Errorf("send fileNode to pkg stream fail: %v", err)
				}
			default:
				return fmt.Errorf("unknown edge %q for FileNode -> PackageNode. Valid values %v", connectCmdFlags.edge, validFileToPackageEdges)
			}
		default:
			return fmt.Errorf("cannot connect %T to FileNode", thisVal)
		}
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close stream fail: %v", err)
	}
	if !res.Success {
		return fmt.Errorf("sending node fail: %v", err)
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
	err = stream.Send(node)
	if err != nil {
		return fmt.Errorf("sending node fail: %v", err)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close stream fail: %v", err)
	}
	if !res.Success {
		return fmt.Errorf("sending node fail: %v", err)
	}
	return nil
}

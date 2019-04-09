package cli

import (
	"fmt"

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
			return err
		}
		err = connectToFileNode(that, args[1:])
		if err != nil {
			return fmt.Errorf("connectToFileNode fail: %v", err)
		}
	case *service.PackageNode:
		that, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageNode{})
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

func connectToFileNode(that *service.FileNode, args []string) error {
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
				return err
			}
			err = addFileNodeEdge(that, this)
			if err != nil {
				return err
			}
		default:
			return fmt.Errorf("cannot connect %T to FileNode", thisVal)
		}
	}
	// ship it back
	err := sendFileNode(that)
	if err != nil {
		return fmt.Errorf("sending FileNode fail: %v", err)
	}
	return nil
}

func connectToPackageNode(that *service.PackageNode, args []string) error {
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
				return err
			}
			if connectCmdFlags.edge != "" && connectCmdFlags.edge != "targets" {
				return fmt.Errorf("unknown edge %q for FileNode -> PackageNode. Valid values %v", connectCmdFlags.edge, validFileToPackageEdges)
			}
			err = stream.Send(this)
			if err != nil {
				return fmt.Errorf("send fileNode to pkg stream fail: %v", err)
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

func addFileNodeEdge(that *service.FileNode, this *service.FileNode) error {
	// default edge
	if connectCmdFlags.edge == "" {
		connectCmdFlags.edge = "derivedFrom"
	}
	// Which edge
	switch connectCmdFlags.edge {
	case "derivedFrom":
		that.DerivedFrom = append(that.DerivedFrom, this)
	case "dependencies":
		that.Dependencies = append(that.Dependencies, this)
	default:
		return fmt.Errorf("unknown edge %q for FileNode -> FileNode. Valid values %v", connectCmdFlags.edge, validFileToFileEdges)
	}
	return nil
}

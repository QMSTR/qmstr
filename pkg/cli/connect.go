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

var validPackageToProjectEdges = []string{
	"packages",
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
		return fmt.Errorf("Failed parsing node %q: %v", args[0], err)
	}
	these, err := getNodesFromArgs(args[1:])
	if err != nil {
		return err
	}
	switch thatVal := thatID.(type) {
	case *service.FileNode:
		that, err := getUniqueFileNode(thatVal)
		if err != nil {
			return err
		}
		theseFileNodes, err := createFileNodesArray(these)
		if err != nil {
			return fmt.Errorf("Only file nodes can be connected to file node: %v", err)
		}
		err = connectToFileNode(that, theseFileNodes)
		if err != nil {
			return fmt.Errorf("Failed connecting file nodes: %v", err)
		}
	case *service.PackageNode:
		that, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageNode{})
		if err != nil {
			return fmt.Errorf("Failed to get package node: %v", err)
		}
		theseFileNodes, err := createFileNodesArray(these)
		if err != nil {
			return fmt.Errorf("Only file nodes can be connected to package node: %v", err)
		}
		err = connectToPackageNode(that, theseFileNodes)
		if err != nil {
			return fmt.Errorf("Failed connecting file nodes to package node: %v", err)
		}
	case *service.ProjectNode:
		that, err := buildServiceClient.GetProjectNode(context.Background(), &service.ProjectNode{Name: thatVal.Name})
		if err != nil {
			return fmt.Errorf("Failed to get project node: %v", err)
		}
		thesePackageNodes, err := createPkgNodesArray(these)
		if err != nil {
			return fmt.Errorf("Only package nodes can be connected to Project node: %v", err)
		}
		err = connectToProjectNode(that, thesePackageNodes)
		if err != nil {
			return fmt.Errorf("Failed connecting package nodes to project node: %v", err)
		}
	default:
		return fmt.Errorf("unsupported node type %T", thatVal)
	}
	return nil
}

func connectToFileNode(that *service.FileNode, these []*service.FileNode) error {
	err := addFileNodeEdge(that, these)
	if err != nil {
		return err
	}
	// ship it back
	err = sendFileNode(that)
	if err != nil {
		return fmt.Errorf("Failed sending FileNode: %v", err)
	}
	return nil
}

func connectToPackageNode(that *service.PackageNode, these []*service.FileNode) error {
	if connectCmdFlags.edge != "" && connectCmdFlags.edge != "targets" {
		return fmt.Errorf("unknown edge %q for FileNode -> PackageNode. Valid values %v", connectCmdFlags.edge, validFileToPackageEdges)
	}
	stream, err := buildServiceClient.Package(context.Background())
	if err != nil {
		return err
	}
	for _, this := range these {
		err = stream.Send(this)
		if err != nil {
			return fmt.Errorf("Failed sending targets: %v", err)
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

func connectToProjectNode(that *service.ProjectNode, these []*service.PackageNode) error {
	if connectCmdFlags.edge != "" && connectCmdFlags.edge != "packages" {
		return fmt.Errorf("unknown edge %q for PackageNode -> ProjectNode. Valid values %v", connectCmdFlags.edge, validPackageToProjectEdges)
	}
	stream, err := buildServiceClient.UpdateProjectNode(context.Background())
	if err != nil {
		return err
	}
	for _, this := range these {
		err = stream.Send(this)
		if err != nil {
			return fmt.Errorf("Failed sending packages: %v", err)
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

func addFileNodeEdge(that *service.FileNode, these []*service.FileNode) error {
	// default edge
	if connectCmdFlags.edge == "" {
		connectCmdFlags.edge = "derivedFrom"
	}
	// Which edge
	switch connectCmdFlags.edge {
	case "derivedFrom":
		that.DerivedFrom = append(that.DerivedFrom, these...)
	case "dependencies":
		that.Dependencies = append(that.Dependencies, these...)
	default:
		return fmt.Errorf("unknown edge %q for FileNode -> FileNode. Valid values %v", connectCmdFlags.edge, validFileToFileEdges)
	}
	return nil
}

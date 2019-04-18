package cli

import (
	"fmt"
	"reflect"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var disconnectCmdFlags = struct {
	edge string
}{}
var deleteNodeMsg *service.DeleteMessage

var disconnectCmd = &cobra.Command{
	Use:   "disconnect [type_of_node:attribute:value] [type_of_node:attribute:value]",
	Short: "Disconnect nodes with specific edges",
	Long: `Usage: qmstrctl disconnect <that> <this>...
Disconnect from Node <that> Node(s) <this>.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpControlService()
		setUpBuildService()
		if err := disconnectNodes(cmd, args); err != nil {
			Log.Fatalf("Disconnect failed: %v", err)
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
			return fmt.Errorf("Only file nodes can be disconnected from file node: %v", err)
		}
		err = disconnectFromFileNode(that, theseFileNodes)
		if err != nil {
			return fmt.Errorf("Failed disconnecting from file node: %v", err)
		}
	case *service.PackageNode:
		that, err := getUniquePackageNode(thatVal)
		if err != nil {
			return err
		}
		theseFileNodes, err := createFileNodesArray(these)
		if err != nil {
			return fmt.Errorf("Only file nodes can be disconnected from package node: %v", err)
		}
		err = disconnectFromPackageNode(that, theseFileNodes)
		if err != nil {
			return fmt.Errorf("Failed disconnecting files from package node: %v", err)
		}
	case *service.ProjectNode:
		that, err := buildServiceClient.GetProjectNode(context.Background(), &service.ProjectNode{})
		if err != nil {
			return fmt.Errorf("Failed to get project node: %v", err)
		}
		thesePkgNodes, err := createPkgNodesArray(these)
		if err != nil {
			return fmt.Errorf("Only package nodes can be disconnected from project node: %v", err)
		}
		err = disconnectFromProjectNode(that, thesePkgNodes)
		if err != nil {
			return fmt.Errorf("Failed disconnecting packages from project node: %v", err)
		}
	default:
		return fmt.Errorf("unsupported node type %T", thatVal)
	}
	return nil
}

func disconnectFromFileNode(that *service.FileNode, these []*service.FileNode) error {
	err := removeFileNodePredicates(that, these)
	if err != nil {
		return err
	}

	// delete edge
	res, err := buildServiceClient.DeleteEdge(context.Background(), deleteNodeMsg)
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("Failed deleting edge: %v", err)
	}

	// ship node back with the modified edge
	err = sendFileNode(that)
	if err != nil {
		return fmt.Errorf("Failed sending FileNode: %v", err)
	}
	return nil
}

func disconnectFromPackageNode(that *service.PackageNode, these []*service.FileNode) error {
	err := removePackageNodePredicates(that, these)
	if err != nil {
		return err
	}
	// delete edge
	res, err := buildServiceClient.DeleteEdge(context.Background(), deleteNodeMsg)
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("Failed deleting edge: %v", err)
	}
	res, err = buildServiceClient.UpdatePackageNode(context.Background(), &service.UpdatePackageNodeMessage{Package: that, Targets: that.Targets})
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("sending package node failed: %v", err)
	}
	return nil
}

func disconnectFromProjectNode(that *service.ProjectNode, these []*service.PackageNode) error {
	err := removeProjectNodePredicates(that, these)
	if err != nil {
		return err
	}
	// delete edge
	res, err := buildServiceClient.DeleteEdge(context.Background(), deleteNodeMsg)
	if err != nil {
		return err
	}
	if !res.Success {
		return fmt.Errorf("Failed deleting edge: %v", err)
	}
	err = updateProjectNode(that, that.Packages)
	if err != nil {
		return err
	}
	return nil
}

func removeFileNodePredicates(that *service.FileNode, these []*service.FileNode) error {
	// default edge
	if disconnectCmdFlags.edge == "" {
		disconnectCmdFlags.edge = "derivedFrom"
	}
	switch disconnectCmdFlags.edge {
	case "derivedFrom":
		deleteNodeMsg = &service.DeleteMessage{Uid: that.Uid, Edge: "derivedFrom"}
		for _, this := range these {
			for i, dr := range that.DerivedFrom {
				if reflect.DeepEqual(this, dr) {
					that.DerivedFrom = append(that.DerivedFrom[:i], that.DerivedFrom[i+1:]...)
				}
			}
		}
	case "dependencies":
		deleteNodeMsg = &service.DeleteMessage{Uid: that.Uid, Edge: "dependencies"}
		for _, this := range these {
			for i, dep := range that.Dependencies {
				if reflect.DeepEqual(this, dep) {
					that.Dependencies = append(that.Dependencies[:i], that.Dependencies[i+1:]...)
				}
			}
		}
	default:
		return fmt.Errorf("unknown edge %q for FileNode -> FileNode. Valid values %v", disconnectCmdFlags.edge, validFileToFileEdges)
	}
	return nil
}

func removePackageNodePredicates(that *service.PackageNode, these []*service.FileNode) error {
	if connectCmdFlags.edge != "" && connectCmdFlags.edge != "targets" {
		return fmt.Errorf("unknown edge %q for FileNode -> PackageNode. Valid values %v", disconnectCmdFlags.edge, validFileToPackageEdges)
	}
	deleteNodeMsg = &service.DeleteMessage{Uid: that.Uid, Edge: "targets"}
	for _, this := range these {
		for i, target := range that.Targets {
			if reflect.DeepEqual(this, target) {
				that.Targets = append(that.Targets[:i], that.Targets[i+1:]...)
			}
		}
	}
	return nil
}

func removeProjectNodePredicates(that *service.ProjectNode, these []*service.PackageNode) error {
	if connectCmdFlags.edge != "" && connectCmdFlags.edge != "packages" {
		return fmt.Errorf("unknown edge %q for PackageNode -> ProjectNode. Valid values %v", disconnectCmdFlags.edge, validPackageToProjectEdges)
	}

	deleteNodeMsg = &service.DeleteMessage{Uid: that.Uid, Edge: "packages"}
	for _, this := range these {
		for i, pkg := range that.Packages {
			if reflect.DeepEqual(this, pkg) {
				that.Packages = append(that.Packages[:i], that.Packages[i+1:]...)
			}
		}
	}
	return nil
}

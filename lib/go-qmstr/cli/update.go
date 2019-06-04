package cli

import (
	"errors"
	"fmt"
	"log"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/net/context"
)

var updateCmd = &cobra.Command{
	Use:   "update [type_of_node:attribute:value]",
	Short: "Update nodes",
	Long:  "Update nodes in the database by adding or modifying their attributes.",
}

var updateFileCmd = &cobra.Command{
	Use:   "file",
	Short: "update a file node",
	Long:  "update a file node described by a node identifier",
	Run:   update,
}

var updatePkgCmd = &cobra.Command{
	Use:   "package",
	Short: "update a package node",
	Long:  "update a package node described by a node identifier",
	Run:   update,
}

var updateProjCmd = &cobra.Command{
	Use:   "project",
	Short: "update a project node",
	Long:  "update a project node described by a node identifier",
	Run:   update,
}

func update(cmd *cobra.Command, args []string) {
	setUpBuildService()
	setUpControlService()
	awaitServer()
	cmdFlags = cmd.Flags()
	err := updateNode(args[0])
	if err != nil {
		log.Fatalf("Failed to update node: %v", err)
	}
}

func init() {
	var err error
	generatedFlags := &pflag.FlagSet{}
	rootCmd.AddCommand(updateCmd)

	err = generateFlags(&service.FileNode{}, generatedFlags)
	if err != nil {
		log.Fatalf("%v", err)
	}
	updateFileCmd.Flags().AddFlagSet(generatedFlags)
	updateCmd.AddCommand(updateFileCmd)

	generatedFlags = &pflag.FlagSet{}
	err = generateFlags(&service.PackageNode{}, generatedFlags)
	if err != nil {
		log.Fatalf("%v", err)
	}
	updatePkgCmd.Flags().AddFlagSet(generatedFlags)
	updateCmd.AddCommand(updatePkgCmd)

	generatedFlags = &pflag.FlagSet{}
	err = generateFlags(&service.ProjectNode{}, generatedFlags)
	if err != nil {
		log.Fatalf("%v", err)
	}
	updateProjCmd.Flags().AddFlagSet(generatedFlags)
	updateCmd.AddCommand(updateProjCmd)
}

func updateNode(nodeIdent string) error {
	var err error
	currentNode, err = ParseNodeID(nodeIdent)
	if err != nil {
		return err
	}
	switch cNode := currentNode.(type) {
	case *service.FileNode:
		currentNode, err = getUniqueFileNode(cNode)
		if err != nil {
			return err
		}
		// set fields of node according to flags
		cmdFlags.Visit(visitNodeFlag)

		err = sendFileNode(currentNode.(*service.FileNode))
		if err != nil {
			return fmt.Errorf("Failed sending FileNode: %v", err)
		}
	case *service.PackageNode:
		currentNode, err := getUniquePackageNode(cNode)
		if err != nil {
			return fmt.Errorf("Failed to get package node: %v", err)
		}
		// set fields of node according to flags
		cmdFlags.Visit(visitNodeFlag)
		res, err := buildServiceClient.CreatePackage(context.Background(), currentNode)
		if err != nil {
			return err
		}
		if !res.Success {
			return errors.New("sending package failed")
		}
	case *service.ProjectNode:
		currentNode, err = buildServiceClient.GetProjectNode(context.Background(), &service.ProjectNode{})
		if err != nil {
			return fmt.Errorf("Failed to get project node: %v", err)
		}
		// set fields of node according to flags
		cmdFlags.Visit(visitNodeFlag)
		res, err := buildServiceClient.CreateProject(context.Background(), currentNode.(*service.ProjectNode))
		if err != nil {
			return err
		}
		if !res.Success {
			return errors.New("sending project failed")
		}
	}
	return nil
}

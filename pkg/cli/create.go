package cli

import (
	"errors"
	"log"
	"reflect"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/net/context"
)

var currentNode interface{}
var cmdFlags *pflag.FlagSet

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create a new node",
	Long:  "create a new node described by a node identifier",
}

var createFileCmd = &cobra.Command{
	Use:   "file",
	Short: "create a new file node",
	Long:  "create a new file node described by a node identifier",
	Run:   create,
}

var createPkgCmd = &cobra.Command{
	Use:   "package",
	Short: "create a new package node",
	Long:  "create a new package node described by a node identifier",
	Run:   create,
}

var createProjCmd = &cobra.Command{
	Use:   "project",
	Short: "create a new project node",
	Long:  "create a new project node described by a node identifier",
	Run:   create,
}

func create(cmd *cobra.Command, args []string) {
	setUpBuildService()
	cmdFlags = cmd.Flags()
	err := createNode(args[0], true)
	if err != nil {
		log.Fatalf("Failed to create node: %v", err)
	}
	tearDownServer()
}

func init() {
	var err error
	generatedFlags := &pflag.FlagSet{}
	rootCmd.AddCommand(createCmd)

	err = generateFlags(&service.FileNode{}, generatedFlags)
	if err != nil {
		log.Fatalf("%v", err)
	}
	createFileCmd.Flags().AddFlagSet(generatedFlags)
	createCmd.AddCommand(createFileCmd)

	generatedFlags = &pflag.FlagSet{}
	err = generateFlags(&service.PackageNode{}, generatedFlags)
	if err != nil {
		log.Fatalf("%v", err)
	}
	createPkgCmd.Flags().AddFlagSet(generatedFlags)
	createCmd.AddCommand(createPkgCmd)

	generatedFlags = &pflag.FlagSet{}
	err = generateFlags(&service.ProjectNode{}, generatedFlags)
	if err != nil {
		log.Fatalf("%v", err)
	}
	createProjCmd.Flags().AddFlagSet(generatedFlags)
	createCmd.AddCommand(createProjCmd)
}

func createNode(nodeIdent string, send bool) error {
	var err error
	currentNode, err = ParseNodeID(nodeIdent)
	if err != nil {
		return err
	}

	// set fields according to flags
	cmdFlags.Visit(visitNodeFlag)

	switch reflect.TypeOf(currentNode) {
	case reflect.TypeOf((*service.FileNode)(nil)):
		if send {
			stream, err := buildServiceClient.Build(context.Background())
			if err != nil {
				return err
			}
			stream.Send(currentNode.(*service.FileNode))
			br, err := stream.CloseAndRecv()
			if err != nil {
				return err
			}
			if !br.Success {
				return errors.New("sending filenode failed")
			}
			return nil
		}
	case reflect.TypeOf((*service.PackageNode)(nil)):
		if send {
			br, err := buildServiceClient.CreatePackage(context.Background(), currentNode.(*service.PackageNode))
			if err != nil {
				return err
			}
			if !br.Success {
				return errors.New("sending package failed")
			}
			return nil
		}
	case reflect.TypeOf((*service.ProjectNode)(nil)):
		if send {
			br, err := buildServiceClient.CreateProject(context.Background(), currentNode.(*service.ProjectNode))
			if err != nil {
				return err
			}
			if !br.Success {
				return errors.New("sending project failed")
			}
			return nil
		}
	}
	return nil
}

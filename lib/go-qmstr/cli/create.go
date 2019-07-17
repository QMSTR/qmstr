package cli

import (
	"errors"
	"log"
	"reflect"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"golang.org/x/net/context"
)

var currentNode interface{}
var cmdFlags *pflag.FlagSet

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new node",
	Long:  "create a new node described by a node identifier",
}

var createPathCmd = &cobra.Command{
	Use:   "pathInfo",
	Short: "Create a new path",
	Long:  "create a new path for a file node",
	Run:   create,
}

var createFileCmd = &cobra.Command{
	Use:   "file",
	Short: "Create a new file node",
	Long:  "create a new file node described by a node identifier",
	Run:   create,
}

var createPkgCmd = &cobra.Command{
	Use:   "package",
	Short: "Create a new package node",
	Long:  "create a new package node described by a node identifier",
	Run:   create,
}

var createProjCmd = &cobra.Command{
	Use:   "project",
	Short: "Create a new project node",
	Long:  "create a new project node described by a node identifier",
	Run:   create,
}

func create(cmd *cobra.Command, args []string) {
	setUpBuildService()
	setUpControlService()
	awaitServer()
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

	err = generateFlags(&service.PathInfo{}, generatedFlags)
	if err != nil {
		log.Fatalf("%v", err)
	}
	createPathCmd.Flags().AddFlagSet(generatedFlags)
	createCmd.AddCommand(createPathCmd)

	generatedFlags = &pflag.FlagSet{}
	err = generateFlags(&service.FileNode{}, generatedFlags)
	if err != nil {
		log.Fatalf("%v", err)
	}
	createFileCmd.Flags().AddFlagSet(generatedFlags)
	createFileCmd.Flags().String("path", "", "Set FileNode's path")
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
	case reflect.TypeOf((*service.PathInfo)(nil)):
		if send {
			br, err := buildServiceClient.CreatePathInfo(context.Background(), currentNode.(*service.PathInfo))
			if err != nil {
				return err
			}
			if !br.Success {
				return errors.New("sending path info failed")
			}
			return nil
		}
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
			pkgNode := currentNode.(*service.PackageNode)
			if pkgNode.Version == "" {
				pkgNode.Version = "default"
			}
			br, err := buildServiceClient.CreatePackage(context.Background(), pkgNode)
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

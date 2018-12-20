package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/gnubuilder"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var (
	target, input []string
	mode          gnubuilder.Mode
)

var nodeCmd = &cobra.Command{
	Use:   "node",
	Short: "Modify nodes in the database",
	Long: `Node is used when we want to modify file nodes in the database.
A child command has to be specified after the 'node' command, to define 
the action to be taken in the database. For instance, to add or dump a node.`,
	Run: func(cmd *cobra.Command, args []string) {},
}

var nodeAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add file nodes into the database",
	Long:  `Add file nodes and the connection between them into the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpBuildServer()
		if err := addFileNodes(args); err != nil {
			Log.Fatalf("Add file nodes failed: %v", err)
		}
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.AddCommand(nodeAddCmd)
	nodeAddCmd.Flags().StringSliceVarP(&target, "output", "o", []string{}, "Output file name")
}

func addFileNodes(args []string) error {
	workdir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Could not get current working dir. %v", err)
	}
	if len(args) < 1 {
		return fmt.Errorf("requires at least one input file")
	}
	if len(target) < 1 {
		return fmt.Errorf("requires at least one output file")
	}
	// create file nodes
	fileNodes, err := createFileNodes(args, workdir)
	if err != nil {
		return fmt.Errorf("Failed to create file nodes: %v", err)
	}
	// add file nodes in the database
	stream, err := buildServiceClient.Build(context.Background())
	defer func() {
		res, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("Failed to close the filenode stream: %v", err)
		}
		if !res.Success {
			log.Fatalln("Server filenode stream failed")
		}
	}()
	if err != nil {
		return fmt.Errorf("could not greet: %v", err)
	}
	for _, fileNode := range fileNodes {
		if err := stream.Send(fileNode); err != nil {
			return fmt.Errorf("Failed to send filenode to server")
		}
	}
	return nil
}

func createFileNodes(args []string, workdir string) ([]*service.FileNode, error) {
	fileNodes := []*service.FileNode{}
	dependencies := []*service.FileNode{}
	fmt.Printf("The input is: %v \nThe output is: %v \n", args, target)
	for idx, input := range args {
		inputFileType := builder.CheckInputFileExt(input)
		inputFileNode := builder.NewFileNode(common.BuildCleanPath(workdir, input, false), inputFileType)
		if len(target) > 1 {
			if len(target) != len(args) {
				return nil, fmt.Errorf("Please provide same number of input and output files")
			}
			targetFileType := builder.CheckInputFileExt(target[idx])
			targetFileNode := builder.NewFileNode(common.BuildCleanPath(workdir, target[idx], false), targetFileType)
			targetFileNode.DerivedFrom = []*service.FileNode{inputFileNode}
			fileNodes = append(fileNodes, targetFileNode)
		} else {
			// when we want to connect multiple inputs to a target
			mode = gnubuilder.ModeLink
			dependencies = append(dependencies, inputFileNode)
		}
	}
	if mode == gnubuilder.ModeLink {
		targetFileType := builder.CheckInputFileExt(target[0])
		targetFileNode := builder.NewFileNode(common.BuildCleanPath(workdir, target[0], false), targetFileType)
		targetFileNode.DerivedFrom = dependencies
		fileNodes = append(fileNodes, targetFileNode)
	}
	return fileNodes, nil
}

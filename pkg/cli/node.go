package cli

import (
	"github.com/spf13/cobra"
)

var target, input []string
var phase string

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
		setUpServer()
		
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(nodeCmd)
	nodeCmd.AddCommand(addCmd)
	addCmd.Flags().StringArrayVarP(&target, "out", "o", []string{}, "Output file name")
}


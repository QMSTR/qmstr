package cli

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use:   "delete [type_of_node:attribute:value]",
	Short: "Delete node from the database",
	Long:  `Delete the provided node from the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpServer()
		if err := deleteNode(args); err != nil {
			Log.Fatalf("Delete failed: %v", err)
		}
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}

func deleteNode(args []string) error {
	return nil
}

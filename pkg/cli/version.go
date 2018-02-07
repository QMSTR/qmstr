package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var major = 0
var minor = 1

// quitCmd represents the quit command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Version of qmstr-cli",
	Long:  `The current version of qmstr-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("The qmstr-cli version: %d.%d\n", major, minor)

	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

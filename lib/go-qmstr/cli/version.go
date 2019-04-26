package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var major = 0
var minor = 3

// quitCmd represents the quit command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of qmstrctl",
	Long:  `Print the version of qmstrctl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("This is qmstrctl version %d.%d.\n", major, minor)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

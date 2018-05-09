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
	Short: "version of qmstrctl",
	Long:  `prints the version of qmstr-cli`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("This is qmstrctl version %d.%d.\n", major, minor)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

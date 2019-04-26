package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// QmstrVersion holds the qmstr version
// The variable is initialiazed in the makefile
var QmstrVersion string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of qmstrctl",
	Long:  `Print the version of qmstrctl`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("This is qmstrctl version %s\n", QmstrVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

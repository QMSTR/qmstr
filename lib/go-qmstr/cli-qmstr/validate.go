package cliqmstr

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a package against a manifest",
	Long: `validate needs at least one argument, the package file. Depending on the package format the manifest can be internal and/or external.
For example: qmstr validate curl-a.b.c.deb curl-a.b.c.deb.spdx`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("validate stub command called.")
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

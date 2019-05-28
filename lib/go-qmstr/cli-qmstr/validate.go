package cliqmstr

import (
	"fmt"
	"log"

	"github.com/QMSTR/qmstr/lib/go-qmstr/validation"
	"github.com/QMSTR/qmstr/modules/manifests"
	"github.com/QMSTR/qmstr/modules/packages"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate a package against a manifest",
	Long: `validate needs at least one argument, the package file. Depending on the package format the manifest can be internal and/or external.
For example: qmstr validate curl-a.b.c.deb curl-a.b.c.deb.spdx`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		pkg, err := packages.PackageFromFile(args[0])
		if err != nil {
			log.Fatalf("Error getting pkg from file %q: %v\n", args[0], err)
		}
		var mani validation.Manifest
		if len(args) == 2 {
			mani, err = manifests.ManifestFromFile(args[1])
			if err != nil {
				log.Fatalf("Error getting manifest from file %q: %v\n", args[1], err)
			}
		}
		err = pkg.Validate(mani)
		if err != nil {
			log.Fatalf("Validation fail: %v", err)
		}
		fmt.Println("Validation successful!")
		return
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

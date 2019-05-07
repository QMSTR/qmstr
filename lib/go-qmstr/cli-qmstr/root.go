package cliqmstr

import (
	"io/ioutil"
	golog "log" // avoid having "log" in the namespace
	"os"

	"github.com/spf13/cobra"
)

var (
	// Enable verbose log output using the Debug logger
	verbose bool
	// Debug receives log messages in verbose mode
	Debug *golog.Logger
	// Log is the standard logger
	Log *golog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qmstr",
	Short: "qmstr performs Quartermaster client-side operations",
	Long: `qmstr performs Quartermaster operations that have only client-side effects.
It does not require a QMSTR master to be available.`,
}

func init() {
	cobra.OnInitialize(initLogging)
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable diagnostics")
}

func initLogging() {
	if verbose {
		Debug = golog.New(os.Stderr, "DEBUG: ", golog.Ldate|golog.Ltime)
	} else {
		Debug = golog.New(ioutil.Discard, "", golog.Ldate|golog.Ltime)
	}
	Log = golog.New(os.Stderr, "", golog.Ldate|golog.Ltime)
}

// Execute the control program and perform the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		golog.Printf("Error: %v!\n", err)
		os.Exit(1)
	}
}

package cli

import (
	"errors"
	"fmt"
	golog "log" // avoid having "log" in the namespace
	"os"

	"github.com/QMSTR/qmstr/pkg/logging"
	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

// global variables
var (
	conn                 *grpc.ClientConn
	controlServiceClient service.ControlServiceClient
	buildServiceClient   service.BuildServiceClient
	// AddressOptional means the command does not require a server address (version, start, ...)
	AddressOptional bool
	address         string
	verbose         bool
	// Debug receives log messages in verbose mode
	Debug *golog.Logger
	// Log is the standard logger
	Log *golog.Logger
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qmstrctl",
	Short: "qmstrctl controls and manages the Quartermaster master",
	Long: `qmstrctl controls and manages the Quartermaster master process.
	It provides commands to run, quit and configure the master.`,
	Run:              func(cmd *cobra.Command, args []string) {},
	PersistentPreRun: SetupLogging,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&address, "cserv", "", "connect to control service")
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "enable diagnostics")
}

// Execute the control program and perform the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		golog.Printf("Error: %v!\n", err)
		os.Exit(1)
	}
}

func setupConnection() error {
	if conn != nil {
		return nil
	}
	if len(address) == 0 {
		address = os.Getenv("QMSTR_MASTER")
	}
	if len(address) == 0 {
		return errors.New("Error: master address not specified")
	}
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return fmt.Errorf("Failed to setup connection to master: %v", err)
	}
	return nil
}

// Set up connection to the control service
func setUpControlService() {
	if err := setupConnection(); err != nil {
		Log.Fatalf("Failed to setup control service connection: %v", err)
	}
	controlServiceClient = service.NewControlServiceClient(conn)
	Debug.Printf("Connection to control service at %v set up\n", address)
}

// Set up connection to the build service
func setUpBuildService() {
	if err := setupConnection(); err != nil {
		Log.Fatalf("Failed to setup build service connection: %v", err)
	}
	buildServiceClient = service.NewBuildServiceClient(conn)
	Debug.Printf("Connection to build service at %v set up\n", address)
}

func tearDownServer() {
	conn.Close()
}

// SetupLogging sets up logging
func SetupLogging(cmd *cobra.Command, args []string) {
	log := logging.Setup(verbose)
	Debug = log.Debug
	Log = log.Log
}

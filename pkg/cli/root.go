package cli

import (
	"fmt"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var conn *grpc.ClientConn
var controlServiceClient service.ControlServiceClient
var address string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "qmstr-cli",
	Short: "Qmstr client connects to qmstr and serves",
	Long: `Qmstr client connects to qmstr and serves depending 
	what we want to do with it. Right now it just quits the server 
	and prints the version of qmstr-cli.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func init() {
	rootCmd.PersistentFlags().StringVar(&address, "cserv", "localhost:50051", "connect to control service")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setUpServer() {
	// Set up server connection
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure())
	fmt.Printf("Connecting to address: %v\n", address)
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	controlServiceClient = service.NewControlServiceClient(conn)
}

func tearDownServer() {
	conn.Close()
}

package cli

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/spf13/cobra"
)

var force bool

// quitCmd represents the quit command
var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: "Quit qmstr",
	Long:  `Run quit if you want to quit qmstr.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpServer()
		quitServer()
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(quitCmd)
	quitCmd.Flags().BoolVarP(&force, "force", "f", false, "force quit")
}

func quitServer() {
	resp, err := buildServiceClient.Quit(context.Background(), &pb.QuitMessage{Kill: force})
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	if !resp.Success {
		fmt.Println("Server responded unsuccessful")
		os.Exit(ReturnCodeResponseFalseError)
	}
}

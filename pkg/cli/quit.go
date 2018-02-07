package cli

import (
	"fmt"

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
		fmt.Println("quit called")

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
	buildServiceClient.Quit(context.Background(), &pb.QuitMessage{Kill: force})
}

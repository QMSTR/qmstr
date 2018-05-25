package cli

import (
	"os"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "prints server status",
	Long:  `prints the server status`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpServer()
		getStatus()
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func getStatus() {
	res, err := controlServiceClient.Status(context.Background(), &service.StatusMessage{Phase: true, Switch: true})
	if err != nil {
		Log.Println("Master server not yet available")
		os.Exit(ReturnCodeServerCommunicationError)
	}
	switching := "not"
	if res.Switching {
		switching = ""
	}
	Log.Printf("Master server is in %s phase and %s switching", res.Phase, switching)
}

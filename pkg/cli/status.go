package cli

import (
	"io"
	"os"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var follow bool

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "prints server status",
	Long:  `prints the server status`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpControlService()
		getStatus()
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVar(&follow, "follow", false, "print status and follow events")
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

	if res.Error != "" {
		Log.Printf("Failure caused by %s\n", res.Error)
	}

	if follow {
		err = printEvents()
		if err != nil {
			Log.Printf("Failed to follow event stream")
			os.Exit(ReturnCodeServerCommunicationError)
		}
	}
}

func printEvents() error {
	stream, err := controlServiceClient.SubscribeEvents(context.Background(), &service.EventMessage{Class: service.EventClass_ALL})
	if err != nil {
		Log.Printf("Could not subscribe to events %v", err)
		return err
	}
	for {
		event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			Log.Printf("failed to receive event %v", err)
			return err
		}
		Log.Printf("Event: %v, Message: %v", event.Class, event.Message)
	}
	return nil
}

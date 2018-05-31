package cli

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
)

var anaCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Start analysis on qmstr-master",
	Long:  `Start analysis phase on the master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpServer()
		startPhase(master.PhaseIDAnalysis)
		tearDownServer()
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Start report on qmstr-master",
	Long:  `Start report phase on the master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpServer()
		startPhase(master.PhaseIDReport)
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(anaCmd)
	rootCmd.AddCommand(reportCmd)
}

func startPhase(phase int32) {
	if verbose {
		go func() {
			stream, err := controlServiceClient.SubscribeEvents(context.Background(), &service.EventMessage{Class: string(master.EventAll)})
			if err != nil {
				Log.Printf("Could not subscribe to events %v", err)
			}
			for {
				event, err := stream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					Log.Printf("failed to receive event %v", err)
					continue
				}
				Log.Printf("Event: %v, Message: %v", event.Class, event.Message)
			}
		}()
	}
	resp, err := controlServiceClient.SwitchPhase(context.Background(), &service.SwitchPhaseMessage{Phase: phase})
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	if !resp.Success {
		fmt.Println("Server responded unsuccessful")
		os.Exit(ReturnCodeResponseFalseError)
	}
}

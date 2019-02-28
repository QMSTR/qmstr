package cli

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
)

var anaCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Start analysis on qmstr-master",
	Long:  `Start analysis phase on the master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpControlService()
		startPhase(service.Phase_ANALYSIS)
		tearDownServer()
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Start report on qmstr-master",
	Long:  `Start report phase on the master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpControlService()
		startPhase(service.Phase_REPORT)
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(anaCmd)
	rootCmd.AddCommand(reportCmd)
}

func startPhase(phase service.Phase) {
	if verbose {
		go printEvents()
	}
	resp, err := controlServiceClient.SwitchPhase(context.Background(), &service.SwitchPhaseMessage{Phase: phase})
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	if !resp.Success {
		fmt.Printf("Server reported failure:\n%s\n", resp.Error)
		os.Exit(ReturnCodeResponseFalseError)
	}
}

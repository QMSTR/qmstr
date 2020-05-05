package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/lib/go-qmstr/config"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spf13/cobra"
)

var anaCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Start analysis on the QMSTR master",
	Long:  `Start analysis phase on the QMSTR master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpControlService()
		startPhase(service.Phase_ANALYSIS)
		tearDownServer()
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Start report on the QMSTR master",
	Long:  `Start report phase on the QMSTR master server.`,
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
	// switch master server phase
	resp, err := controlServiceClient.SwitchPhase(context.Background(), &service.SwitchPhaseMessage{Phase: phase})
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	// export master config from the server
	var config config.MasterConfig
	err = json.Unmarshal([]byte(resp.MasterConfig), &config)
	if err != nil {
		log.Fatal(err)
	}

	if !resp.Success {
		fmt.Printf("Server reported failure:\n%s\n", resp.Error)
		os.Exit(ReturnCodeResponseFalseError)
	}
}

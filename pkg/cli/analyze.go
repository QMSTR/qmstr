package cli

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/spf13/cobra"
)

var anaCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Start analysis on qmstr-master",
	Long:  `Start analysis described in provided YAML file on the master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpServer()
		startAnalysis()
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(anaCmd)
}

func startAnalysis() {
	resp, err := buildServiceClient.Analyze(context.Background(), &buildservice.AnalysisMessage{Async: false})
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	if !resp.Success {
		fmt.Println("Server responded unsuccessful")
		os.Exit(ReturnCodeResponseFalseError)
	}
}

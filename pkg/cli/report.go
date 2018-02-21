package cli

import (
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/spf13/cobra"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Start report generation on qmstr-master",
	Long:  `Generate report described in provided YAML file on the master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpServer()
		submitReportRequest()
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

func submitReportRequest() {
	bsrStream, err := buildServiceClient.Report(context.Background(), &buildservice.ReportMessage{Async: false})
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}

	for {
		resp, err := bsrStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Report failed %v", err)
			os.Exit(ReturnCodeServerCommunicationError)
		}
		fmt.Printf("[%s] %s\n", resp.Reporter, resp.ResponseMessage)
	}
}

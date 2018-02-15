package cli

import (
	"fmt"
	"os"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Start report generation on qmstr-master",
	Long:  `Generate report described in provided YAML file on the master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("Please provide one yaml file")
			os.Exit(ReturnCodeParameterError)
		}

		repMsg := &buildservice.ReportMessage{}

		setUpServer()
		unmarshalReportRequest(args[0], repMsg)
		submitReportRequest(repMsg)
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(reportCmd)
}

func submitReportRequest(repMsg *buildservice.ReportMessage) {
	resp, err := buildServiceClient.Report(context.Background(), repMsg)
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	if !resp.Success {
		fmt.Println("Server responded unsuccessful")
		os.Exit(ReturnCodeResponseFalseError)
	}
	fmt.Println(resp.ResponseMessage)
}

func unmarshalReportRequest(yamlFile string, request *buildservice.ReportMessage) {
	data, err := consumeFile(yamlFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(ReturnCodeSysError)
	}

	err = yaml.Unmarshal(data, &request)
	if err != nil {
		fmt.Printf("Wrong YAML file format %v", err)
		os.Exit(ReturnCodeFormatError)
	}
}

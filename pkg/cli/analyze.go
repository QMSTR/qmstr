package cli

import (
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var anaCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Start analysis on qmstr-master",
	Long:  `Start analysis described in provided YAML file on the master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 1 {
			fmt.Println("Please provide one yaml file")
			os.Exit(ReturnCodeParameterError)
		}

		anaMsg := &buildservice.AnalysisMessage{}

		setUpServer()
		unmarshalAnalysisRequest(args[0], anaMsg)
		submitAnalysis(anaMsg)
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(anaCmd)
}

func submitAnalysis(anaMsg *buildservice.AnalysisMessage) {
	resp, err := buildServiceClient.Analyze(context.Background(), anaMsg)
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	if !resp.Success {
		fmt.Println("Server responded unsuccessful")
		os.Exit(ReturnCodeResponseFalseError)
	}
}

func unmarshalAnalysisRequest(yamlFile string, request *buildservice.AnalysisMessage) {
	f, err := os.Open(yamlFile)
	if err != nil {
		fmt.Printf("Failed to open file %v", err)
		os.Exit(ReturnCodeSysError)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Printf("Failed to read from file %v", err)
		os.Exit(ReturnCodeSysError)
	}

	err = yaml.Unmarshal(data, &request)
	if err != nil {
		fmt.Printf("Wrong YAML file format %v", err)
		os.Exit(ReturnCodeFormatError)
	}
}

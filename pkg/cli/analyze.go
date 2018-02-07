package cli

import (
	"io/ioutil"
	"log"
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
			log.Fatalln("Please provide one yaml file")
		}

		anaMsg := &buildservice.AnalysisMessage{}

		setUpServer()
		unmarshalAnalysisRequest(args[0], anaMsg)
		log.Printf("About to submit %v", anaMsg)
		submitAnalysis(anaMsg)
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(anaCmd)
}

func submitAnalysis(anaMsg *buildservice.AnalysisMessage) {
	buildServiceClient.Analyze(context.Background(), anaMsg)
}

func unmarshalAnalysisRequest(yamlFile string, request *buildservice.AnalysisMessage) {
	f, err := os.Open(yamlFile)
	if err != nil {
		log.Fatalf("Failed to open file %v", err)
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("Failed to read from file %v", err)
	}

	err = yaml.Unmarshal(data, &request)
	if err != nil {
		log.Fatalf("Wrong YAML file format %v", err)
	}
}

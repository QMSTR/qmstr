package cli

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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
		masterConfig := startPhase(service.Phase_ANALYSIS)
		startAnalyzerModules(masterConfig)
		tearDownServer()
	},
}

var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Start report on the QMSTR master",
	Long:  `Start report phase on the QMSTR master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpControlService()
		masterConfig := startPhase(service.Phase_REPORT)
		startReporterModules(masterConfig)
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(anaCmd)
	rootCmd.AddCommand(reportCmd)
}

// start analysis/reporting phase and return the master config
func startPhase(phase service.Phase) config.MasterConfig {
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
	return config
}

// start the analyzer modules configured in the qmstr.yaml
func startAnalyzerModules(masterConfig config.MasterConfig) {
	//loop through the analyzers in master config
	for idx, anaConfig := range masterConfig.Analysis {
		analyzerName := anaConfig.Analyzer
		// Initialize analyzer
		_, err := controlServiceClient.InitModule(context.Background(), &service.InitModuleRequest{ModuleName: analyzerName, ExtraConfig: anaConfig.TrustLevel})
		if err != nil {
			fmt.Printf("Failed initializing module %s: %v\n", analyzerName, err)
			os.Exit(ReturnCodeServerCommunicationError)
		}

		// Run analyzer module
		cmd := exec.Command(analyzerName, "--aserv", masterConfig.Server.RPCAddress, "--aid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logModuleError(analyzerName, out)
			errMsg := fmt.Sprintf("Analyzer %s failed: %v", analyzerName, err)
			controlServiceClient.ShutdownModule(context.Background(), &service.ShutdownModuleRequest{Message: errMsg, DB: true})
			os.Exit(ReturnCodeServerCommunicationError)
		}
		msg := fmt.Sprintf("Analyzer %s successfully finished", analyzerName)
		if _, err := controlServiceClient.ShutdownModule(context.Background(), &service.ShutdownModuleRequest{Message: msg, DB: true}); err != nil {
			fmt.Printf("Failed shutting down the module %s: %v\n", analyzerName, err)
			os.Exit(ReturnCodeServerCommunicationError)
		}
		log.Printf("Analyzer %s finished successfully:\n%s\n", analyzerName, out)
	}
	// Now we need to notify the server that
	// the analysis phase has finished its tasks
	// and each ready to move on to the next phase
	controlServiceClient.ClosePhase(context.Background(), &service.Request{})
}

// start the reporter modules configured in the qmstr.yaml
func startReporterModules(masterConfig config.MasterConfig) {
	//loop through the reporters in master config
	for idx, reporterConfig := range masterConfig.Reporting {
		reporterName := reporterConfig.Reporter
		// Initialize reporter
		_, err := controlServiceClient.InitModule(context.Background(), &service.InitModuleRequest{ModuleName: reporterName})
		if err != nil {
			fmt.Printf("Failed initializing module %s: %v\n", reporterName, err)
			os.Exit(ReturnCodeServerCommunicationError)
		}

		// Run reporter module
		cmd := exec.Command(reporterName, "--rserv", masterConfig.Server.RPCAddress, "--rid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logModuleError(reporterName, out)
			errMsg := fmt.Sprintf("Reporter %s failed: %v", reporterName, err)
			controlServiceClient.ShutdownModule(context.Background(), &service.ShutdownModuleRequest{Message: errMsg, DB: false})
			os.Exit(ReturnCodeServerCommunicationError)
		}
		msg := fmt.Sprintf("Reporter %s successfully finished", reporterName)
		controlServiceClient.ShutdownModule(context.Background(), &service.ShutdownModuleRequest{Message: msg, DB: false})
		log.Printf("Reporter %s finished successfully:\n%s\n", reporterName, out)
	}
	// No need to call ClosePhase() here.
	// Once the reporters have finished we can quit the master server
}

func logModuleError(moduleName string, output []byte) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s failed with:\n", moduleName))
	s := bufio.NewScanner(strings.NewReader(string(output)))
	for s.Scan() {
		buffer.WriteString(fmt.Sprintf("\t--> %s\n", s.Text()))
	}
	log.Println(buffer.String())
}

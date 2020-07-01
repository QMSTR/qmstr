package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/lib/go-qmstr/cli"
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
		close(cli.PingAnalyzer) // Ping modules to start!
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

		// PING MODULES TO START

		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(anaCmd)
	rootCmd.AddCommand(reportCmd)
}

// start analysis/reporting phase
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
	if !resp.Success {
		fmt.Printf("Server reported failure:\n%s\n", resp.Error)
		os.Exit(ReturnCodeResponseFalseError)
	}
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

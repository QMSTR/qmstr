package cli

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spf13/cobra"
)

// ModulesAreDone is a signal channel to close
// analysis/reporting phase once modules are done
var ModulesAreDone chan struct{}

var anaCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Start analysis on the QMSTR master",
	Long:  `Start analysis phase on the QMSTR master server.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpControlService()
		startPhase(service.Phase_ANALYSIS)
		startAnalyzers()
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
		startReporters()
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

func startAnalyzers() {
	close(PingAnalyzer) // Ping modules to start!
	waitModulesToFinish()

	// Close analysis phase
	controlServiceClient.ClosePhase(context.Background(), &service.Request{})
}

func startReporters() {
	close(PingReporter) // Ping modules to start!
	waitModulesToFinish()
	// No need to close phase
}

func waitModulesToFinish() {
	// wait until all modules have finished
	ModulesAreDone = make(chan struct{})
	log.Printf("Waiting for modules to finish.. \n")
	<-ModulesAreDone // <-- THIS MAY NOT WORK!!! Select{}
	log.Printf("All modules have finished! \n")
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

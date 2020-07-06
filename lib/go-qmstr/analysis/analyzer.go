package analysis

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	goflag "flag"

	"github.com/QMSTR/qmstr/lib/go-qmstr/cli"
	"github.com/QMSTR/qmstr/lib/go-qmstr/config"
	"github.com/QMSTR/qmstr/lib/go-qmstr/module"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	flag "github.com/spf13/pflag"
)

type Analyzer struct {
	module.MasterClient
	module AnalyzerModule
	id     int32
	name   string
}

type AnalyzerModule interface {
	Configure(configMap map[string]string) error
	Analyze(masterClient *module.MasterClient, token int64) error
	PostAnalyze() error
}

// CountAnalyzers counts the analyzers that run in the current process
var CountAnalyzers int32

func NewAnalyzer(anaModule AnalyzerModule) *Analyzer {
	var serviceAddress string
	var anaID int32
	CountAnalyzers++
	// TODO: Connect to QMSTRADDRESS
	flag.StringVar(&serviceAddress, "aserv", "", "connect to analyzer service address")
	flag.Int32Var(&anaID, "aid", -1, "unique analyzer id")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	mc := module.NewMasterClient(serviceAddress)

	return &Analyzer{MasterClient: mc, id: anaID, module: anaModule}
}

func (a *Analyzer) GetModuleName() string {
	return a.name
}

func (a *Analyzer) RunAnalyzerModule() error {
	configResp, err := a.AnaSvcClient.GetAnalyzerConfig(context.Background(), &service.AnalyzerConfigRequest{AnalyzerID: a.id})
	if err != nil {
		log.Printf("Could not get configuration %v\n", err)
		return fmt.Errorf("could not get analyzer configuration %v", err)
	}
	var analyzerConfig *config.Analysis
	err = json.Unmarshal([]byte(configResp.AnalyzerConfig), &analyzerConfig)
	if err != nil {
		return fmt.Errorf("fail unmarshalling analyzer config %v", err)
	}
	a.name = analyzerConfig.Analyzer
	cacheDir := analyzerConfig.Config["cachedir"]
	outDir := analyzerConfig.Config["outputdir"]

	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create cache directory for module %s %v", a.GetModuleName(), err)
	}

	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory for module %s %v", a.GetModuleName(), err)
	}

	// Initialize analyzer
	_, err = a.CtrlSvcClient.InitModule(context.Background(), &service.InitModuleRequest{
		ModuleName: a.name, ExtraConfig: analyzerConfig.TrustLevel})
	if err != nil {
		return fmt.Errorf("%v: %v", err, cli.ReturnCodeServerCommunicationError)
	}

	err = a.module.Configure(analyzerConfig.Config)
	if err != nil {
		return fmt.Errorf("failed to configure analyzer module %s %v", a.GetModuleName(), err)
	}

	err = a.module.Analyze(&a.MasterClient, configResp.Token)
	if err != nil {
		return fmt.Errorf("Analysis failed for analyzer module %s %v", a.GetModuleName(), err)
	}

	msg := fmt.Sprintf("Analyzer %s finished successfully ", a.name)
	log.Println(msg)
	// Ping master server that the analyzer finished
	a.CtrlSvcClient.ShutdownModule(context.Background(), &service.ShutdownModuleRequest{Message: msg, DB: true})

	return nil
}

// ReduceAnalyzersCounter is called everytime an analyzer finishes its process.
// When it reaches 0, it sends a signal to close the analysis phase
func ReduceAnalyzersCounter() {
	CountAnalyzers--
	if CountAnalyzers == 0 { // all analyzers have finished
		// close analysis phase
		close(cli.ModulesAreDone)
	}
	if CountAnalyzers < 0 {
		log.Printf("WARNING: Analyzers count cannot be minus")
	}
}

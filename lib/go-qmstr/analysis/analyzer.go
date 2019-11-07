package analysis

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	goflag "flag"

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

func NewAnalyzer(anaModule AnalyzerModule) *Analyzer {
	var serviceAddress string
	var anaID int32
	flag.StringVar(&serviceAddress, "aserv", "localhost:50051", "Analyzer service address")
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

	a.name = configResp.Name
	cacheDir := configResp.ConfigMap["cachedir"]
	outDir := configResp.ConfigMap["outputdir"]

	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create cache directory for module %s %v", a.GetModuleName(), err)
	}

	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory for module %s %v", a.GetModuleName(), err)
	}

	err = a.module.Configure(configResp.ConfigMap)
	if err != nil {
		return fmt.Errorf("failed to configure analyzer module %s %v", a.GetModuleName(), err)
	}

	err = a.module.Analyze(&a.MasterClient, configResp.Token)
	if err != nil {
		return fmt.Errorf("Analysis failed for analyzer module %s %v", a.GetModuleName(), err)
	}
	return nil
}

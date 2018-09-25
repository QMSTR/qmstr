package analysis

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	goflag "flag"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

type Analyzer struct {
	analysisService service.AnalysisServiceClient
	controlService  service.ControlServiceClient
	module          AnalyzerModule
	id              int32
	name            string
}

type AnalyzerModule interface {
	Configure(configMap map[string]string) error
	Analyze(controlService service.ControlServiceClient, analysisService service.AnalysisServiceClient, token int64, session string) error
	PostAnalyze() error
}

func NewAnalyzer(module AnalyzerModule) *Analyzer {
	var serviceAddress string
	var anaID int32
	flag.StringVar(&serviceAddress, "aserv", "localhost:50051", "Analyzer service address")
	flag.Int32Var(&anaID, "aid", -1, "unique analyzer id")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	anaServiceClient := service.NewAnalysisServiceClient(conn)
	controlServiceClient := service.NewControlServiceClient(conn)

	return &Analyzer{id: anaID, module: module, analysisService: anaServiceClient, controlService: controlServiceClient}
}

func (a *Analyzer) GetModuleName() string {
	return a.name
}

func (a *Analyzer) RunAnalyzerModule() error {
	configResp, err := a.analysisService.GetAnalyzerConfig(context.Background(), &service.AnalyzerConfigRequest{AnalyzerID: a.id})
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

	err = a.module.Analyze(a.controlService, a.analysisService, configResp.Token, configResp.Session)
	if err != nil {
		return fmt.Errorf("Analysis failed for analyzer module %s %v", a.GetModuleName(), err)
	}
	return nil
}

func CreateInfoNode(infoType string, dataNodes ...*service.InfoNode_DataNode) *service.InfoNode {
	return &service.InfoNode{
		Type:      infoType,
		DataNodes: dataNodes,
	}
}

func CreateWarningNode(warning string) *service.InfoNode {
	return CreateInfoNode("warning", &service.InfoNode_DataNode{
		Type: "warning_message",
		Data: warning,
	})
}

func CreateErrorNode(errorMes string) *service.InfoNode {
	return CreateInfoNode("error", &service.InfoNode_DataNode{
		Type: "error_message",
		Data: errorMes,
	})
}

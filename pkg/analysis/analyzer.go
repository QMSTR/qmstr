package analysis

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/service"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

type Analyzer struct {
	analysisService service.AnalysisServiceClient
	module          AnalyzerModule
	id              int32
	name            string
	cacheDir        string
}

type AnalyzerModule interface {
	Configure(configMap map[string]string) error
	Analyze(node *service.FileNode) (*service.InfoNodeSlice, error)
}

func NewAnalyzer(module AnalyzerModule) *Analyzer {
	var serviceAddress string
	var anaID int32
	flag.StringVar(&serviceAddress, "aserv", "localhost:50051", "Analyzer service address")
	flag.Int32Var(&anaID, "aid", -1, "unique analyzer id")
	flag.Parse()

	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	anaServiceClient := service.NewAnalysisServiceClient(conn)

	return &Analyzer{id: anaID, module: module, analysisService: anaServiceClient}
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
	a.cacheDir = configResp.CacheDir

	err = a.module.Configure(configResp.ConfigMap)
	if err != nil {
		return fmt.Errorf("failed to configure analyzer module %s %v", a.GetModuleName(), err)
	}

	nodeResp, err := a.analysisService.GetNodes(context.Background(), &service.NodeRequest{Type: configResp.TypeSelector})
	if err != nil {
		return fmt.Errorf("could not get nodes %v", err)
	}

	resultMap := map[string]*service.InfoNodeSlice{}

	for _, node := range nodeResp.FileNodes {
		for _, substitution := range configResp.PathSub {
			node.Path = strings.Replace(node.Path, substitution.Old, substitution.New, 1)
		}

		infoNodeSlice, err := a.module.Analyze(node)
		if err != nil {
			return fmt.Errorf("analyzer %s failed %v", a.name, err)
		}

		if len(infoNodeSlice.Inodes) > 0 {
			resultMap[node.Hash] = infoNodeSlice
		}
	}

	anaresp, err := a.analysisService.SendNodes(context.Background(), &service.AnalysisMessage{ResultMap: resultMap, Token: configResp.Token})
	if err != nil {
		return fmt.Errorf("failed to send nodes %v", err)
	}
	if !anaresp.Success {
		return errors.New("Server could not process nodes")
	}

	return nil
}

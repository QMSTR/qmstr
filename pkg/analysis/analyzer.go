package analysis

import (
	"log"
	"os"
	"strings"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/service"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

type Analyzer struct {
	analysisService service.AnalysisServiceClient
	plugin          AnalyzerPlugin
	id              int32
}

type AnalyzerPlugin interface {
	Configure(configMap map[string]string) error
	Analyze(node *service.FileNode) (*service.InfoNodeSlice, error)
}

func NewAnalyzer(plugin AnalyzerPlugin) *Analyzer {
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

	return &Analyzer{id: anaID, plugin: plugin, analysisService: anaServiceClient}
}

func (a *Analyzer) RunPlugin() {
	configResp, err := a.analysisService.GetConfig(context.Background(), &service.ConfigRequest{AnalyzerID: a.id})
	if err != nil {
		log.Printf("Could not get configuration %v\n", err)
		os.Exit(666)
	}

	a.plugin.Configure(configResp.ConfigMap)

	nodeResp, err := a.analysisService.GetNodes(context.Background(), &service.NodeRequest{Type: configResp.TypeSelector})
	if err != nil {
		log.Printf("Could not get nodes %v\n", err)
		os.Exit(667)
	}

	resultMap := map[string]*service.InfoNodeSlice{}

	for _, node := range nodeResp.FileNodes {
		for _, substitution := range configResp.PathSub {
			node.Path = strings.Replace(node.Path, substitution.Old, substitution.New, 1)
		}

		infoNodeSlice, err := a.plugin.Analyze(node)
		if err != nil {
			log.Printf("Analysis failed %v\n", err)
			os.Exit(670)
		}

		resultMap[node.Hash] = infoNodeSlice
	}

	anaresp, err := a.analysisService.SendNodes(context.Background(), &service.AnalysisMessage{ResultMap: resultMap, Token: configResp.Token})
	if err != nil {
		log.Printf("Failed to send nodes %v\n", err)
		os.Exit(668)
	}
	if !anaresp.Success {
		log.Println("Server could not process nodes")
		os.Exit(669)
	}

}

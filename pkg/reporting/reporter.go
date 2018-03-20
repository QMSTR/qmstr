package reporting

import (
	"io"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/master"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/service"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

type Reporter struct {
	reportingService service.ReportServiceClient
	plugin           ReporterPlugin
	id               int32
	name             string
}

type ReporterPlugin interface {
	Configure(configMap map[string]string) error
	Report(node *service.FileNode) error
	PostReport() error
}

func NewReporter(plugin ReporterPlugin) *Reporter {
	var serviceAddress string
	var anaID int32
	flag.StringVar(&serviceAddress, "rserv", "localhost:50051", "Reporting service address")
	flag.Int32Var(&anaID, "rid", -1, "unique reporter id")
	flag.Parse()

	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	reportingServiceClient := service.NewReportServiceClient(conn)

	return &Reporter{id: anaID, plugin: plugin, reportingService: reportingServiceClient}
}

func (r *Reporter) RunReporterPlugin() {
	configResp, err := r.reportingService.GetReporterConfig(context.Background(), &service.ReporterConfigRequest{ReporterID: r.id})
	if err != nil {
		log.Printf("Could not get configuration %v\n", err)
		os.Exit(master.ReturnReportServiceCommFailed)
	}

	r.plugin.Configure(configResp.ConfigMap)

	respStream, err := r.reportingService.GetReportNodes(context.Background(), &service.ReportRequest{Type: configResp.TypeSelector})
	if err != nil {
		log.Printf("Could not get nodes %v\n", err)
		os.Exit(master.ReturnReportServiceCommFailed)
	}

	for {
		resp, err := respStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Reporter %s failed %v\n", r.name, err)
			os.Exit(master.ReturnReporterFailed)

		}
		err = r.plugin.Report(resp.FileNode)
		if err != nil {
			log.Printf("Reporter %s failed %v\n", r.name, err)
			os.Exit(master.ReturnReporterFailed)
		}
	}

}

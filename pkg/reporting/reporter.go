package reporting

import (
	"fmt"
	"io"
	"log"

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

// RunReporterPlugin is the main driver function for each reporter.
func (r *Reporter) RunReporterPlugin() error {
	configResp, err := r.reportingService.GetReporterConfig(context.Background(), &service.ReporterConfigRequest{ReporterID: r.id})
	if err != nil {
		return fmt.Errorf("could not get configuration %v", err)
	}

	r.plugin.Configure(configResp.ConfigMap)

	respStream, err := r.reportingService.GetReportNodes(context.Background(), &service.ReportRequest{Type: configResp.TypeSelector})
	if err != nil {
		return fmt.Errorf("could not get nodes %v", err)
	}

	for {
		resp, err := respStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reporter %s failed %v", r.name, err)

		}
		err = r.plugin.Report(resp.FileNode)
		if err != nil {
			return fmt.Errorf("reporter %s failed %v", r.name, err)
		}
	}

	if err := r.plugin.PostReport(); err != nil {
		return fmt.Errorf("reporter %s failed inn PostReport: %v", r.name, err)
	}

	return nil
}

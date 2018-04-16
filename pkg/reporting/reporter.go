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
	module           ReporterModule
	id               int32
	name             string
}

type ReporterModule interface {
	Configure(configMap map[string]string) error
	Report(node *service.PackageNode) error
	PostReport() error
}

func NewReporter(name string, module ReporterModule) *Reporter {
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

	return &Reporter{id: anaID, module: module, reportingService: reportingServiceClient, name: name}
}

func (r *Reporter) GetModuleName() string {
	return r.name
}

// RunReporterModule is the main driver function for each reporter.
func (r *Reporter) RunReporterModule() error {
	configResp, err := r.reportingService.GetReporterConfig(context.Background(), &service.ReporterConfigRequest{ReporterID: r.id})
	if err != nil {
		return fmt.Errorf("could not get configuration %v", err)
	}

	r.module.Configure(configResp.ConfigMap)

	respStream, err := r.reportingService.GetReportNodes(context.Background(), &service.ReportRequest{Session: configResp.Session})
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
		err = r.module.Report(resp.PackageNode)
		if err != nil {
			return fmt.Errorf("reporter %s failed %v", r.name, err)
		}
	}

	if err := r.module.PostReport(); err != nil {
		return fmt.Errorf("reporter %s failed in PostReport: %v", r.name, err)
	}

	return nil
}

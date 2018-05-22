package reporting

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	goflag "flag"

	"github.com/QMSTR/qmstr/pkg/service"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

type Reporter struct {
	reportingService service.ReportServiceClient
	controlService   service.ControlServiceClient
	module           ReporterModule
	id               int32
	name             string
}

type ReporterModule interface {
	Configure(configMap map[string]string) error
	Report(node *service.PackageNode) error
	PostReport() error
}

func NewReporter(module ReporterModule) *Reporter {
	var serviceAddress string
	var anaID int32
	flag.StringVar(&serviceAddress, "rserv", "localhost:50051", "Reporting service address")
	flag.Int32Var(&anaID, "rid", -1, "unique reporter id")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	reportingServiceClient := service.NewReportServiceClient(conn)
	controlServiceClient := service.NewControlServiceClient(conn)

	return &Reporter{id: anaID, module: module, reportingService: reportingServiceClient, controlService: controlServiceClient}
}

func (r *Reporter) GetModuleName() string {
	return r.name
}

// RunReporterModule is the main driver function for each reporter.
func (r *Reporter) RunReporterModule() error {
	configResp, err := r.reportingService.GetReporterConfig(context.Background(), &service.ReporterConfigRequest{ReporterID: r.id})
	if err != nil {
		return fmt.Errorf("could not get reporter configuration %v", err)
	}

	// Set module name
	r.name = configResp.Name
	cacheDir := configResp.ConfigMap["cachedir"]
	outDir := configResp.ConfigMap["outputdir"]

	err = os.MkdirAll(cacheDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create cache directory for module %s %v", r.GetModuleName(), err)
	}

	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory for module %s %v", r.GetModuleName(), err)
	}

	err = r.module.Configure(configResp.ConfigMap)
	if err != nil {
		return fmt.Errorf("failed to configure reporter module %s %v", r.GetModuleName(), err)
	}

	resp, err := r.controlService.GetPackageNode(context.Background(), &service.PackageRequest{Session: configResp.Session})
	if err != nil {
		return fmt.Errorf("could not get package node %v", err)
	}

	err = r.module.Report(resp.PackageNode)
	if err != nil {
		return fmt.Errorf("reporter %s failed %v", r.name, err)
	}

	if err := r.module.PostReport(); err != nil {
		return fmt.Errorf("reporter %s failed in PostReport: %v", r.name, err)
	}

	return nil
}

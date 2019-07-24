package reporting

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	goflag "flag"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"
)

// Reporter contains the fields provided to every report
type Reporter struct {
	reportingService service.ReportServiceClient
	controlService   service.ControlServiceClient
	module           ReporterModule
	id               int32
	name             string
}

// ReporterModule defines the methods required to implement a reporter.
type ReporterModule interface {
	Configure(configMap map[string]string) error
	Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient) error
	PostReport() error
}

// NewReporter creates a new reporter.
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

// GetModuleName returns the module name
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
		return fmt.Errorf("failed to create cache directory \"%s\" for module %s: %v", cacheDir, r.GetModuleName(), err)
	}

	err = os.MkdirAll(outDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to create output directory \"%s\" for module %s: %v", outDir, r.GetModuleName(), err)
	}

	err = r.module.Configure(configResp.ConfigMap)
	if err != nil {
		return fmt.Errorf("failed to configure reporter module %s: %v", r.GetModuleName(), err)
	}

	err = r.module.Report(r.controlService, r.reportingService)
	if err != nil {
		return fmt.Errorf("reporter %s failed: %v", r.name, err)
	}

	if err := r.module.PostReport(); err != nil {
		return fmt.Errorf("reporter %s failed in PostReport: %v", r.name, err)
	}

	return nil
}

func GetTrgtPathInfo(trgt *service.FileNode, pkgNode *service.PackageNode) *service.PathInfo {
	for _, pathInfo := range trgt.Paths {
		if pathInfo.Link == pkgNode.Name {
			return pathInfo
		}
	}
	path, err := service.GetPathInfo(trgt)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

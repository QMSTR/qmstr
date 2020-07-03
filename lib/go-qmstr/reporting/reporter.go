package reporting

import (
	"fmt"
	"log"
	"os"

	"golang.org/x/net/context"

	goflag "flag"

	"github.com/QMSTR/qmstr/lib/go-qmstr/cli"
	"github.com/QMSTR/qmstr/lib/go-qmstr/module"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	flag "github.com/spf13/pflag"
)

// Reporter contains the fields provided to every report
type Reporter struct {
	module.MasterClient
	module ReporterModule
	id     int32
	name   string
}

// ReporterModule defines the methods required to implement a reporter.
type ReporterModule interface {
	Configure(configMap map[string]string) error
	Report(masterClient *module.MasterClient) error
	PostReport() error
}

var CountReporters int32

// NewReporter creates a new reporter.
func NewReporter(repModule ReporterModule) *Reporter {
	var serviceAddress string
	var rprID int32
	CountReporters++
	flag.StringVar(&serviceAddress, "rserv", "localhost:50051", "Reporting service address")
	flag.Int32Var(&rprID, "rid", -1, "unique reporter id")
	flag.CommandLine.AddGoFlagSet(goflag.CommandLine)
	flag.Parse()

	mc := module.NewMasterClient(serviceAddress)

	return &Reporter{MasterClient: mc, id: rprID, module: repModule}
}

// GetModuleName returns the module name
func (r *Reporter) GetModuleName() string {
	return r.name
}

// RunReporterModule is the main driver function for each reporter.
func (r *Reporter) RunReporterModule() error {
	configResp, err := r.RptSvcClient.GetReporterConfig(context.Background(), &service.ReporterConfigRequest{ReporterID: r.id})
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

	// Initialize reporter
	_, err = r.CtrlSvcClient.InitModule(context.Background(), &service.InitModuleRequest{
		ModuleName: r.name})
	if err != nil {
		return fmt.Errorf("%v: %v", err, cli.ReturnCodeServerCommunicationError)
	}
	err = r.module.Configure(configResp.ConfigMap)
	if err != nil {
		return fmt.Errorf("failed to configure reporter module %s: %v", r.GetModuleName(), err)
	}

	err = r.module.Report(&r.MasterClient)
	if err != nil {
		return fmt.Errorf("reporter %s failed: %v", r.name, err)
	}

	if err := r.module.PostReport(); err != nil {
		return fmt.Errorf("reporter %s failed in PostReport: %v", r.name, err)
	}

	msg := fmt.Sprintf("Reporter %s finished successfully", r.name)
	log.Println(msg)
	// Ping master server that the reporter finished
	r.CtrlSvcClient.ShutdownModule(context.Background(), &service.ShutdownModuleRequest{Message: msg, DB: false})

	return nil
}

// ReduceReportersCounter is called everytime a reporter finishes its process.
// When it reaches 0, it sends a signal to close the reporting phase
func ReduceReportersCounter() {
	CountReporters--
	if CountReporters == 0 { // all reporters have finished
		// close reporting phase
		close(cli.ModulesAreDone)
	}
	if CountReporters < 0 {
		log.Printf("WARNING: Reporters count cannot be minus")
	}
}

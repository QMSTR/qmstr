package master

import (
	"errors"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseReport struct {
	genericServerPhase
	config []config.Reporting
}

func (phase *serverPhaseReport) Activate() error {
	log.Println("Reporting activated")
	for idx, reporterConfig := range phase.config {
		reporterName := reporterConfig.Reporter

		cmd := exec.Command(reporterName, "--rserv", phase.serverConfig.RPCAddress, "--rid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logModuleError(reporterName, out)
			return err
		}
		log.Printf("Reporter %s finished successfully: %s\n", reporterName, out)
	}
	return nil
}

func (phase *serverPhaseReport) Shutdown() error {
	return nil
}

func (phase *serverPhaseReport) GetPhaseId() int32 {
	return phase.phaseId
}

func (phase *serverPhaseReport) Build(in *service.BuildMessage) (*service.BuildResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseReport) GetReporterConfig(in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	idx := in.ReporterID
	if idx < 0 || idx >= int32(len(phase.config)) {
		return nil, fmt.Errorf("Invalid reporter id %d", idx)
	}
	config := phase.config[idx]

	config.Config["name"] = config.Name

	// Set cacheDir, if not overriden
	if _, ok := config.Config["cacheDir"]; !ok {
		config.Config["cacheDir"] = filepath.Join(phase.serverConfig.CacheDir, config.Reporter, config.PosixName)
	}
	// Set output dir, if not overriden
	if _, ok := config.Config["outputDir"]; !ok {
		config.Config["outputDir"] = filepath.Join(phase.serverConfig.OutputDir, config.Reporter, config.PosixName)
	}

	return &service.ReporterConfigResponse{ConfigMap: config.Config, Session: phase.session,
		Name: config.Name}, nil
}

func (phase *serverPhaseReport) GetAnalyzerConfig(in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseReport) GetNodes(in *service.NodeRequest) (*service.NodeResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseReport) SendNodes(in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseReport) GetPackageNode(in *service.ReportRequest) (*service.ReportResponse, error) {
	node, err := phase.db.GetPackageNode(in.Session)
	if err != nil {
		return nil, err
	}
	return &service.ReportResponse{PackageNode: node}, nil
}

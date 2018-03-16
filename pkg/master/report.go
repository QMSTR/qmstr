package master

import (
	"errors"
	"fmt"
	"log"
	"os/exec"

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

		cmd := exec.Command(reporterName, "--rserv", phase.rpcAddress, "--rid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Printf("Reporter %s failed with: %s\n", reporterName, out)
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
	return &service.ReporterConfigResponse{ConfigMap: config.Config, TypeSelector: config.Selector}, nil
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

func (phase *serverPhaseReport) GetReportNodes(in *service.ReportRequest, streamServer service.ReportService_GetReportNodesServer) error {
	log.Println("Nodes requested")
	nodes, err := phase.db.GetFileNodesByType(in.Type, true)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		err = streamServer.Send(&service.ReportResponse{FileNode: node})
		if err != nil {
			return err
		}
	}
	return nil
}

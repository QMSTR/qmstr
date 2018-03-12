package master

import (
	"errors"
	"log"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseReport struct {
	genericServerPhase
	config []config.Reporting
}

func (phase *serverPhaseReport) Activate() error {
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

func (phase *serverPhaseReport) GetConfig(in *service.ConfigRequest) (*service.ConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseReport) GetNodes(in *service.NodeRequest) (*service.NodeResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseReport) SendNodes(in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (s *serverPhaseReport) Report(in *service.ReportRequest, streamServer service.ReportService_ReportServer) error {
	log.Println("Report requested")
	return nil
}

package master

import (
	"errors"
	"log"

	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseReport struct {
	genericServerPhase
}

func (phase *serverPhaseReport) Activate() bool {
	return false
}

func (phase *serverPhaseReport) GetPhaseId() int32 {
	return phase.phaseId
}

func (phase *serverPhaseReport) Build(in *service.BuildMessage) (*service.BuildResponse, error) {
	return nil, errors.New("Get  off")
}

func (phase *serverPhaseReport) GetNodes(in *service.NodeRequest) (*service.NodeResponse, error) {
	return nil, errors.New("Get  off")
}

func (phase *serverPhaseReport) SendNodes(in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	return nil, errors.New("Get  off")
}

func (s *serverPhaseReport) Report(in *service.ReportRequest, streamServer service.ReportService_ReportServer) error {
	log.Println("Report requested")
	return nil
}

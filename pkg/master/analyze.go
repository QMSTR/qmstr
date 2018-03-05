package master

import (
	"errors"
	"log"

	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseAnalysis struct {
	genericServerPhase
}

func (phase *serverPhaseAnalysis) Activate() bool {
	return false
}

func (phase *serverPhaseAnalysis) GetPhaseId() int32 {
	return phase.phaseId
}

func (phase *serverPhaseAnalysis) Build(in *service.BuildMessage) (*service.BuildResponse, error) {
	return nil, errors.New("Get  off")
}

func (s *serverPhaseAnalysis) GetNodes(in *service.NodeRequest) (*service.NodeResponse, error) {
	log.Println("Nodes requested")
	nodes, err := s.db.GetFileNodesByType(in.Type, true)
	if err != nil {
		return nil, err
	}
	resp := &service.NodeResponse{FileNodes: nodes}
	return resp, nil
}

func (s *serverPhaseAnalysis) SendNodes(in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	log.Println("Nodes received")

	nodes := in.FileNodes

	for _, node := range nodes {
		// TODO remove everything but InfoNodes and UID
		s.db.AlterNode(node)
	}

	return &service.AnalysisResponse{Success: true}, nil
}

func (phase *serverPhaseAnalysis) Report(in *service.ReportRequest, streamServer service.ReportService_ReportServer) error {
	return errors.New("Get  off")
}

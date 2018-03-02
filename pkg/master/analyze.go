package master

import (
	"errors"
	"log"

	pb "github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseAnalysis struct {
	genericServerPhase
}

func (phase *serverPhaseAnalysis) Build(in *pb.BuildMessage) (*pb.BuildResponse, error) {
	return nil, errors.New("Get  off")
}

func (s *serverPhaseAnalysis) GetNodes(in *pb.NodeRequest) (*pb.NodeResponse, error) {
	log.Println("Nodes requested")

	return &pb.NodeResponse{FileNodes: nil}, nil
}

func (s *serverPhaseAnalysis) SendNodes(in *pb.AnalysisMessage) (*pb.AnalysisResponse, error) {
	log.Println("Nodes received")

	return &pb.AnalysisResponse{Success: true}, nil
}

func (phase *serverPhaseAnalysis) Report(in *pb.ReportRequest, streamServer pb.ReportService_ReportServer) error {
	return errors.New("Get  off")
}

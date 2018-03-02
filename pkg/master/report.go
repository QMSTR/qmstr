package master

import (
	"errors"
	"log"

	pb "github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseReport struct {
	genericServerPhase
}

func (phase *serverPhaseReport) Build(in *pb.BuildMessage) (*pb.BuildResponse, error) {
	return nil, errors.New("Get  off")
}

func (phase *serverPhaseReport) GetNodes(in *pb.NodeRequest) (*pb.NodeResponse, error) {
	return nil, errors.New("Get  off")
}

func (phase *serverPhaseReport) SendNodes(in *pb.AnalysisMessage) (*pb.AnalysisResponse, error) {
	return nil, errors.New("Get  off")
}

func (s *serverPhaseReport) Report(in *pb.ReportRequest, streamServer pb.ReportService_ReportServer) error {
	log.Println("Report requested")
	return nil
}

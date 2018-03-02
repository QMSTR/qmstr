package master

import (
	"errors"

	pb "github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseBuild struct {
	genericServerPhase
}

func (phase *serverPhaseBuild) Build(in *pb.BuildMessage) (*pb.BuildResponse, error) {
	return &pb.BuildResponse{Success: true}, nil
}

func (phase *serverPhaseBuild) GetNodes(in *pb.NodeRequest) (*pb.NodeResponse, error) {
	return nil, errors.New("Get  off")
}

func (phase *serverPhaseBuild) SendNodes(in *pb.AnalysisMessage) (*pb.AnalysisResponse, error) {
	return nil, errors.New("Get  off")
}

func (phase *serverPhaseBuild) Report(in *pb.ReportRequest, streamServer pb.ReportService_ReportServer) error {
	return errors.New("Get  off")
}

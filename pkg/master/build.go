package master

import (
	"errors"
	"log"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseBuild struct {
	genericServerPhase
}

func newBuildPhase(session string, masterConfig *config.MasterConfig, db *database.DataBase) serverPhase {
	return &serverPhaseBuild{
		genericServerPhase{Name: "Build", session: session, masterConfig: masterConfig, db: db},
	}
}

func (phase *serverPhaseBuild) Activate() error {
	return nil
}

func (phase *serverPhaseBuild) Shutdown() error {
	phase.db.AwaitBuildComplete()
	return nil
}

func (phase *serverPhaseBuild) GetPhaseID() int32 {
	return PhaseIDBuild
}

func (phase *serverPhaseBuild) Build(in *service.BuildMessage) (*service.BuildResponse, error) {
	for _, node := range in.FileNodes {
		log.Printf("Adding file node %s", node.Path)
		phase.db.AddFileNode(node)
	}
	return &service.BuildResponse{Success: true}, nil
}

func (phase *serverPhaseBuild) GetReporterConfig(in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseBuild) GetAnalyzerConfig(in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseBuild) GetNodes(in *service.NodeRequest) (*service.NodeResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseBuild) SendNodes(in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	return nil, errors.New("Wrong phase")
}

package master

import (
	"log"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseBuild struct {
	genericServerPhase
}

func newBuildPhase(session string, masterConfig *config.MasterConfig, db *database.DataBase, server *server) serverPhase {
	return &serverPhaseBuild{
		genericServerPhase{Name: "Build", session: session, masterConfig: masterConfig, db: db, server: server},
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

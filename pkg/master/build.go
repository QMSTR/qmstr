package master

import (
	"log"
	"path/filepath"
	"strings"

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
		for _, substitution := range phase.masterConfig.Server.PathSub {
			node.Path = strings.Replace(node.Path, substitution.Old, substitution.New, 1)
		}
		relPath, err := filepath.Rel(phase.masterConfig.Server.BuildPath, node.Path)
		if err != nil {
			return nil, err
		}
		node.Path = relPath
		log.Printf("Adding file node %s", node.Path)
		phase.db.AddFileNode(node)
	}
	return &service.BuildResponse{Success: true}, nil
}

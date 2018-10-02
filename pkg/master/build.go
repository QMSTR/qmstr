package master

import (
	"log"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
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
	phase.db.CloseInsertQueue()
	return nil
}

func (phase *serverPhaseBuild) GetPhaseID() service.Phase {
	return service.Phase_BUILD
}

func (phase *serverPhaseBuild) Build(in *service.BuildMessage) (*service.BuildResponse, error) {
	buildPath := phase.masterConfig.Server.BuildPath
	pathSub := phase.masterConfig.Server.PathSub
	for _, node := range in.FileNodes {
		err := common.SetRelativePath(node, buildPath, pathSub)
		if err != nil {
			return nil, err
		}
		for _, derNode := range node.DerivedFrom {
			err := common.SetRelativePath(derNode, buildPath, pathSub)
			if err != nil {
				return nil, err
			}
		}
		log.Printf("Adding file node %s", node.Path)
		phase.db.AddBuildFileNode(node)
	}
	return &service.BuildResponse{Success: true}, nil
}

func (phase *serverPhaseBuild) ExportGraph(in *service.ExportRequest) (*service.ExportResponse, error) {
	phase.db.Sync()
	err := phase.requestExport()
	if err != nil {
		return nil, err
	}
	return &service.ExportResponse{Success: true}, nil
}

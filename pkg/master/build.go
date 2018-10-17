package master

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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
	phase.db.OpenInsertQueue()
	return nil
}

func (phase *serverPhaseBuild) Shutdown() error {
	phase.db.CloseInsertQueue()
	return nil
}

func (phase *serverPhaseBuild) GetPhaseID() service.Phase {
	return service.Phase_BUILD
}

func (phase *serverPhaseBuild) Build(stream service.BuildService_BuildServer) error {
	buildPath := phase.masterConfig.Server.BuildPath
	pathSub := phase.masterConfig.Server.PathSub
	for {
		fileNode, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&service.BuildResponse{
				Success: true,
			})
		}
		if err != nil {
			return err
		}

		common.SanitizeFileNode(fileNode, buildPath, pathSub, phase.db, "")
		log.Printf("Adding file node %s", fileNode.Path)
		phase.db.AddFileNode(fileNode)
	}
}

func (phase *serverPhaseBuild) PushFile(in *service.PushFileMessage) (*service.PushFileResponse, error) {
	pushDir := filepath.Join(common.ContainerBuildDir, common.ContainerPushFilesDirName)
	if err := os.MkdirAll(pushDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("failed to create dir for uploaded files: %v", err)
	}
	var filename string
	if in.Name != "" {
		filename = in.Name
	} else {
		filename = in.Hash
	}
	f, err := os.Create(filepath.Join(pushDir, filename))
	if err != nil {
		return nil, fmt.Errorf("failed to create uploaded file: %v", err)
	}
	defer f.Close()

	_, err = io.Copy(f, bytes.NewReader(in.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to write uploaded file: %v", err)
	}

	f.Sync()
	f.Seek(0, 0)
	hash, err := common.Hash(f)
	if err != nil {
		return nil, fmt.Errorf("failed to verify uploaded file: %v", err)
	}
	if hash != in.Hash {
		return nil, fmt.Errorf("failed to verify uploaded file %s != %s", hash, in.Hash)
	}

	return &service.PushFileResponse{Path: f.Name()}, nil
}

func (phase *serverPhaseBuild) ExportGraph(in *service.ExportRequest) (*service.ExportResponse, error) {
	phase.db.Sync()
	err := phase.requestExport()
	if err != nil {
		return nil, err
	}
	return &service.ExportResponse{Success: true}, nil
}

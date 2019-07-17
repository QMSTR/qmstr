package master

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/config"
	"github.com/QMSTR/qmstr/lib/go-qmstr/database"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

type serverPhaseBuild struct {
	genericServerPhase
}

func newBuildPhase(masterConfig *config.MasterConfig, db *database.DataBase, server *server, done bool) serverPhase {
	return &serverPhaseBuild{
		genericServerPhase{Name: "Build", masterConfig: masterConfig, db: db, server: server, done: done},
	}
}

func (phase *serverPhaseBuild) Activate() error {
	phase.server.publishEvent(&service.Event{Class: service.EventClass_PHASE, Message: "Activating build phase"})
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

		if !fileNode.IsValid() {
			return errors.New("invalid file node")
		}

		if err := common.SanitizeFileNode(fileNode, buildPath, pathSub, phase.db, ""); err != nil {
			return fmt.Errorf("failed sanitising file node %v: %v", fileNode, err)
		}
		// log the last path added to the node
		log.Printf("Adding file node %s with hash:%s", service.GetFilePath(fileNode), fileNode.Hash)
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

func (phase *serverPhaseBuild) CreatePathInfo(in *service.PathInfo) (*service.BuildResponse, error) {
	if !in.IsValid() {
		return nil, errors.New("invalid path info node")
	}
	if _, err := phase.db.GetPathInfobyPath(in.Path); err != database.ErrNoSuchPath {
		return nil, errors.New("path already created")
	}
	log.Printf("Adding path info node %s", in.Path)
	phase.db.AddPathInfo(in)
	return &service.BuildResponse{Success: true}, nil
}

func (phase *serverPhaseBuild) CreatePackage(in *service.PackageNode) (*service.BuildResponse, error) {
	if !in.IsValid() {
		return nil, errors.New("invalid package node")
	}
	if _, err := phase.db.GetPackageNodeByName(in.Name); err != database.ErrNoSuchPackage {
		return nil, errors.New("package already created")
	}
	log.Printf("Adding package node %s", in.Name)
	phase.db.AddPackageNode(in)
	return &service.BuildResponse{Success: true}, nil
}

func (phase *serverPhaseBuild) CreateProject(in *service.ProjectNode) (*service.BuildResponse, error) {
	if !in.IsValid() {
		return nil, errors.New("invalid project node")
	}

	projectNode, err := phase.db.GetProjectNode()
	switch err {
	case database.ErrNoProjectNode:
		// update requested
		if in.Uid != "" {
			return nil, errors.New("can not update non existing project")
		}
	case nil:
		// create requested
		if in.Uid == "" {
			return nil, errors.New("project node already created")
		}
		// update different project node
		if in.Uid != projectNode.Uid {
			return nil, errors.New("can not update existing project with different uid")
		}
	default:
		return nil, err
	}

	log.Printf("Adding project node %s", in)
	phase.db.AddProjectNode(in)
	return &service.BuildResponse{Success: true}, nil
}

func (phase *serverPhaseBuild) DeleteNode(stream service.BuildService_DeleteNodeServer) error {
	for {
		deleteNodeReq, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&service.BuildResponse{
				Success: true,
			})
		}
		if err != nil {
			return err
		}
		_, err = database.Delete(phase.db, deleteNodeReq.Uid)
		if err != nil {
			return err
		}
	}
}

func (phase *serverPhaseBuild) DeleteEdge(in *service.DeleteMessage) (*service.BuildResponse, error) {
	_, err := database.DeleteEdge(phase.db, in.Uid, in.Edge)
	if err != nil {
		return nil, errors.New("Failed removing edge")
	}
	return &service.BuildResponse{Success: true}, nil
}

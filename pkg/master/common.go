package master

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhase interface {
	GetPhaseID() service.Phase
	getName() string
	Activate() error
	Shutdown() error
	getDataBase() (*database.DataBase, error)
	getError() string
	getMasterConfig() *config.MasterConfig
	Build(service.BuildService_BuildServer) error
	GetAnalyzerConfig(*service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error)
	GetReporterConfig(*service.ReporterConfigRequest) (*service.ReporterConfigResponse, error)
	SendInfoNodes(service.AnalysisService_SendInfoNodesServer) error
	SendFileNode(service.AnalysisService_SendFileNodeServer) error
	SendPackageNode(service.AnalysisService_SendPackageNodeServer) error
	SendDiagnosticNode(service.AnalysisService_SendDiagnosticNodeServer) error
	GetFileNode(*service.FileNode, service.ControlService_GetFileNodeServer) error
	GetDiagnosticNode(*service.DiagnosticNode, service.ControlService_GetDiagnosticNodeServer) error
	GetBOM(*service.BOMRequest) (*service.BOM, error)
	GetInfoData(*service.InfoDataRequest) (*service.InfoDataResponse, error)
	ExportSnapshot(*service.ExportRequest) (*service.ExportResponse, error)
	requestExport() error
	getPostInitPhase() service.Phase
	PushFile(*service.PushFileMessage) (*service.PushFileResponse, error)
	CreateProject(*service.ProjectNode) (*service.BuildResponse, error)
	CreatePackage(*service.PackageNode) (*service.BuildResponse, error)
}

type genericServerPhase struct {
	Name          string
	db            *database.DataBase
	masterConfig  *config.MasterConfig
	server        *server
	postInitPhase *service.Phase
}

func (gsp *genericServerPhase) getPostInitPhase() service.Phase {
	if gsp.postInitPhase == nil {
		return service.Phase_BUILD
	}
	return *gsp.postInitPhase
}

func (gsp *genericServerPhase) getDataBase() (*database.DataBase, error) {
	if gsp.db == nil {
		return nil, errors.New("Database not yet available")
	}
	return gsp.db, nil
}

func (gsp *genericServerPhase) getError() string {
	return ""
}

func (gsp *genericServerPhase) getMasterConfig() *config.MasterConfig {
	return gsp.masterConfig
}

func (gsp *genericServerPhase) getName() string {
	return gsp.Name
}

func (gsp *genericServerPhase) Build(stream service.BuildService_BuildServer) error {
	return errors.New("Wrong phase")
}

func (gsp *genericServerPhase) PushFile(in *service.PushFileMessage) (*service.PushFileResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) CreatePackage(in *service.PackageNode) (*service.BuildResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) CreateProject(in *service.ProjectNode) (*service.BuildResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) GetReporterConfig(in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) GetAnalyzerConfig(in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) SendInfoNodes(stream service.AnalysisService_SendInfoNodesServer) error {
	return errors.New("Wrong phase")
}

func (gsp *genericServerPhase) SendFileNode(stream service.AnalysisService_SendFileNodeServer) error {
	return errors.New("Wrong phase")
}

func (gsp *genericServerPhase) SendPackageNode(stream service.AnalysisService_SendPackageNodeServer) error {
	return errors.New("Wrong phase")
}

func (gsp *genericServerPhase) SendDiagnosticNode(stream service.AnalysisService_SendDiagnosticNodeServer) error {
	return errors.New("Wrong phase")
}

func (gsp *genericServerPhase) GetFileNode(in *service.FileNode, stream service.ControlService_GetFileNodeServer) error {
	db, err := gsp.getDataBase()
	if err != nil {
		return err
	}
	nodeFiles, err := db.GetFileNodesByFileNode(in, true)
	if err != nil {
		return err
	}
	for _, nodeFile := range nodeFiles {
		if gsp.server.currentPhase.GetPhaseID() == service.Phase_ANALYSIS {
			nodeFile.Path = filepath.Join(gsp.masterConfig.Server.BuildPath, nodeFile.Path)
		}
		stream.Send(nodeFile)
	}
	return nil
}

func (gsp *genericServerPhase) GetDiagnosticNode(in *service.DiagnosticNode, stream service.ControlService_GetDiagnosticNodeServer) error {
	db, err := gsp.getDataBase()
	if err != nil {
		return err
	}
	diagnosticNodes, err := db.GetDiagnosticNodeBySeverity(in)
	if err != nil {
		return err
	}
	for _, diagnosticNode := range diagnosticNodes {
		stream.Send(diagnosticNode)
	}
	return nil
}

func (gsp *genericServerPhase) GetBOM(in *service.BOMRequest) (*service.BOM, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) GetInfoData(in *service.InfoDataRequest) (*service.InfoDataResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) ExportSnapshot(in *service.ExportRequest) (*service.ExportResponse, error) {
	err := gsp.requestExport()
	if err != nil {
		return nil, err
	}
	pushedFilesDir := filepath.Join(common.ContainerBuildDir, common.ContainerPushFilesDirName)
	fi, err := os.Stat(pushedFilesDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("No pushed files for export found.")
			return &service.ExportResponse{Success: true}, nil
		} else {
			return nil, err
		}
	}
	if fi.IsDir() {
		pushExportDir := filepath.Join(common.ContainerGraphExportDir, common.ContainerPushFilesDirName)
		if err := os.MkdirAll(pushExportDir, os.ModePerm); err != nil {
			return nil, err
		}

		err := filepath.Walk(pushedFilesDir, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			log.Printf("Exporting %s", path)
			src, err := os.Open(path)
			if err != nil {
				return err
			}
			defer src.Close()
			dest, err := os.Create(filepath.Join(pushExportDir, filepath.Base(path)))
			if err != nil {
				return err
			}
			defer dest.Close()

			// TODO verify copied version
			_, err = io.Copy(dest, src)
			if err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return &service.ExportResponse{Success: true}, nil
}

func (gsp *genericServerPhase) requestExport() error {
	// clean the dir
	if err := os.RemoveAll(common.ContainerGraphExportDir); err != nil {
		return fmt.Errorf("failed to clean the export dir: %v", err)
	}
	// create dir
	if err := os.Mkdir(common.ContainerGraphExportDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create the export dir: %v", err)
	}

	resp, err := http.Get("http://localhost:8080/admin/export")
	if err != nil {
		return fmt.Errorf("failed sending export request to dgraph: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed reading dgraph response: %v", err)
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to export graph dgraph answered: [%d] %s", resp.StatusCode, body)
	}

	log.Printf("dgraph export: %s", body)

	return nil
}

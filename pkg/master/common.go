package master

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

type serverPhase interface {
	GetPhaseID() service.Phase
	getName() string
	Activate() error
	Shutdown() error
	getDataBase() (*database.DataBase, error)
	getSession() string
	getError() string
	getMasterConfig() *config.MasterConfig
	Build(*service.BuildMessage) (*service.BuildResponse, error)
	GetAnalyzerConfig(*service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error)
	GetReporterConfig(*service.ReporterConfigRequest) (*service.ReporterConfigResponse, error)
	SendInfoNodes(service.AnalysisService_SendInfoNodesServer) error
	SendFileNode(service.AnalysisService_SendFileNodeServer) error
	SendPackageNode(service.AnalysisService_SendPackageNodeServer) error
	GetFileNode(*service.FileNode, service.ControlService_GetFileNodeServer) error
	GetBOM(*service.BOMRequest) (*service.BOM, error)
	GetInfoData(*service.InfoDataRequest) (*service.InfoDataResponse, error)
	ExportGraph(*service.ExportRequest) (*service.ExportResponse, error)
	requestExport() error
	getPostInitPhase() service.Phase
	PushFile(*service.PushFileMessage) (*service.BuildResponse, error)
}

type genericServerPhase struct {
	Name          string
	db            *database.DataBase
	session       string
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

func (gsp *genericServerPhase) getSession() string {
	return gsp.session
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

func (gsp *genericServerPhase) Build(in *service.BuildMessage) (*service.BuildResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) PushFile(in *service.PushFileMessage) (*service.BuildResponse, error) {
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

func (gsp *genericServerPhase) GetFileNode(in *service.FileNode, stream service.ControlService_GetFileNodeServer) error {
	db, err := gsp.getDataBase()
	if err != nil {
		return err
	}
	nodeFiles, err := db.GetFileNodesByFileNode(in, true)

	for _, nodeFile := range nodeFiles {
		stream.Send(nodeFile)
	}
	return nil
}

func (gsp *genericServerPhase) GetBOM(in *service.BOMRequest) (*service.BOM, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) GetInfoData(in *service.InfoDataRequest) (*service.InfoDataResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) ExportGraph(in *service.ExportRequest) (*service.ExportResponse, error) {
	err := gsp.requestExport()
	if err != nil {
		return nil, err
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

package master

import (
	"errors"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhase interface {
	GetPhaseID() int32
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
	SendFileNodes(service.AnalysisService_SendFileNodesServer) error
	GetFileNode(*service.FileNode, service.ControlService_GetFileNodeServer) error
	GetBOM(*service.BOMRequest) (*service.BOM, error)
	GetInfoData(*service.InfoDataRequest) (*service.InfoDataResponse, error)
}

type genericServerPhase struct {
	Name         string
	db           *database.DataBase
	session      string
	masterConfig *config.MasterConfig
	server       *server
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

func (gsp *genericServerPhase) GetReporterConfig(in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) GetAnalyzerConfig(in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (gsp *genericServerPhase) SendInfoNodes(stream service.AnalysisService_SendInfoNodesServer) error {
	return errors.New("Wrong phase")
}

func (gsp *genericServerPhase) SendFileNodes(stream service.AnalysisService_SendFileNodesServer) error {
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

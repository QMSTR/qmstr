package master

import (
	"errors"
	"fmt"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseInit struct {
	genericServerPhase
}

func newInitServerPhase(session string, masterConfig *config.MasterConfig) *serverPhaseInit {
	return &serverPhaseInit{genericServerPhase{Name: "Init", session: session, masterConfig: masterConfig}}
}

func (phase *serverPhaseInit) Activate() error {
	// Connect to database (dgraph)
	db, err := database.Setup(phase.masterConfig.Server.DBAddress, phase.masterConfig.Server.DBWorkers)
	if err != nil {
		return fmt.Errorf("Could not setup database: %v", err)
	}
	phase.db = db

	phase.initPackage(phase.session)
	return nil
}

func (phase *serverPhaseInit) initPackage(session string) {
	rootPackageNode := &service.PackageNode{Name: phase.masterConfig.Name}
	tmpInfoNode := &service.InfoNode{Type: "metadata", NodeType: service.NodeTypeInfoNode}
	for key, val := range phase.masterConfig.MetaData {
		tmpInfoNode.DataNodes = append(tmpInfoNode.DataNodes, &service.InfoNode_DataNode{Type: key, Data: val, NodeType: service.NodeTypeDataNode})
	}

	if len(tmpInfoNode.DataNodes) > 0 {
		rootPackageNode.AdditionalInfo = []*service.InfoNode{tmpInfoNode}
	}

	rootPackageNode.Session = session
	phase.db.AddPackageNode(rootPackageNode)
}

func (phase *serverPhaseInit) Shutdown() error {
	return nil
}

func (phase *serverPhaseInit) GetPhaseID() int32 {
	return phaseIDInit
}

func (phase *serverPhaseInit) Build(in *service.BuildMessage) (*service.BuildResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseInit) GetAnalyzerConfig(in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseInit) GetNodes(in *service.NodeRequest) (*service.NodeResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseInit) SendNodes(in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseInit) GetReporterConfig(in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

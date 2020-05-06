package master

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/QMSTR/qmstr/lib/go-qmstr/config"
	"github.com/QMSTR/qmstr/lib/go-qmstr/database"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

type serverPhaseAnalysis struct {
	genericServerPhase
	currentAnalyzer *service.Analyzer
	currentToken    int64
	finished        chan interface{}
}

var src = rand.NewSource(time.Now().UnixNano())

func newAnalysisPhase(masterConfig *config.MasterConfig, db *database.DataBase, server *server, done bool) serverPhase {
	return &serverPhaseAnalysis{
		genericServerPhase{Name: "Analysis", masterConfig: masterConfig, db: db, server: server, done: done},
		nil, src.Int63(), make(chan interface{}, 1)}
}

func (phase *serverPhaseAnalysis) Activate() error {
	log.Println("Analysis activated")
	phase.server.publishEvent(&service.Event{Class: service.EventClass_PHASE, Message: "Activating analysis phase"})
	return nil
}

func (phase *serverPhaseAnalysis) InitModule(in *service.InitModuleRequest) (*service.InitModuleResponse, error) {
	analyzer, err := phase.db.GetAnalyzerByName(in.ModuleName)
	if err != nil {
		return nil, err
	}
	phase.currentAnalyzer = analyzer
	src.Seed(phase.currentToken)
	phase.currentToken = src.Int63()
	phase.currentAnalyzer.TrustLevel = in.ExtraConfig

	log.Printf("Running analyzer %s ...\n", in.ModuleName)
	phase.db.OpenInsertQueue()
	phase.server.publishEvent(&service.Event{Class: service.EventClass_MODULE, Message: fmt.Sprintf("Running analyzer %s", in.ModuleName)})
	return &service.InitModuleResponse{}, nil
}

func (phase *serverPhaseAnalysis) Shutdown() error {
	phase.finished <- nil
	phase.done = true
	phase.server.persistPhase()
	log.Println("Analysis phase finished")
	return nil
}

func (phase *serverPhaseAnalysis) GetPhaseID() service.Phase {
	return service.Phase_ANALYSIS
}

func (phase *serverPhaseAnalysis) GetAnalyzerConfig(in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	idx := in.AnalyzerID
	if idx < 0 || idx >= int32(len(phase.masterConfig.Analysis)) {
		return nil, fmt.Errorf("Invalid analyzer id %d", idx)
	}
	config := phase.masterConfig.Analysis[idx]

	if config.Config == nil {
		config.Config = make(map[string]string)
	}

	config.Config["name"] = config.Name
	config.Config["cachedir"] = filepath.Join(ServerCacheDir, config.Analyzer, config.PosixName)
	config.Config["outputdir"] = filepath.Join(ServerOutputDir, config.Analyzer, config.PosixName)

	// Set path substitution, if not overriden
	if config.PathSub == nil || len(config.PathSub) == 0 {
		config.PathSub = phase.masterConfig.Server.PathSub
	}
	phase.currentAnalyzer.PathSub = config.PathSub
	return &service.AnalyzerConfigResponse{ConfigMap: config.Config, PathSub: config.PathSub,
		Token: phase.currentToken, Name: config.Name}, nil
}

func (phase *serverPhaseAnalysis) SendInfoNodes(stream service.AnalysisService_SendInfoNodesServer) error {
	for {
		infoNodeReq, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&service.SendResponse{
				Success: true,
			})
		}
		if err != nil {
			return err
		}
		if infoNodeReq.Token != phase.currentToken {
			log.Println("Analyzer supplied wrong token")
			return errors.New("wrong token supplied")
		}
		for _, infoNode := range infoNodeReq.Infonodes {
			infoNode.Analyzer = []*service.Analyzer{phase.currentAnalyzer}
		}
		err = phase.db.AddInfoNodes(infoNodeReq.Uid, infoNodeReq.Infonodes)
		if err != nil {
			return err
		}
	}
}

func (phase *serverPhaseAnalysis) SendDiagnosticNode(stream service.AnalysisService_SendDiagnosticNodeServer) error {
	for {
		diagnosticNodeReq, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&service.SendResponse{
				Success: true,
			})
		}
		if err != nil {
			return err
		}
		if diagnosticNodeReq.Token != phase.currentToken {
			log.Println("Analyzer supplied wrong token")
			return errors.New("wrong token supplied")
		}
		diagnosticNode := diagnosticNodeReq.Diagnosticnode
		diagnosticNode.Analyzer = []*service.Analyzer{phase.currentAnalyzer}
		err = phase.db.AddDiagnosticNode(diagnosticNodeReq.Uid, diagnosticNode)
		if err != nil {
			return err
		}
	}
}

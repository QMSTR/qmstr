package master

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
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

func newAnalysisPhase(masterConfig *config.MasterConfig, db *database.DataBase, server *server) serverPhase {
	return &serverPhaseAnalysis{
		genericServerPhase{Name: "Analysis", masterConfig: masterConfig, db: db, server: server},
		nil, src.Int63(), make(chan interface{}, 1)}
}

func (phase *serverPhaseAnalysis) Activate() error {
	log.Println("Analysis activated")
	phase.server.publishEvent(&service.Event{Class: service.EventClass_PHASE, Message: "Activating analysis phase"})
	for idx, anaConfig := range phase.masterConfig.Analysis {
		analyzerName := anaConfig.Analyzer

		analyzer, err := phase.db.GetAnalyzerByName(analyzerName)
		if err != nil {
			return err
		}
		phase.currentAnalyzer = analyzer
		src.Seed(phase.currentToken)
		phase.currentToken = src.Int63()
		phase.currentAnalyzer.TrustLevel = anaConfig.TrustLevel

		log.Printf("Running analyzer %s ...\n", analyzerName)
		phase.db.OpenInsertQueue()
		phase.server.publishEvent(&service.Event{Class: service.EventClass_MODULE, Message: fmt.Sprintf("Running analyzer %s", analyzerName)})
		cmd := exec.Command(analyzerName, "--aserv", phase.masterConfig.Server.RPCAddress, "--aid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logModuleError(analyzerName, out)
			errMsg := fmt.Sprintf("Analyzer %s failed", analyzerName)
			phase.server.publishEvent(&service.Event{Class: service.EventClass_MODULE, Message: errMsg})
			phase.db.CloseInsertQueue()
			return errors.New(errMsg)
		}
		phase.db.CloseInsertQueue()
		phase.server.publishEvent(&service.Event{Class: service.EventClass_MODULE, Message: fmt.Sprintf("Analyzer %s successfully finished", analyzerName)})
		log.Printf("Analyzer %s finished successfully:\n%s\n", analyzerName, out)
	}

	phase.finished <- nil
	log.Println("Analysis phase finished")
	return nil
}

func (phase *serverPhaseAnalysis) Shutdown() error {
	log.Println("Waiting for analysis to be finished")
	<-phase.finished
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
		infoNode := infoNodeReq.Infonode
		infoNode.Analyzer = []*service.Analyzer{phase.currentAnalyzer}
		err = phase.db.AddInfoNodes(infoNodeReq.Uid, infoNode)
		if err != nil {
			return err
		}
	}
}

func (phase *serverPhaseAnalysis) SendPackageNode(stream service.AnalysisService_SendPackageNodeServer) error {
	for {
		pkgNodeReq, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&service.SendResponse{
				Success: true,
			})
		}
		if err != nil {
			return err
		}
		if pkgNodeReq.Token != phase.currentToken {
			log.Println("Analyzer supplied wrong token")
			return errors.New("wrong token supplied")
		}
		pkgNode := pkgNodeReq.Packagenode

		for _, target := range pkgNode.Targets {
			err = common.SetRelativePath(target, phase.masterConfig.Server.BuildPath, nil)
			if err != nil {
				return err
			}
			log.Printf("Adding file node %v to package targets.", target.Path)
		}

		phase.db.AddPackageNode(pkgNode)
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
		err = phase.db.AddDiagnosticNodes(diagnosticNodeReq.Uid, diagnosticNode)
		if err != nil {
			return err
		}
	}

}

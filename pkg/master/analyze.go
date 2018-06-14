package master

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseAnalysis struct {
	genericServerPhase
	currentAnalyzer *service.Analyzer
	currentToken    int64
	finished        chan interface{}
}

var src = rand.NewSource(time.Now().UnixNano())

func newAnalysisPhase(session string, masterConfig *config.MasterConfig, db *database.DataBase, server *server) serverPhase {
	return &serverPhaseAnalysis{
		genericServerPhase{Name: "Analysis", session: session, masterConfig: masterConfig, db: db, server: server},
		nil, src.Int63(), make(chan interface{}, 1)}
}

func (phase *serverPhaseAnalysis) Activate() error {
	log.Println("Analysis activated")
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
		phase.server.publishEvent(&service.Event{Class: string(EventModule), Message: fmt.Sprintf("Running analyzer %s", analyzerName)})
		cmd := exec.Command(analyzerName, "--aserv", phase.masterConfig.Server.RPCAddress, "--aid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logModuleError(analyzerName, out)
			phase.server.publishEvent(&service.Event{Class: string(EventModule), Message: fmt.Sprintf("Analyzer %s failed", analyzerName)})
			return err
		}
		phase.server.publishEvent(&service.Event{Class: string(EventModule), Message: fmt.Sprintf("Analyzer %s successfully finished", analyzerName)})
		log.Printf("Analyzer %s finished successfully: %s\n", analyzerName, out)
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

func (phase *serverPhaseAnalysis) GetPhaseID() int32 {
	return PhaseIDAnalysis
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

	// Set cachedir, if not overriden
	if _, ok := config.Config["cachedir"]; !ok {
		config.Config["cachedir"] = filepath.Join(phase.masterConfig.Server.CacheDir, config.Analyzer, config.PosixName)
	}
	// Set output dir, if not overriden
	if _, ok := config.Config["outputdir"]; !ok {
		config.Config["outputdir"] = filepath.Join(phase.masterConfig.Server.OutputDir, config.Analyzer, config.PosixName)
	}
	// Set path substitution, if not overriden
	if config.PathSub == nil || len(config.PathSub) == 0 {
		config.PathSub = phase.masterConfig.Server.PathSub
	}
	phase.currentAnalyzer.PathSub = config.PathSub
	return &service.AnalyzerConfigResponse{ConfigMap: config.Config, PathSub: config.PathSub,
		Token: phase.currentToken, Name: config.Name, Session: phase.session}, nil
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
		infoNode.NodeType = service.NodeTypeInfoNode
		infoNode.Analyzer = []*service.Analyzer{phase.currentAnalyzer}
		err = phase.db.AddInfoNodes(infoNodeReq.Uid, infoNode)
		if err != nil {
			return err
		}
	}
}

func (phase *serverPhaseAnalysis) SendFileNodes(stream service.AnalysisService_SendFileNodesServer) error {
	for {
		fileNodeReq, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&service.SendResponse{
				Success: true,
			})
		}
		if err != nil {
			return err
		}
		if fileNodeReq.Token != phase.currentToken {
			log.Println("Analyzer supplied wrong token")
			return errors.New("wrong token supplied")
		}
		fileNode := fileNodeReq.Filenode
		fileNode.NodeType = service.NodeTypeFileNode
		err = phase.db.AddFileNodes(fileNodeReq.Uid, fileNode)
		if err != nil {
			return err
		}
	}
}

func (phase *serverPhaseAnalysis) GetFileNode(in *service.FileNode, stream service.ControlService_GetFileNodeServer) error {
	// TODO get rid of code duplication
	db, err := phase.getDataBase()
	if err != nil {
		return err
	}
	nodeFiles, err := db.GetFileNodesByFileNode(in, true)
	if err != nil {
		return err
	}

	for _, nodeFile := range nodeFiles {
		for _, substitution := range phase.currentAnalyzer.PathSub {
			nodeFile.Path = strings.Replace(nodeFile.Path, substitution.Old, substitution.New, 1)
		}
		stream.Send(nodeFile)
	}
	return nil
}

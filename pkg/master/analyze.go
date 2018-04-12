package master

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os/exec"
	"time"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/service"
)

type serverPhaseAnalysis struct {
	genericServerPhase
	config          []config.Analysis
	currentAnalyzer *service.Analyzer
	currentToken    int64
	finished        chan interface{}
}

var src = rand.NewSource(time.Now().UnixNano())

func newAnalysisPhase(genericPhase genericServerPhase, anaConfig []config.Analysis) *serverPhaseAnalysis {
	return &serverPhaseAnalysis{genericPhase, anaConfig, nil, src.Int63(), make(chan interface{}, 1)}
}

func (phase *serverPhaseAnalysis) Activate() error {
	log.Println("Analysis activated")
	for idx, anaConfig := range phase.config {
		analyzerName := anaConfig.Analyzer

		analyzer, err := phase.db.GetAnalyzerByName(analyzerName)
		if err != nil {
			return err
		}
		phase.currentAnalyzer = analyzer
		src.Seed(phase.currentToken)
		phase.currentToken = src.Int63()

		log.Printf("Running analyzer %s ...\n", analyzerName)
		cmd := exec.Command(analyzerName, "--aserv", phase.rpcAddress, "--aid", fmt.Sprintf("%d", idx))
		out, err := cmd.CombinedOutput()
		if err != nil {
			logModuleError(analyzerName, out)
			return err
		}
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

func (phase *serverPhaseAnalysis) GetPhaseId() int32 {
	return phase.phaseId
}

func (phase *serverPhaseAnalysis) Build(in *service.BuildMessage) (*service.BuildResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseAnalysis) GetReporterConfig(in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	return nil, errors.New("Wrong phase")
}

func (phase *serverPhaseAnalysis) GetAnalyzerConfig(in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	idx := in.AnalyzerID
	if idx < 0 || idx >= int32(len(phase.config)) {
		return nil, fmt.Errorf("Invalid analyzer id %d", idx)
	}
	config := phase.config[idx]
	return &service.AnalyzerConfigResponse{ConfigMap: config.Config, TypeSelector: config.Selector, PathSub: config.PathSub, Token: phase.currentToken}, nil
}

func (phase *serverPhaseAnalysis) GetNodes(in *service.NodeRequest) (*service.NodeResponse, error) {
	log.Println("Nodes requested")
	nodes, err := phase.db.GetFileNodesByType(in.Type, true)
	if err != nil {
		return nil, err
	}
	resp := &service.NodeResponse{FileNodes: nodes}
	return resp, nil
}

func (phase *serverPhaseAnalysis) SendNodes(in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	log.Println("Nodes received")

	if in.Token != phase.currentToken {
		fmt.Println("Analyzer supplied wrong token")
		return nil, errors.New("wrong token supplied")
	}

	for hash, inodes := range in.ResultMap {
		log.Printf("Processing node %s with %d info nodes\n", hash, len(inodes.Inodes))
		fileNode, err := phase.db.GetNodeByHash(hash, true)
		if err != nil {
			return &service.AnalysisResponse{Success: false}, err
		}
		for idx, inode := range inodes.Inodes {
			infoNode, err := phase.db.GetInfoNodeByDataNode(inode.Type, inode.DataNodes...)
			if err != nil {
				return nil, err
			}
			// prevent inserting data nodes twice
			infoNode.DataNodes = nil
			infoNode.Analyzer = append(infoNode.Analyzer, phase.currentAnalyzer)
			inodes.Inodes[idx] = infoNode
		}
		fileNode.AdditionalInfo = append(fileNode.AdditionalInfo, inodes.Inodes...)
		phase.db.AlterFileNode(fileNode)
	}

	return &service.AnalysisResponse{Success: true}, nil
}

func (phase *serverPhaseAnalysis) GetReportNodes(in *service.ReportRequest, streamServer service.ReportService_GetReportNodesServer) error {
	return errors.New("Wrong phase")
}

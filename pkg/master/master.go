package master

import (
	"fmt"
	"log"
	"net"
	"path/filepath"
	"sync"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/analysis"
	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/report"
	"google.golang.org/grpc"
)

var quitServer chan interface{}

type server struct {
	db             *database.DataBase
	analyzerQueue  chan analysis.Analysis
	analysisClosed chan bool
	serverMutex    *sync.Mutex
	analysisDone   bool
	analysis       []config.Analysis
	reporting      []config.Reporting
}

func (s *server) Report(in *pb.ReportMessage, streamServer pb.BuildService_ReportServer) error {
	log.Println("Report requested")

	s.serverMutex.Lock()

	for _, currentReport := range s.reporting {
		nodeSelector := currentReport.Selector

		var reporter report.Reporter
		switch currentReport.ReportType {
		case "license":
			reporter = report.NewLicenseReporter()
		case "copyrightHolder":
			reporter = report.NewCopyrightHolderReporter()
		default:
			streamServer.Send(&pb.ReportResponse{Success: false, ResponseMessage: fmt.Sprintf("No such reporter %s", currentReport.ReportType)})
		}

		nodes, err := s.db.GetNodesByType(nodeSelector, true, currentReport.Name)
		if err != nil {
			streamServer.Send(&pb.ReportResponse{Success: false, ResponseMessage: err.Error()})
		}

		nodeRefs := []*database.Node{}
		for _, node := range nodes {
			nodeRefs = append(nodeRefs, &node)
		}

		reportResponse, err := reporter.Generate(nodeRefs)
		reportResponse.Reporter = currentReport.ReportType
		if err != nil {
			streamServer.Send(&pb.ReportResponse{Success: false, ResponseMessage: err.Error()})
		}
		streamServer.Send(reportResponse)
	}
	return nil
}

func (s *server) Analyze(ctx context.Context, in *pb.AnalysisMessage) (*pb.AnalysisResponse, error) {
	log.Println("Analysis requested")

	s.serverMutex.Lock()
	if s.analysisDone {
		s.serverMutex.Unlock()
		return &pb.AnalysisResponse{Success: false}, fmt.Errorf("Analysis already done")
	}

	s.db.AwaitBuildComplete()

	go func() {
		log.Println("Analysis queue worker started")
		for ana := range s.analyzerQueue {
			analysis.RunAnalysis(ana)
		}
		s.analysisClosed <- true
	}()

	for _, currentAnalysis := range s.analysis {
		nodeSelector := currentAnalysis.Selector
		analyzerSelector := currentAnalysis.Analyzer

		var analyzer analysis.Analyzer
		var err error
		switch analyzerSelector {
		case "spdx":
			analyzer = analysis.NewSpdxAnalyzer(currentAnalysis.Config, s.db)
		case "ninka":
			analyzer = analysis.NewNinkaAnalyzer(currentAnalysis.Config, s.db)
		case "scancode":
			analyzer, err = analysis.NewScancodeAnalyzer(currentAnalysis.Config, s.db)
			if err != nil {
				return &pb.AnalysisResponse{Success: false}, err
			}
		default:
			return &pb.AnalysisResponse{Success: false}, fmt.Errorf("No such analyzer %s", analyzerSelector)
		}

		nodes, err := s.db.GetNodesByType(nodeSelector, false, "")
		if err != nil {
			return &pb.AnalysisResponse{Success: false}, err
		}

		anaNodes := []analysis.AnalysisNode{}
		for _, node := range nodes {
			anaNodes = append(anaNodes, analysis.NewAnalysisNode(node, currentAnalysis.PathSub, s.db))
		}

		s.analyzerQueue <- analysis.Analysis{Name: analyzerSelector, Nodes: anaNodes, Analyzer: analyzer}
	}
	// wait for analysis to finish
	close(s.analyzerQueue)
	log.Println("Waiting for analyzers to finish")
	s.analysisDone = <-s.analysisClosed
	log.Println("All analyzers finished")

	s.serverMutex.Unlock()
	return &pb.AnalysisResponse{Success: true}, nil
}

func (s *server) Build(ctx context.Context, in *pb.BuildMessage) (*pb.BuildResponse, error) {
	// Compiliation
	for _, compile := range in.GetCompilations() {
		log.Printf("Compiled %v", compile)

		uidTrgt, err := s.db.HasNode(compile.Target.GetHash())
		if err != nil {
			return &pb.BuildResponse{Success: false}, err
		}

		// no such node exist
		if uidTrgt == "" {
			src := database.NewNode(compile.Source.GetPath(), compile.Source.GetHash())
			src.Type = database.ArtifactTypeSrc
			trgt := database.NewNode(compile.Target.GetPath(), compile.Target.GetHash())
			trgt.DerivedFrom = []*database.Node{&src}
			trgt.Type = database.ArtifactTypeObj

			s.db.AddNode(&trgt)
		}
	}

	// Linking
	for _, bin := range in.GetBinary() {
		log.Printf("Linked target: %v", bin)

		uidTrgt, err := s.db.HasNode(bin.Target.GetHash())
		if err != nil {
			return &pb.BuildResponse{Success: false}, err
		}

		deps := []*database.Node{}
		// no such node exist
		if uidTrgt == "" {
			for _, dep := range bin.GetInput() {
				depNode := database.NewNode(dep.GetPath(), dep.GetHash())
				depNode.Name = filepath.Base(dep.GetPath())
				deps = append(deps, &depNode)
			}
			trgt := database.NewNode(bin.Target.GetPath(), bin.Target.GetHash())
			trgt.DerivedFrom = deps
			trgt.Type = database.ArtifactTypeLink

			s.db.AddNode(&trgt)
		}
	}
	return &pb.BuildResponse{Success: true}, nil
}

func (s *server) Log(ctx context.Context, in *pb.LogMessage) (*pb.LogResponse, error) {
	log.Printf("REMOTE: %s", string(in.Msg))
	return &pb.LogResponse{Success: true}, nil
}

func (s *server) Quit(ctx context.Context, in *pb.QuitMessage) (*pb.QuitResponse, error) {
	if in.Kill {
		log.Fatalf("qmstr was killed hard by client")
	}

	// Wait for pending tasks to complete e.g. synchronize channels

	// Schedule shutdown
	quitServer <- nil

	return &pb.QuitResponse{Success: true}, nil
}

func ListenAndServe(configfile string) error {
	masterConfig, err := config.ReadConfigFromFile(configfile)
	if err != nil {
		return err
	}

	// Connect to backend database (dgraph)
	db, err := database.Setup(masterConfig.Server.DBAddress, masterConfig.Server.DBWorkers)
	if err != nil {
		return fmt.Errorf("Could not setup database: %v", err)
	}

	analyzerQueue := make(chan analysis.Analysis, 100)
	analysisClosed := make(chan bool)

	// Setup buildservice
	lis, err := net.Listen("tcp", masterConfig.Server.RPCAddress)
	if err != nil {
		return fmt.Errorf("Failed to setup socket and listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBuildServiceServer(s, &server{
		db:             db,
		analyzerQueue:  analyzerQueue,
		serverMutex:    &sync.Mutex{},
		analysisClosed: analysisClosed,
		analysisDone:   false,
		analysis:       masterConfig.Analysis,
		reporting:      masterConfig.Reporting,
	})

	quitServer = make(chan interface{})
	go func() {
		<-quitServer
		log.Println("qmstr-master terminated by client")
		s.GracefulStop()
		close(quitServer)
		quitServer = nil
	}()

	log.Printf("qmstr-master listening on %s\n", masterConfig.Server.RPCAddress)
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("Failed to start rpc service %v", err)
	}
	return nil
}

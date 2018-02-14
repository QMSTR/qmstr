package master

import (
	"fmt"
	"log"
	"net"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/analysis"
	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/database"
	"google.golang.org/grpc"
)

var quitServer chan interface{}

type server struct {
	db            *database.DataBase
	analyzerQueue chan analysis.Analysis
}

func (s *server) Report(ctx context.Context, in *pb.ReportMessage) (*pb.ReportResponse, error) {
	log.Printf("Report requested: %s for %s", in.ReportType, in.Target)
	return &pb.ReportResponse{Success: false}, nil
}

func (s *server) Analyze(ctx context.Context, in *pb.AnalysisMessage) (*pb.AnalysisResponse, error) {
	log.Printf("Analysis requested: %s for %s", in.Analyzer, in.Selector)
	nodeSelector := in.Selector
	analyzerSelector := in.Analyzer

	var analyzer analysis.Analyzer
	switch analyzerSelector {
	case "spdx":
		analyzer = analysis.NewSpdxAnalyzer(in.Config, s.db)
	default:
		return &pb.AnalysisResponse{Success: false}, fmt.Errorf("No such analyzer %s", analyzerSelector)
	}

	nodes, err := s.db.GetNodesByType(nodeSelector)
	if err != nil {
		return &pb.AnalysisResponse{Success: false}, err
	}

	anaNodes := []analysis.AnalysisNode{}
	for _, node := range nodes {
		anaNodes = append(anaNodes, analysis.NewAnalysisNode(node, in.PathSub, s.db))
	}

	s.analyzerQueue <- analysis.Analysis{Nodes: anaNodes, Analyzer: analyzer}

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
			trgt.DerivedFrom = []database.Node{src}
			trgt.Type = database.ArtifactTypeObj

			s.db.AddNode(trgt)
		}
	}

	// Linking
	for _, bin := range in.GetBinary() {
		log.Printf("Linked target: %v", bin)

		uidTrgt, err := s.db.HasNode(bin.Target.GetHash())
		if err != nil {
			return &pb.BuildResponse{Success: false}, err
		}

		deps := []database.Node{}
		// no such node exist
		if uidTrgt == "" {
			for _, dep := range bin.GetInput() {
				depNode := database.NewNode(dep.GetPath(), dep.GetHash())
				depNode.Name = filepath.Base(dep.GetPath())
				deps = append(deps, depNode)
			}
			trgt := database.NewNode(bin.Target.GetPath(), bin.Target.GetHash())
			trgt.DerivedFrom = deps
			trgt.Type = database.ArtifactTypeLink

			s.db.AddNode(trgt)
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

func ListenAndServe(rpcAddr string, dbAddr string) error {

	// Connect to backend database (dgraph)
	db, err := database.Setup(dbAddr)
	if err != nil {
		return fmt.Errorf("Could not setup database: %v", err)
	}

	analyzerQueue := make(chan analysis.Analysis, 100)

	// Setup buildservice
	lis, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		return fmt.Errorf("Failed to setup socket and listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBuildServiceServer(s, &server{
		db:            db,
		analyzerQueue: analyzerQueue,
	})

	quitServer = make(chan interface{})
	go func() {
		<-quitServer
		log.Println("qmstr terminated by client")
		s.GracefulStop()
		close(quitServer)
		quitServer = nil
	}()

	go func() {
		fmt.Println("Analysis queue worker started")
		for ana := range analyzerQueue {
			analysis.RunAnalysis(ana)
		}

	}()

	log.Print("qmstr master running")
	if err := s.Serve(lis); err != nil {
		return fmt.Errorf("Failed to start rpc service %v", err)
	}
	return nil
}

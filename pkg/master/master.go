package master

import (
	"fmt"
	"log"
	"net"
	"sync"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	pb "github.com/QMSTR/qmstr/pkg/service"
	"google.golang.org/grpc"
)

var quitServer chan interface{}

type serverPhase interface {
	Build(*pb.BuildMessage) (*pb.BuildResponse, error)
	GetNodes(*pb.NodeRequest) (*pb.NodeResponse, error)
	SendNodes(*pb.AnalysisMessage) (*pb.AnalysisResponse, error)
	Report(*pb.ReportRequest, pb.ReportService_ReportServer) error
}

type genericServerPhase struct {
	Name string
}

type server struct {
	db             *database.DataBase
	analysisClosed chan bool
	serverMutex    *sync.Mutex
	analysisDone   bool
	analysis       []config.Analysis
	reporting      []config.Reporting
	currentPhase   serverPhase
}

func (s *server) Build(ctx context.Context, in *pb.BuildMessage) (*pb.BuildResponse, error) {
	return s.currentPhase.Build(in)
}

func (s *server) GetNodes(ctx context.Context, in *pb.NodeRequest) (*pb.NodeResponse, error) {
	return s.currentPhase.GetNodes(in)
}

func (s *server) SendNodes(ctx context.Context, in *pb.AnalysisMessage) (*pb.AnalysisResponse, error) {
	return s.currentPhase.SendNodes(in)
}

func (s *server) Report(in *pb.ReportRequest, streamServer pb.ReportService_ReportServer) error {
	return s.currentPhase.Report(in, streamServer)
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

	analysisClosed := make(chan bool)

	// Setup buildservice
	lis, err := net.Listen("tcp", masterConfig.Server.RPCAddress)
	if err != nil {
		return fmt.Errorf("Failed to setup socket and listen: %v", err)
	}
	s := grpc.NewServer()
	serverImpl := &server{
		db:             db,
		serverMutex:    &sync.Mutex{},
		analysisClosed: analysisClosed,
		analysisDone:   false,
		analysis:       masterConfig.Analysis,
		reporting:      masterConfig.Reporting,
	}
	pb.RegisterBuildServiceServer(s, serverImpl)
	pb.RegisterAnalysisServiceServer(s, serverImpl)
	pb.RegisterReportServiceServer(s, serverImpl)
	pb.RegisterControlServiceServer(s, serverImpl)

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

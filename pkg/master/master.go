package master

import (
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
	"google.golang.org/grpc"
)

var quitServer chan interface{}
var phaseMap map[int32]serverPhase

type serverPhase interface {
	GetPhaseId() int32
	Activate() bool
	Build(*service.BuildMessage) (*service.BuildResponse, error)
	GetNodes(*service.NodeRequest) (*service.NodeResponse, error)
	SendNodes(*service.AnalysisMessage) (*service.AnalysisResponse, error)
	Report(*service.ReportRequest, service.ReportService_ReportServer) error
}

type genericServerPhase struct {
	Name    string
	phaseId int32
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

func (s *server) Build(ctx context.Context, in *service.BuildMessage) (*service.BuildResponse, error) {
	return s.currentPhase.Build(in)
}

func (s *server) GetNodes(ctx context.Context, in *service.NodeRequest) (*service.NodeResponse, error) {
	return s.currentPhase.GetNodes(in)
}

func (s *server) SendNodes(ctx context.Context, in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	return s.currentPhase.SendNodes(in)
}

func (s *server) Report(in *service.ReportRequest, streamServer service.ReportService_ReportServer) error {
	return s.currentPhase.Report(in, streamServer)
}

func (s *server) SwitchPhase(ctx context.Context, in *service.SwitchPhaseMessage) (*service.SwitchPhaseResponse, error) {
	requestedPhase := in.Phase
	if requestedPhase <= s.currentPhase.GetPhaseId() {
		errMsg := fmt.Sprintf("Illegal phase transition %d->%d requested", s.currentPhase.GetPhaseId(), requestedPhase)
		log.Println(errMsg)
		return &service.SwitchPhaseResponse{Success: false}, errors.New(errMsg)
	}
	if phase, ok := phaseMap[requestedPhase]; ok {
		log.Printf("Switching to %d phase", requestedPhase)
		s.currentPhase = phase
		s.currentPhase.Activate()
		return &service.SwitchPhaseResponse{Success: true}, nil
	}
	return &service.SwitchPhaseResponse{Success: false}, fmt.Errorf("Invalid phase requested %d", requestedPhase)
}

func (s *server) Log(ctx context.Context, in *service.LogMessage) (*service.LogResponse, error) {
	log.Printf("REMOTE: %s", string(in.Msg))
	return &service.LogResponse{Success: true}, nil
}

func (s *server) Quit(ctx context.Context, in *service.QuitMessage) (*service.QuitResponse, error) {
	if in.Kill {
		log.Fatalf("qmstr was killed hard by client")
	}

	// Wait for pending tasks to complete e.g. synchronize channels

	// Schedule shutdown
	quitServer <- nil

	return &service.QuitResponse{Success: true}, nil
}

func InitAndRun(configfile string) error {
	masterConfig, err := config.ReadConfigFromFile(configfile)
	if err != nil {
		return err
	}

	phaseMap = map[int32]serverPhase{
		1: &serverPhaseBuild{genericServerPhase{Name: "Build phase", phaseId: 1}},
		2: &serverPhaseAnalysis{genericServerPhase{Name: "Analysis phase", phaseId: 2}},
		3: &serverPhaseReport{genericServerPhase{Name: "Reporting phase", phaseId: 3}},
	}

	// Connect to backend database (dgraph)
	db, err := database.Setup(masterConfig.Server.DBAddress, masterConfig.Server.DBWorkers)
	if err != nil {
		return fmt.Errorf("Could not setup database: %v", err)
	}

	// Setup buildservice
	lis, err := net.Listen("tcp", masterConfig.Server.RPCAddress)
	if err != nil {
		return fmt.Errorf("Failed to setup socket and listen: %v", err)
	}
	s := grpc.NewServer()
	serverImpl := &server{
		db:             db,
		serverMutex:    &sync.Mutex{},
		analysisClosed: make(chan bool),
		analysisDone:   false,
		analysis:       masterConfig.Analysis,
		reporting:      masterConfig.Reporting,
	}
	service.RegisterBuildServiceServer(s, serverImpl)
	service.RegisterAnalysisServiceServer(s, serverImpl)
	service.RegisterReportServiceServer(s, serverImpl)
	service.RegisterControlServiceServer(s, serverImpl)

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

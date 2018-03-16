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
	Activate() error
	Shutdown() error
	Build(*service.BuildMessage) (*service.BuildResponse, error)
	GetAnalyzerConfig(*service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error)
	GetReporterConfig(*service.ReporterConfigRequest) (*service.ReporterConfigResponse, error)
	GetNodes(*service.NodeRequest) (*service.NodeResponse, error)
	SendNodes(*service.AnalysisMessage) (*service.AnalysisResponse, error)
	GetReportNodes(*service.ReportRequest, service.ReportService_GetReportNodesServer) error
}

type genericServerPhase struct {
	Name       string
	phaseId    int32
	db         *database.DataBase
	rpcAddress string
}

type server struct {
	db             *database.DataBase
	analysisClosed chan bool
	serverMutex    *sync.Mutex
	analysisDone   bool
	currentPhase   serverPhase
}

func (s *server) Build(ctx context.Context, in *service.BuildMessage) (*service.BuildResponse, error) {
	return s.currentPhase.Build(in)
}

func (s *server) GetAnalyzerConfig(ctx context.Context, in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	return s.currentPhase.GetAnalyzerConfig(in)
}

func (s *server) GetReporterConfig(ctx context.Context, in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	return s.currentPhase.GetReporterConfig(in)
}

func (s *server) GetNodes(ctx context.Context, in *service.NodeRequest) (*service.NodeResponse, error) {
	return s.currentPhase.GetNodes(in)
}

func (s *server) SendNodes(ctx context.Context, in *service.AnalysisMessage) (*service.AnalysisResponse, error) {
	return s.currentPhase.SendNodes(in)
}

func (s *server) GetReportNodes(in *service.ReportRequest, streamServer service.ReportService_GetReportNodesServer) error {
	return s.currentPhase.GetReportNodes(in, streamServer)
}

func (s *server) SwitchPhase(ctx context.Context, in *service.SwitchPhaseMessage) (*service.SwitchPhaseResponse, error) {
	requestedPhase := in.Phase
	if requestedPhase <= s.currentPhase.GetPhaseId() {
		errMsg := fmt.Sprintf("Illegal phase transition %d->%d requested", s.currentPhase.GetPhaseId(), requestedPhase)
		log.Println(errMsg)
		return &service.SwitchPhaseResponse{Success: false}, errors.New(errMsg)
	}
	if phase, ok := phaseMap[requestedPhase]; ok {
		log.Printf("Switching to phase %d", requestedPhase)
		s.currentPhase.Shutdown()
		s.currentPhase = phase
		err := s.currentPhase.Activate()
		if err != nil {
			return &service.SwitchPhaseResponse{Success: false}, err
		}
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

	phaseMap = map[int32]serverPhase{
		1: &serverPhaseBuild{genericServerPhase{Name: "Build phase", phaseId: 1, db: db}},
		2: newAnalysisPhase(genericServerPhase{Name: "Analysis phase", phaseId: 2, db: db, rpcAddress: masterConfig.Server.RPCAddress},
			masterConfig.Analysis),
		3: &serverPhaseReport{genericServerPhase{Name: "Reporting phase", phaseId: 3, db: db, rpcAddress: masterConfig.Server.RPCAddress}, masterConfig.Reporting},
	}

	s := grpc.NewServer()
	serverImpl := &server{
		db:             db,
		serverMutex:    &sync.Mutex{},
		analysisClosed: make(chan bool),
		analysisDone:   false,
		currentPhase:   phaseMap[1],
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

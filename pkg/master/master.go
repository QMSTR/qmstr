package master

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/net/context"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
	"google.golang.org/grpc"
)

var quitServer chan interface{}
var phaseMap map[int32]func(string, *config.MasterConfig, *database.DataBase) serverPhase

func init() {
	phaseMap = map[int32]func(string, *config.MasterConfig, *database.DataBase) serverPhase{
		PhaseIDBuild:    newBuildPhase,
		PhaseIDAnalysis: newAnalysisPhase,
		PhaseIDReport:   newReportPhase,
	}
}

type serverPhase interface {
	GetPhaseID() int32
	getName() string
	Activate() error
	Shutdown() error
	getDataBase() (*database.DataBase, error)
	getSession() string
	getMasterConfig() *config.MasterConfig
	Build(*service.BuildMessage) (*service.BuildResponse, error)
	GetAnalyzerConfig(*service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error)
	GetReporterConfig(*service.ReporterConfigRequest) (*service.ReporterConfigResponse, error)
	GetNodes(*service.NodeRequest) (*service.NodeResponse, error)
	SendNodes(*service.AnalysisMessage) (*service.AnalysisResponse, error)
}

type genericServerPhase struct {
	Name         string
	db           *database.DataBase
	session      string
	masterConfig *config.MasterConfig
}

func (gsp *genericServerPhase) getDataBase() (*database.DataBase, error) {
	if gsp.db == nil {
		return nil, errors.New("Database not yet available")
	}
	return gsp.db, nil
}

func (gsp *genericServerPhase) getSession() string {
	return gsp.session
}

func (gsp *genericServerPhase) getMasterConfig() *config.MasterConfig {
	return gsp.masterConfig
}

func (gsp *genericServerPhase) getName() string {
	return gsp.Name
}

type server struct {
	analysisClosed     chan bool
	serverMutex        *sync.Mutex
	analysisDone       bool
	currentPhase       serverPhase
	pendingPhaseSwitch int64
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

func (s *server) GetPackageNode(ctx context.Context, in *service.PackageRequest) (*service.PackageResponse, error) {
	db, err := s.currentPhase.getDataBase()
	if err != nil {
		return nil, err
	}
	node, err := db.GetPackageNode(in.Session)
	if err != nil {
		return nil, err
	}
	return &service.PackageResponse{PackageNode: node}, nil
}

func (s *server) Status(ctx context.Context, in *service.StatusMessage) (*service.StatusResponse, error) {
	resp := service.StatusResponse{}
	resp.PhaseID = s.currentPhase.GetPhaseID()
	if in.Phase {
		resp.Phase = s.currentPhase.getName()
	}
	if in.Switch {
		resp.Switching = atomic.LoadInt64(&s.pendingPhaseSwitch) == 1
	}
	return &resp, nil
}

func (s *server) SwitchPhase(ctx context.Context, in *service.SwitchPhaseMessage) (*service.SwitchPhaseResponse, error) {
	requestedPhase := in.Phase
	err := s.switchPhase(requestedPhase)
	if err != nil {
		return &service.SwitchPhaseResponse{Success: false}, err
	}
	return &service.SwitchPhaseResponse{Success: true}, nil
}

func (s *server) switchPhase(requestedPhase int32) error {
	if !atomic.CompareAndSwapInt64(&s.pendingPhaseSwitch, 0, 1) {
		errMsg := "denied there is a pending phase transition"
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if requestedPhase <= s.currentPhase.GetPhaseID() {
		errMsg := fmt.Sprintf("Illegal phase transition %d->%d requested", s.currentPhase.GetPhaseID(), requestedPhase)
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if phaseCtor, ok := phaseMap[requestedPhase]; ok {
		log.Printf("Switching to phase %d", requestedPhase)
		s.currentPhase.Shutdown()
		db, err := s.currentPhase.getDataBase()
		if err != nil {
			return err
		}
		s.currentPhase = phaseCtor(s.currentPhase.getSession(), s.currentPhase.getMasterConfig(), db)
		s.pendingPhaseSwitch = 0
		err = s.currentPhase.Activate()
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("Invalid phase requested %d", requestedPhase)
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

// InitAndRun sets up and runs the grpc services and the dgraph database connection
func InitAndRun(masterConfig *config.MasterConfig) (chan error, error) {
	// Setup buildservice
	lis, err := net.Listen("tcp", masterConfig.Server.RPCAddress)
	if err != nil {
		return nil, fmt.Errorf("Failed to setup socket and listen: %v", err)
	}

	sessionBytes := make([]byte, 32)
	rand.Read(sessionBytes)
	session := fmt.Sprintf("%x", sessionBytes)

	s := grpc.NewServer()
	serverImpl := &server{
		serverMutex:    &sync.Mutex{},
		analysisClosed: make(chan bool),
		analysisDone:   false,
		currentPhase:   newInitServerPhase(session, masterConfig),
	}
	service.RegisterBuildServiceServer(s, serverImpl)
	service.RegisterAnalysisServiceServer(s, serverImpl)
	service.RegisterReportServiceServer(s, serverImpl)
	service.RegisterControlServiceServer(s, serverImpl)

	// start grpc service as soon as possible to allow clients to connect and get status feedback
	masterRun := make(chan error)
	go func() {
		log.Printf("qmstr-master listening on %s\n", masterConfig.Server.RPCAddress)
		err := s.Serve(lis)
		masterRun <- err
	}()

	// Activate init phase
	err = serverImpl.currentPhase.Activate()
	if err != nil {
		masterRun <- err
		return nil, err
	}

	serverImpl.switchPhase(PhaseIDBuild)

	quitServer = make(chan interface{})
	go func() {
		<-quitServer
		log.Println("qmstr-master terminated by client")
		s.GracefulStop()
		close(quitServer)
		quitServer = nil
	}()

	return masterRun, nil
}

func logModuleError(moduleName string, output []byte) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("%s failed with:\n", moduleName))
	s := bufio.NewScanner(strings.NewReader(string(output)))
	for s.Scan() {
		buffer.WriteString(fmt.Sprintf("\t--> %s\n", s.Text()))
	}
	log.Println(buffer.String())
}

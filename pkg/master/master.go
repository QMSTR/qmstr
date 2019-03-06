package master

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
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
var phaseMap map[service.Phase]func(*config.MasterConfig, *database.DataBase, *server) serverPhase

func init() {
	phaseMap = map[service.Phase]func(*config.MasterConfig, *database.DataBase, *server) serverPhase{
		service.Phase_BUILD:    newBuildPhase,
		service.Phase_ANALYSIS: newAnalysisPhase,
		service.Phase_REPORT:   newReportPhase,
	}
}

type server struct {
	analysisClosed     chan bool
	serverMutex        *sync.Mutex
	analysisDone       bool
	currentPhase       serverPhase
	pendingPhaseSwitch int64
	eventChannels      map[service.EventClass][]chan *service.Event
	eventMutex         *sync.RWMutex
}

func (s *server) Build(stream service.BuildService_BuildServer) error {
	return s.currentPhase.Build(stream)
}

func (s *server) PushFile(ctx context.Context, in *service.PushFileMessage) (*service.PushFileResponse, error) {
	return s.currentPhase.PushFile(in)
}

func (s *server) CreatePackage(ctx context.Context, in *service.PackageNode) (*service.BuildResponse, error) {
	return s.currentPhase.CreatePackage(in)
}

func (s *server) CreateProject(ctx context.Context, in *service.ProjectNode) (*service.BuildResponse, error) {
	return s.currentPhase.CreateProject(in)
}

func (s *server) GetAnalyzerConfig(ctx context.Context, in *service.AnalyzerConfigRequest) (*service.AnalyzerConfigResponse, error) {
	return s.currentPhase.GetAnalyzerConfig(in)
}

func (s *server) GetReporterConfig(ctx context.Context, in *service.ReporterConfigRequest) (*service.ReporterConfigResponse, error) {
	return s.currentPhase.GetReporterConfig(in)
}

func (s *server) SendInfoNodes(stream service.AnalysisService_SendInfoNodesServer) error {
	return s.currentPhase.SendInfoNodes(stream)
}

func (s *server) SendFileNode(stream service.AnalysisService_SendFileNodeServer) error {
	return s.currentPhase.SendFileNode(stream)
}

func (s *server) SendPackageNode(stream service.AnalysisService_SendPackageNodeServer) error {
	return s.currentPhase.SendPackageNode(stream)
}

func (s *server) SendDiagnosticNode(stream service.AnalysisService_SendDiagnosticNodeServer) error {
	return s.currentPhase.SendDiagnosticNode(stream)
}

func (s *server) GetBOM(ctx context.Context, in *service.BOMRequest) (*service.BOM, error) {
	return s.currentPhase.GetBOM(in)
}

func (s *server) Package(stream service.BuildService_PackageServer) error {
	db, err := s.currentPhase.getDataBase()
	if err != nil {
		return err
	}
	pkg, err := db.GetPackageNode()
	if err != nil {
		return err
	}
	for {
		fl, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("recv fileNode fail: %v", err)
		}
		pkg.Targets = append(pkg.Targets, fl)
	}
	db.AddPackageNode(pkg)
	return stream.SendAndClose(&service.BuildResponse{Success: true})
}

func (s *server) GetPackageNode(ctx context.Context, in *service.PackageRequest) (*service.PackageNode, error) {
	db, err := s.currentPhase.getDataBase()
	if err != nil {
		return nil, err
	}
	node, err := db.GetPackageNode()
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (s *server) SendBuildError(ctx context.Context, in *service.InfoNode) (*service.BuildResponse, error) {
	db, err := s.currentPhase.getDataBase()
	if err != nil {
		return nil, err
	}
	node, err := db.GetPackageNode()
	if err != nil {
		return nil, err
	}

	if err := db.AddInfoNodes(node.Uid, in); err != nil {
		return nil, err
	}

	return &service.BuildResponse{Success: true}, nil
}

func (s *server) GetFileNode(in *service.FileNode, stream service.ControlService_GetFileNodeServer) error {
	return s.currentPhase.GetFileNode(in, stream)
}

func (s *server) GetDiagnosticNode(in *service.DiagnosticNode, stream service.ControlService_GetDiagnosticNodeServer) error {
	return s.currentPhase.GetDiagnosticNode(in, stream)
}

func (s *server) GetInfoData(ctx context.Context, in *service.InfoDataRequest) (*service.InfoDataResponse, error) {
	return s.currentPhase.GetInfoData(in)
}

func (s *server) ExportSnapshot(ctx context.Context, in *service.ExportRequest) (*service.ExportResponse, error) {
	return s.currentPhase.ExportSnapshot(in)
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
	resp.Error = s.currentPhase.getError()
	return &resp, nil
}

func (s *server) SwitchPhase(ctx context.Context, in *service.SwitchPhaseMessage) (*service.SwitchPhaseResponse, error) {
	requestedPhase := in.Phase
	err := s.switchPhase(requestedPhase)
	if err != nil {
		return &service.SwitchPhaseResponse{Success: false, Error: err.Error()}, nil
	}
	return &service.SwitchPhaseResponse{Success: true}, nil
}

func (s *server) switchPhase(requestedPhase service.Phase) error {
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
		defer s.persistPhase()
		s.publishEvent(&service.Event{Class: service.EventClass_PHASE, Message: fmt.Sprintf("Switching to phase %d", requestedPhase)})
		err := s.currentPhase.Shutdown()
		if err != nil {
			// switch to failure phase
			s.enterFailureServerPhase(err)
			return err
		}
		db, err := s.currentPhase.getDataBase()
		if err != nil {
			// switch to failure phase
			s.enterFailureServerPhase(err)
			return err
		}
		s.currentPhase = phaseCtor(s.currentPhase.getMasterConfig(), db, s)
		s.pendingPhaseSwitch = 0
		err = s.currentPhase.Activate()
		if err != nil {
			s.publishEvent(&service.Event{Class: service.EventClass_PHASE, Message: "Entering failure phase"})
			s.enterFailureServerPhase(err)
			return err
		}
		s.publishEvent(&service.Event{Class: service.EventClass_PHASE, Message: fmt.Sprintf("Switched to phase %d", requestedPhase)})
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

	s := grpc.NewServer()
	serverImpl := &server{
		serverMutex:    &sync.Mutex{},
		analysisClosed: make(chan bool),
		analysisDone:   false,
		currentPhase:   newInitServerPhase(masterConfig),
		eventChannels: map[service.EventClass][]chan *service.Event{
			service.EventClass_ALL:    []chan *service.Event{},
			service.EventClass_MODULE: []chan *service.Event{},
			service.EventClass_PHASE:  []chan *service.Event{},
		},
		eventMutex: &sync.RWMutex{},
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
		serverImpl.enterFailureServerPhase(err)
		return masterRun, err
	}

	serverImpl.switchPhase(serverImpl.currentPhase.getPostInitPhase())

	quitServer = make(chan interface{})
	go func() {
		<-quitServer
		log.Println("qmstr-master terminated by client")
		s.GracefulStop()
		close(quitServer)
		masterRun <- nil
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

func (s *server) persistPhase() error {
	if s.pendingPhaseSwitch != 0 {
		return errors.New("Can not persist phase while switching")
	}
	db, err := s.currentPhase.getDataBase()
	if err != nil {
		return err
	}

	if _, err := db.AddQmstrStateNode(&service.QmstrStateNode{Phase: s.currentPhase.GetPhaseID()}); err != nil {
		return err
	}
	return nil
}

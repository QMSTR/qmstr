//go:generate protoc -I ../../pkg/service --go_out=plugins=grpc:../../pkg/service ../../pkg/service/datamodel.proto ../../pkg/service/analyzerservice.proto ../../pkg/service/buildservice.proto ../../pkg/service/controlservice.proto  ../../pkg/service/reportservice.proto
package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/compiler"
	pb "github.com/QMSTR/qmstr/pkg/service"
	"github.com/QMSTR/qmstr/pkg/wrapper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	addrEnv  = "QMSTR_ADDRESS"
	debugEnv = "QMSTR_DEBUG"
)

var (
	buildServiceClient   pb.BuildServiceClient
	controlServiceClient pb.ControlServiceClient
	logger               *log.Logger
	conn                 *grpc.ClientConn
	debug                bool
)

var address = "localhost:50051"

func initLogging() {
	var infoWriter io.Writer
	infoWriter = wrapper.NewRemoteLogWriter(controlServiceClient)
	logger = log.New(infoWriter, "", log.Ldate|log.Ltime)
}

func main() {
	_, debug = os.LookupEnv(debugEnv)
	_, difAddress := os.LookupEnv(addrEnv)

	if difAddress {
		address = os.Getenv(addrEnv)
	}
	// Set up server connection
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()
	buildServiceClient = pb.NewBuildServiceClient(conn)
	controlServiceClient = pb.NewControlServiceClient(conn)

	initLogging()

	commandLine := os.Args
	if debug {
		logger.Printf("QMSTR called via %v", commandLine)
	}

	// find out where we are
	workingDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Could not get current working dir.")
	}
	if debug {
		logger.Printf("Wrapper running in %s", workingDir)
	}

	w := wrapper.NewWrapper(commandLine, logger, debug)
	w.Wrap()
	compiler := compiler.GetCompiler(w.Program, workingDir, logger, debug)
	buildMsg, err := compiler.Analyze(commandLine)
	if err == nil {
		send_result(buildMsg)
	}
}

func send_result(buildmsg *pb.BuildMessage) error {
	r, err := buildServiceClient.Build(context.Background(), buildmsg)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if !r.Success {
		return errors.New("Server failure")
	}
	return nil
}

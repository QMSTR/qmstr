//go:generate protoc -I ../../proto --go_out=plugins=grpc:../../pkg/service ../../proto/datamodel.proto ../../proto/analyzerservice.proto ../../proto/buildservice.proto ../../proto/controlservice.proto  ../../proto/reportservice.proto
package main

import (
	"errors"
	"io"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/builder"
	pb "github.com/QMSTR/qmstr/pkg/service"
	"github.com/QMSTR/qmstr/pkg/wrapper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	addrEnv  = "QMSTR_MASTER"
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
	build := builder.GetBuilder(w.Program, workingDir, logger, debug)
	buildMsg, err := build.Analyze(commandLine)
	if err == nil {
		sendResult(buildMsg)
	}
}

func sendResult(buildmsg *pb.BuildMessage) error {
	r, err := buildServiceClient.Build(context.Background(), buildmsg)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if !r.Success {
		return errors.New("Server failure")
	}
	return nil
}

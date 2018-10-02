package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/logging"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/QMSTR/qmstr/pkg/wrapper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	buildServiceClient   service.BuildServiceClient
	controlServiceClient service.ControlServiceClient
	logger               *log.Logger
	conn                 *grpc.ClientConn
	debug                bool
)

var address = "localhost:50051"

func initLogging() {
	var infoWriter io.Writer
	infoWriter = logging.NewRemoteLogWriter(controlServiceClient)
	logger = log.New(infoWriter, "", log.Ldate|log.Ltime)
}

func main() {
	_, debug = os.LookupEnv(common.QMSTRDEBUGENV)
	_, difAddress := os.LookupEnv(common.QMSTRADDRENV)

	if difAddress {
		address = os.Getenv(common.QMSTRADDRENV)
	}
	// Set up server connection
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	defer conn.Close()
	buildServiceClient = service.NewBuildServiceClient(conn)
	controlServiceClient = service.NewControlServiceClient(conn)

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

	w, err := wrapper.NewWrapper(commandLine, workingDir, logger, debug)
	if err != nil {
		log.Fatalf("failed to create wrapper for %s: %v", commandLine, err)
	}
	w.Wrap()
	buildMsg, err := w.Builder.Analyze(commandLine)
	if err == nil {
		if buildMsg != nil {
			sendResult(buildMsg)
		}
	} else {
		logger.Printf("%s failed for \"%s\": %v", w.Builder.GetName(), commandLine, err)
		sendBuildError(fmt.Sprintf("Failed to analyze build [%s] due to %v", commandLine, err))
	}
}

func sendBuildError(errMsg string) error {
	errNode := service.CreateErrorNode(errMsg)
	r, err := buildServiceClient.SendBuildError(context.Background(), errNode)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if !r.Success {
		return errors.New("Server failure")
	}
	return nil
}

func sendResult(buildmsg *service.BuildMessage) error {
	r, err := buildServiceClient.Build(context.Background(), buildmsg)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if !r.Success {
		return errors.New("Server failure")
	}
	return nil
}

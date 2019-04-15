package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/logging"
	"github.com/QMSTR/qmstr/pkg/service"
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
	defer w.Exit()
	w.Wrap()

	fileNodes, err := w.Builder.Analyze(commandLine)
	switch err {
	case nil:
		stream, err := buildServiceClient.Build(context.Background())
		defer func() {
			res, err := stream.CloseAndRecv()
			if err != nil {
				log.Fatalf("Failed to close the filenode stream: %v", err)
			}
			if !res.Success {
				log.Fatalln("Server filenode stream failed")
			}
		}()
		if err != nil {
			log.Fatalf("could not greet: %v", err)
		}
		for _, fileNode := range fileNodes {
			pushFileMsg, err := w.Builder.GetPushFile()
			if err != nil {
				if err != builder.ErrNoPushFile {
					errMsg := fmt.Sprintf("%s failed to get file for upload: %v", w.Builder.GetName(), err)
					sendBuildException(service.ExceptionType_ERROR, errMsg)
					logger.Println(errMsg)
				}
			}
			if pushFileMsg != nil {
				remotePath, err := pushFile(pushFileMsg)
				if err != nil {
					errMsg := fmt.Sprintf("%s failed to upload file", pushFileMsg.Hash)
					sendBuildException(service.ExceptionType_ERROR, errMsg)
					logger.Println(errMsg)
				}
				for _, dep := range fileNode.DerivedFrom {
					if dep.Name == "-" {
						dep.Hash = pushFileMsg.Hash
						dep.Path = remotePath
					}
				}
			}
			if err := stream.Send(fileNode); err != nil {
				log.Fatalf("Failed to send filenode to server")
			}
		}
	case builder.ErrBuilderModeNotImplemented:
		logger.Printf("WARNING for %s: \"%s\": %v", w.Builder.GetName(), commandLine, err)
		sendBuildException(service.ExceptionType_WARNING, fmt.Sprintf("Warning while analyzing [%s]: %v", commandLine, err))
	default:
		logger.Printf("%s failed for \"%s\": %v", w.Builder.GetName(), commandLine, err)
		sendBuildException(service.ExceptionType_ERROR, fmt.Sprintf("Failed to analyze build [%s] due to %v", commandLine, err))
	}

}

func pushFile(pushMsg *service.PushFileMessage) (string, error) {
	r, err := buildServiceClient.PushFile(context.Background(), pushMsg)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	return r.Path, nil
}

func sendBuildException(exType service.ExceptionType, msg string) error {
	var exNode *service.InfoNode
	switch exType {
	case service.ExceptionType_ERROR:
		exNode = service.CreateErrorNode(msg)
	case service.ExceptionType_WARNING:
		exNode = service.CreateWarningNode(msg)
	}
	r, err := buildServiceClient.SendBuildError(context.Background(), exNode)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if !r.Success {
		return errors.New("Server failure")
	}
	return nil
}

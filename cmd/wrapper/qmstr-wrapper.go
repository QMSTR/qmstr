//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto

package main

import (
	"errors"
	"log"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var buildServiceClient pb.BuildServiceClient

func main() {
	// Set up server connection
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	buildServiceClient = pb.NewBuildServiceClient(conn)

}

func send_result(buildmsg pb.BuildMessage) error {
	r, err := buildServiceClient.Build(context.Background(), &buildmsg)
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	if !r.Cool {
		return errors.New("Server failure")
	}
	return nil
}

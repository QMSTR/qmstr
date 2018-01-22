//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto

package main

import (
	"log"
	"net"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var quitServer chan interface{}

type server struct{}

func (s *server) Build(ctx context.Context, in *pb.BuildMessage) (*pb.BuildResponse, error) {
	// TODO implement handling
	return &pb.BuildResponse{Success: true}, nil
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

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to setup socket and listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBuildServiceServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)

	quitServer = make(chan interface{})
	go func() {
		<-quitServer
		log.Println("qmstr terminated by client")
		s.GracefulStop()
		close(quitServer)
		quitServer = nil
	}()

	log.Print("About to serve")

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

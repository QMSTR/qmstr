//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto

package main

import (
	"log"
	"net"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/QMSTR/qmstr/pkg/dgraph"
	"github.com/dgraph-io/dgraph/client"
	"github.com/dgraph-io/dgraph/protos/api"
)

const (
	serverPort = ":50051"
	clientPort = "localhost:9080"
)

var cl *client.Dgraph
var quitServer chan interface{}

type server struct{}

func (s *server) Build(ctx context.Context, in *pb.BuildMessage) (*pb.BuildResponse, error) {
	for _, bin := range in.Binary {
		log.Printf("Linked target: %v", bin)
	}
	for _, compile := range in.Compilations {
		log.Printf("Compiled %v", compile)
	}
	dgraph.AddData(cl, in)
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
	lis, err := net.Listen("tcp", serverPort)
	if err != nil {
		log.Fatalf("Failed to setup socket and listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBuildServiceServer(s, &server{})

	// Create a client connection
	conn, err := grpc.Dial(clientPort, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to the dgraph server: %v", err)
	}
	defer conn.Close()
	cl = client.NewDgraphClient(api.NewDgraphClient(conn))
	// Set up our dgraph schema
	dgraph.Setup(cl)

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

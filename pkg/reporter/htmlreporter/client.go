package htmlreporter

import (
	"log"

	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

var conn *grpc.ClientConn

//var buildServiceClient pb.BuildServiceClient

// ConnectToMaster connects to the QMSTR master process.
func ConnectToMaster() {
	// Set up server connection
	var err error
	conn, err = grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}
	//buildServiceClient = pb.NewBuildServiceClient(conn)
}

// Temp blabla
func Temp() {
	//buildServiceClient.Log(context.Background(), &pb.LogMessage{Msg: []byte("HTML Reporter starting...")})
}

// DisconnectFromMaster disconnects from the QMSTR master process.
func DisconnectFromMaster() {
	//conn.Close()
}

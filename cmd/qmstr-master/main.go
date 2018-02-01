//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto

package main

import (
	"log"

	"github.com/QMSTR/qmstr/pkg/master"
	flag "github.com/spf13/pflag"
)

func main() {

	rpcAddr := flag.String("rpcAddr", ":50051", "Set the address and port to bind to.")
	dbAddr := flag.String("dbAddr", "localhost:9080", "Set the address and port of the backing database.")

	flag.Parse()

	if err := master.ListenAndServe(*rpcAddr, *dbAddr); err != nil {
		log.Fatalf("Cannot start QMSTR Server: %v\n", err)
	}
}

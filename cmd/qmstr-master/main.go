//go:generate protoc -I ../../pkg/service --go_out=plugins=grpc:../../pkg/service ../../pkg/service/*.proto

package main

import (
	"log"

	"github.com/QMSTR/qmstr/pkg/master"
	flag "github.com/spf13/pflag"
)

func main() {

	configFile := flag.String("config", "qmstr.yaml", "Set the qmstr configuration file.")
	flag.Parse()

	if err := master.ListenAndServe(*configFile); err != nil {
		log.Fatalf("Cannot start QMSTR Server: %v\n", err)
	}
}

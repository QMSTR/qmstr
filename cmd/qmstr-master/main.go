//go:generate protoc -I ../../proto --go_out=plugins=grpc:../../pkg/service ../../proto/datamodel.proto ../../proto/analyzerservice.proto ../../proto/buildservice.proto ../../proto/controlservice.proto  ../../proto/reportservice.proto

package main

import (
	"log"

	"github.com/QMSTR/qmstr/pkg/master"
	flag "github.com/spf13/pflag"
)

func main() {

	configFile := flag.String("config", "qmstr.yaml", "Set the qmstr configuration file.")
	flag.Parse()

	if err := master.InitAndRun(*configFile); err != nil {
		log.Fatalf("Cannot start QMSTR Server: %v\n", err)
	}
}

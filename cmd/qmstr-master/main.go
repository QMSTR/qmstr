package main

import (
	"log"

	"github.com/QMSTR/qmstr/pkg/master"
	flag "github.com/spf13/pflag"
)

func main() {
	configFile := flag.String("config", "qmstr.yaml", "Set the qmstr configuration file.")
	flag.Parse()

	masterRun, err := master.InitAndRun(*configFile)
	if err != nil {
		log.Fatalf("Cannot start QMSTR Server: %v\n", err)
	}

	err = <-masterRun
	if err != nil {
		log.Fatalf("QMSTR master failed: %v\n", err)
	}
}

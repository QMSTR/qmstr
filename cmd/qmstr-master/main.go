package main

import (
	"log"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/master"
	flag "github.com/spf13/pflag"
)

func main() {
	configFile := flag.String("config", "qmstr.yaml", "Set the qmstr configuration file.")
	flag.Parse()

	masterConfig, err := config.ReadConfigFromFile(*configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration %v", err)
	}

	masterRun, err := master.InitAndRun(&masterConfig)
	if err != nil {
		log.Fatalf("Cannot start QMSTR Server: %v\n", err)
	}

	err = <-masterRun
	if err != nil {
		log.Fatalf("QMSTR master failed: %v\n", err)
	}
}

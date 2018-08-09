package main

import (
	"log"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/service"
	flag "github.com/spf13/pflag"
)

func main() {
	configFile := flag.String("config", "qmstr.yaml", "Set the qmstr configuration file.")
	var pathSubstitution []string
	flag.StringSliceVar(&pathSubstitution, "pathsub", nil, "Set path substitution e.g. old,new")
	flag.Parse()

	masterConfig, err := config.ReadConfigFromFiles(*configFile)
	if err != nil {
		log.Fatalf("Failed to read configuration %v", err)
	}

	if pathSubstitution != nil {
		if len(pathSubstitution)%2 != 0 {
			log.Fatalln("Path substitution provided via commandline is invalid")
		}
		for i := 0; i < len(pathSubstitution); i += 2 {
			masterConfig.Server.PathSub = append(masterConfig.Server.PathSub, &service.PathSubstitution{Old: pathSubstitution[i], New: pathSubstitution[i+1]})
		}
		log.Printf("Using following PathSubstitution: %v", masterConfig.Server.PathSub)
	}

	masterRun, err := master.InitAndRun(masterConfig)
	if err != nil {
		log.Printf("Cannot start QMSTR Server: %v\n", err)
	}

	err = <-masterRun
	if err != nil {
		log.Fatalf("QMSTR master failed: %v\n", err)
	}
	log.Println("QMSTR master quit")
}

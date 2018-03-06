package config

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/service"
	yaml "gopkg.in/yaml.v2"
)

type Analysis struct {
	Selector string
	Analyzer string
	PathSub  []*service.PathSubstitution
	Config   map[string]string
}

type Reporting struct {
	Selector   string
	Name       string
	ReportType string
}

type ServerConfig struct {
	RPCAddress string
	DBAddress  string
	DBWorkers  int
}

type MasterConfig struct {
	Server    *ServerConfig
	Analysis  []Analysis
	Reporting []Reporting
}

func getDefaultConfig() *MasterConfig {
	return &MasterConfig{
		Server: &ServerConfig{DBWorkers: 2, RPCAddress: ":50051", DBAddress: "localhost:9080"},
	}
}

func ReadConfigFromFile(configfile string) (*MasterConfig, error) {
	log.Printf("Reading configuration from %s\n", configfile)
	data, err := ConsumeFile(configfile)
	if err != nil {
		return nil, err
	}

	configuration := getDefaultConfig()

	err = yaml.Unmarshal(data, configuration)
	if err != nil {
		return nil, err
	}

	return configuration, nil
}

func ConsumeFile(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return data, nil
}

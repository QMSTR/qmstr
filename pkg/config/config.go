package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/service"
	yaml "gopkg.in/yaml.v2"
)

type Analysis struct {
	Name       string
	PosixName  string
	Analyzer   string
	TrustLevel int64
	PathSub    []*service.PathSubstitution
	Config     map[string]string
}

type Reporting struct {
	Name      string
	PosixName string
	Reporter  string
	Config    map[string]string
}

type ServerConfig struct {
	RPCAddress string
	DBAddress  string
	DBWorkers  int
	OutputDir  string
	CacheDir   string
	ImageName  string `yaml:"image"`
	Debug      bool
	ExtraEnv   map[string]string
	ExtraMount map[string]string
	BuildPath  string
	PathSub    []*service.PathSubstitution
}

type MasterConfig struct {
	Name      string
	MetaData  map[string]string
	Server    *ServerConfig
	Analysis  []Analysis
	Reporting []Reporting
}

type QmstrConfig struct {
	Project *MasterConfig
}

func getDefaultConfig() *QmstrConfig {
	return &QmstrConfig{
		Project: &MasterConfig{
			Server: &ServerConfig{DBWorkers: 2, RPCAddress: ":50051", DBAddress: "localhost:9080",
				ExtraEnv: map[string]string{}, ExtraMount: map[string]string{},
			},
		},
	}
}

func ReadConfigFromFiles(configfiles ...string) (*MasterConfig, error) {
	fileNotExistCount := 0
	config := getDefaultConfig()
	for _, configfile := range configfiles {
		if _, err := os.Stat(configfile); os.IsNotExist(err) {
			log.Printf("File %s not found", configfile)
			fileNotExistCount++
			continue
		}

		log.Printf("Reading configuration from %s\n", configfile)
		data, err := ConsumeFile(configfile)
		if err != nil {
			return nil, err
		}

		if err := readConfig(data, config); err != nil {
			return nil, fmt.Errorf("Failed to read config from %s: %v", configfile, err)
		}

	}

	if fileNotExistCount == len(configfiles) {
		return nil, errors.New("No configuration file found")
	}

	return config.Project, nil
}

func ReadConfigFromBytes(data []byte) (*MasterConfig, error) {
	config := getDefaultConfig()
	err := readConfig(data, config)
	if err != nil {
		return nil, err
	}
	return config.Project, err
}

func readConfig(data []byte, configuration *QmstrConfig) error {
	err := yaml.Unmarshal(data, configuration)
	if err != nil {
		return err
	}
	err = validateConfig(configuration.Project)
	if err != nil {
		return err
	}
	return nil
}

func SerializeConfig(config *MasterConfig) ([]byte, error) {
	data, err := yaml.Marshal(QmstrConfig{Project: config})
	return data, err
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

package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"

	"github.com/QMSTR/qmstr/pkg/service"
	yaml "gopkg.in/yaml.v2"
)

type Analysis struct {
	Name     string `yaml:"name"`
	Selector string
	Analyzer string
	PathSub  []*service.PathSubstitution
	Config   map[string]string
}

type Reporting struct {
	Name     string `yaml:"name"`
	Reporter string
	Config   map[string]string
}

type ServerConfig struct {
	RPCAddress string
	DBAddress  string
	DBWorkers  int
}

type MasterConfig struct {
	Name      string
	MetaData  map[string]string
	Server    *ServerConfig
	Analysis  []Analysis
	Reporting []Reporting
}

type QmstrConfig struct {
	Package *MasterConfig
}

func getDefaultConfig() *QmstrConfig {
	return &QmstrConfig{
		Package: &MasterConfig{
			Server: &ServerConfig{DBWorkers: 2, RPCAddress: ":50051", DBAddress: "localhost:9080"},
		},
	}
}

func ReadConfigFromFile(configfile string) (*MasterConfig, error) {
	log.Printf("Reading configuration from %s\n", configfile)
	data, err := ConsumeFile(configfile)
	if err != nil {
		return nil, err
	}
	return readConfig(data)
}

func readConfig(data []byte) (*MasterConfig, error) {
	configuration := getDefaultConfig()
	err := yaml.Unmarshal(data, configuration)
	if err != nil {
		return nil, err
	}
	err = validateConfig(configuration.Package)
	if err != nil {
		return nil, err
	}
	return configuration.Package, nil
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

func validateConfig(configuration *MasterConfig) error {
	if configuration == nil {
		return fmt.Errorf("empty configuration -- check indentation")
	}

	// Validate analyzers
	for idx, analyzer := range configuration.Analysis {
		err := validateFields(analyzer, "Name", "Analyzer")
		if err != nil {
			return fmt.Errorf("%d. analyzer misconfigured %v", idx+1, err)
		}
	}
	// Validate reporters
	for idx, reporter := range configuration.Reporting {
		err := validateFields(reporter, "Name", "Reporter")
		if err != nil {
			return fmt.Errorf("%d. reporter misconfigured %v", idx+1, err)
		}
	}
	return nil
}

func validateFields(structure interface{}, fields ...string) error {
	v := reflect.ValueOf(structure)
	for _, field := range fields {
		f := v.FieldByName(field)
		if !f.IsValid() || f.Kind() != reflect.String || f.String() == "" {
			return fmt.Errorf("%s invalid", field)
		}
	}
	return nil
}

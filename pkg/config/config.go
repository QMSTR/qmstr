package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"regexp"
	"strings"

	"github.com/QMSTR/qmstr/pkg/service"
	yaml "gopkg.in/yaml.v2"
)

type Analysis struct {
	Name      string `yaml:"name"`
	PosixName string
	Selector  string
	Analyzer  string
	PathSub   []*service.PathSubstitution
	Config    map[string]string
}

type Reporting struct {
	Name      string `yaml:"name"`
	PosixName string
	Reporter  string
	Config    map[string]string
}

type ServerConfig struct {
	RPCAddress string
	DBAddress  string
	DBWorkers  int
	CacheDir   string
	OutputDir  string
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
			// TODO make default output and cache dir platform independent
			Server: &ServerConfig{DBWorkers: 2, RPCAddress: ":50051", DBAddress: "localhost:9080",
				CacheDir: "/var/cache/qmstr", OutputDir: "/var/qmstr"},
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

	serveraddress := strings.Split(configuration.Server.RPCAddress, ":")
	if len(serveraddress) != 2 {
		return errors.New("Invalid RPC address")
	}

	uniqueFields := map[string]map[string]struct{}{}
	uniqueFields["Name"] = map[string]struct{}{}
	uniqueFields["PosixName"] = map[string]struct{}{}

	// Validate analyzers
	for idx, analyzer := range configuration.Analysis {
		if analyzer.PosixName == "" {
			analyzer.PosixName = posixFullyPortableFilename(analyzer.Name)
		}
		err := validateFields(analyzer, uniqueFields, "Name", "Analyzer", "PosixName")
		if err != nil {
			return fmt.Errorf("%d. analyzer misconfigured %v", idx+1, err)
		}
	}
	// Validate reporters
	for idx, reporter := range configuration.Reporting {
		if reporter.PosixName == "" {
			reporter.PosixName = posixFullyPortableFilename(reporter.Name)
		}
		err := validateFields(reporter, uniqueFields, "Name", "Reporter")
		if err != nil {
			return fmt.Errorf("%d. reporter misconfigured %v", idx+1, err)
		}
	}
	return nil
}

func validateFields(structure interface{}, uniqueFields map[string]map[string]struct{}, fields ...string) error {
	v := reflect.ValueOf(structure)
	for _, field := range fields {
		trackSet := map[string]struct{}{}
		if val, ok := uniqueFields[field]; ok {
			trackSet = val
		}
		f := v.FieldByName(field)
		if !f.IsValid() || f.Kind() != reflect.String || f.String() == "" {
			return fmt.Errorf("%s invalid", field)
		}
		if _, ok := trackSet[f.String()]; ok {
			return fmt.Errorf("duplicate value of %s in %s", f.String(), field)
		}
		trackSet[f.String()] = struct{}{}
	}

	return nil
}

func posixFullyPortableFilename(filename string) string {
	nonPosixChars := regexp.MustCompile(`[^A-Za-z0-9\._-]`)
	posixFilename := nonPosixChars.ReplaceAllString(filename, "_")
	return posixFilename
}

// GetRPCPort returns the configured port for qmstr's grpc service
func (mc *MasterConfig) GetRPCPort() (string, error) {
	err := validateConfig(mc)
	if err != nil {
		return "", err
	}
	return strings.Split(mc.Server.RPCAddress, ":")[1], nil
}

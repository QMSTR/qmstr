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
	Name       string `yaml:"name"`
	PosixName  string
	Analyzer   string
	TrustLevel int64
	PathSub    []*service.PathSubstitution
	Config     map[string]string
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
	OutputDir  string
	CacheDir   string
	ImageName  string `yaml:"image"`
	Debug      bool
	ExtraEnv   map[string]string
	ExtraMount map[string]string
	PathSub    []*service.PathSubstitution
}

type MasterConfig struct {
	Name        string
	BuildConfig string
	MetaData    map[string]string
	Server      *ServerConfig
	Analysis    []Analysis
	Reporting   []Reporting
}

type QmstrConfig struct {
	Package *MasterConfig
}

func getDefaultConfig() *QmstrConfig {
	return &QmstrConfig{
		Package: &MasterConfig{
			BuildConfig: "default",
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

		readConfig(data, config)
	}

	if fileNotExistCount == len(configfiles) {
		return nil, errors.New("No configuration file found")
	}

	return config.Package, nil
}

func ReadConfigFromBytes(data []byte) (*MasterConfig, error) {
	config := getDefaultConfig()
	err := readConfig(data, config)
	if err != nil {
		return nil, err
	}
	return config.Package, err
}

func readConfig(data []byte, configuration *QmstrConfig) error {
	err := yaml.Unmarshal(data, configuration)
	if err != nil {
		return err
	}
	err = validateConfig(configuration.Package)
	if err != nil {
		return err
	}
	return nil
}

func SerializeConfig(config *MasterConfig) ([]byte, error) {
	data, err := yaml.Marshal(QmstrConfig{Package: config})
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
		err := validateFields(reporter, uniqueFields, "Name", "Reporter", "PosixName")
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

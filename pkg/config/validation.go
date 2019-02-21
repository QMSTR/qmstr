package config

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/QMSTR/qmstr/pkg/common"
)

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
			analyzer.PosixName = common.GetPosixFullyPortableFilename(analyzer.Name)
		}
		err := validateFields(analyzer, uniqueFields, "Name", "Analyzer", "PosixName")
		if err != nil {
			return fmt.Errorf("%d. analyzer misconfigured %v", idx+1, err)
		}
	}
	// Validate reporters
	for idx, reporter := range configuration.Reporting {
		if reporter.PosixName == "" {
			reporter.PosixName = common.GetPosixFullyPortableFilename(reporter.Name)
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

// GetRPCPort returns the configured port for qmstr's grpc service
func (mc *MasterConfig) GetRPCPort() (string, error) {
	err := validateConfig(mc)
	if err != nil {
		return "", err
	}
	return strings.Split(mc.Server.RPCAddress, ":")[1], nil
}

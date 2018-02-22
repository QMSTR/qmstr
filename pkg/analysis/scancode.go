package analysis

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"

	"github.com/QMSTR/qmstr/pkg/database"
)

type ScancodeAnalyzer struct {
	ScanData interface{}
	db       *database.DataBase
}

func NewScancodeAnalyzer(config map[string]string, db *database.DataBase) (*ScancodeAnalyzer, error) {
	if workdir, ok := config["workdir"]; ok {
		scanData := scancode(workdir, runtime.NumCPU())
		return &ScancodeAnalyzer{ScanData: scanData, db: db}, nil
	}
	return nil, errors.New("scancode analyzer configuration missing \"workdir\"")
}

func scancode(workdir string, jobs int) interface{} {
	cmdlineargs := []string{"--quiet", "--full-root"}
	if jobs > 1 {
		cmdlineargs = append(cmdlineargs, "--processes", fmt.Sprintf("%d", jobs))
	}
	cmd := exec.Command("scancode", append(cmdlineargs, workdir)...)
	log.Printf("Calling %s", cmd.Args)
	scanResult, err := cmd.Output()
	if err != nil {
		log.Printf("Scandir failed %s", err)
	}
	re := regexp.MustCompile("{.+")
	jsonScanResult := re.Find(scanResult)
	var scanData interface{}
	err = json.Unmarshal(jsonScanResult, &scanData)
	if err != nil {
		log.Printf("parsing scan data failed %s", err)
	}
	log.Printf("SCANDATA: %v", scanData)
	return scanData
}

func (scan *ScancodeAnalyzer) Analyze(node *AnalysisNode) error {
	log.Printf("Analyze %s", node.GetPath())
	licenses := scan.detectLicenses(node.GetPath())

	// TODO merge with currently set licenses
	for _, license := range licenses {
		err := node.SetLicense(license)
		if err != nil {
			return fmt.Errorf("failed to set license %v: %v", license, err)
		}
	}
	return nil
}

func (scan *ScancodeAnalyzer) detectLicenses(srcFilePath string) []*database.License {
	licenses := []*database.License{}
	scanDatamap := scan.ScanData.(map[string]interface{})
	for _, file := range scanDatamap["files"].([]interface{}) {
		fileData := file.(map[string]interface{})
		if fileData["path"] == srcFilePath {
			log.Printf("Found %s", srcFilePath)
			for _, licenseData := range fileData["licenses"].([]interface{}) {
				license := licenseData.(map[string]interface{})
				licenses = append(licenses, &database.License{
					Key:            license["key"].(string),
					SpdxIdentifier: license["spdx_license_key"].(string),
				})
			}
		}
	}
	if len(licenses) == 0 {
		licenses = []*database.License{&database.License{Key: "UNKNOWN"}}
	}
	return licenses
}

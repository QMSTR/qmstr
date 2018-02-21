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
	fmt.Printf("Calling %s", cmd.Args)
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
	licenses, err := scan.detectLicenses(node.GetPath())
	if err != nil {
		node.SetLicense("UNKNOWN")
		return err
	}

	// TODO merge with currently set licenses
	for _, license := range licenses {
		node.SetLicense(license)
	}
	return nil
}

func (scan *ScancodeAnalyzer) detectLicenses(srcFilePath string) ([]string, error) {
	licenses := []string{}
	scanDatamap := scan.ScanData.(map[string]interface{})
	for _, file := range scanDatamap["files"].([]interface{}) {
		fileData := file.(map[string]interface{})
		if fileData["path"] == srcFilePath {
			fmt.Printf("Found %s", srcFilePath)
			for _, license := range fileData["licenses"].([]interface{}) {
				licenses = append(licenses, license.(map[string]interface{})["spdx_license_key"].(string))
			}
		}
	}
	return licenses, nil
}

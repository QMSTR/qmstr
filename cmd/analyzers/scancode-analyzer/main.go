package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/service"
)

type ScancodeAnalyzer struct {
	scanData interface{}
}

func main() {
	analyzer := analysis.NewAnalyzer(&ScancodeAnalyzer{})
	analyzer.RunAnalyzerPlugin()
}

func (scanalyzer *ScancodeAnalyzer) Configure(configMap map[string]string) error {
	if workdir, ok := configMap["workdir"]; ok {
		scanalyzer.scanData = scancode(workdir, runtime.NumCPU())
	}
	return nil
}

func (scanalyzer *ScancodeAnalyzer) Analyze(node *service.FileNode) (*service.InfoNodeSlice, error) {
	log.Printf("Analyzing file %s", node.Path)

	retVal := &service.InfoNodeSlice{Inodes: []*service.InfoNode{}}

	licenseInfo, err := scanalyzer.detectLicenses(node.GetPath())
	if err == nil {
		retVal.Inodes = append(retVal.Inodes, licenseInfo...)
	}
	copyrightHolders, err := scanalyzer.detectCopyrightHolders(node.GetPath())
	if err == nil {
		retVal.Inodes = append(retVal.Inodes, copyrightHolders)
	}

	return retVal, nil
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
		log.Printf("scancode failed %s", err)
	}
	re := regexp.MustCompile("{.+")
	jsonScanResult := re.Find(scanResult)
	var scanData interface{}
	err = json.Unmarshal(jsonScanResult, &scanData)
	if err != nil {
		log.Printf("parsing scan data failed %s", err)
	}
	return scanData
}

func (scanalyzer *ScancodeAnalyzer) detectLicenses(srcFilePath string) ([]*service.InfoNode, error) {
	licenseNodes := []*service.InfoNode{}
	scanDatamap := scanalyzer.scanData.(map[string]interface{})
	for _, file := range scanDatamap["files"].([]interface{}) {
		fileData := file.(map[string]interface{})
		if fileData["path"] == srcFilePath {
			log.Printf("Found %s", srcFilePath)
			for _, licenseData := range fileData["licenses"].([]interface{}) {
				license := licenseData.(map[string]interface{})
				tempDataNodes := []*service.InfoNode_DataNode{&service.InfoNode_DataNode{
					Type: "key",
					Data: license["key"].(string),
				},
					&service.InfoNode_DataNode{
						Type: "score",
						Data: strconv.FormatFloat(license["score"].(float64), 'f', 2, 64),
					},
				}

				spdxIdent := license["spdx_license_key"].(string)

				if spdxIdent != "" {
					tempDataNodes = append(tempDataNodes, &service.InfoNode_DataNode{
						Type: "spdxIdentifier",
						Data: spdxIdent,
					})
				}

				licenseNodes = append(licenseNodes, &service.InfoNode{
					Type:      "license",
					DataNodes: tempDataNodes,
				})
			}
			return licenseNodes, nil
		}
	}
	return nil, fmt.Errorf("No license found for %s", srcFilePath)
}

func (scanalyzer *ScancodeAnalyzer) detectCopyrightHolders(srcFilePath string) (*service.InfoNode, error) {
	copyrightHolders := []*service.InfoNode_DataNode{}
	scanDatamap := scanalyzer.scanData.(map[string]interface{})
	for _, file := range scanDatamap["files"].([]interface{}) {
		fileData := file.(map[string]interface{})
		if fileData["path"] == srcFilePath {
			for _, copyright := range fileData["copyrights"].([]interface{}) {
				copyrightData := copyright.(map[string]interface{})
				for _, copyrightHolder := range copyrightData["holders"].([]interface{}) {
					log.Printf("Found copyright holder %s", copyrightHolder)
					copyrightHolders = append(copyrightHolders, &service.InfoNode_DataNode{
						Type: "copyrightHolder",
						Data: copyrightHolder.(string),
					})
				}
			}
			copyHolderInfoNode := &service.InfoNode{
				Type:      "copyright",
				DataNodes: copyrightHolders,
			}
			return copyHolderInfoNode, nil
		}
	}
	return nil, fmt.Errorf("No copyright info found for %s", srcFilePath)
}

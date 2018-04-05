package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"syscall"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/service"
)

const scancodeErrorCount = "Errors count:\\s*(\\d+)"

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
	copyrights, err := scanalyzer.detectCopyrights(node.GetPath())
	if err == nil {
		retVal.Inodes = append(retVal.Inodes, copyrights)
	}
	return retVal, nil
}

func scancodeExitHandler(err error, output []byte) {
	// print to stdout so master server can log it
	fmt.Print(output)
	if output != nil {
		// scancode might have failed, let's see ...
		re := regexp.MustCompile(scancodeErrorCount)
		errors := re.FindSubmatch(output)
		if len(errors) > 1 && string(errors[1]) == "0" {
			return
		}
	}
	if err != nil {
		log.Printf("Scancode failed %s\n", err)
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// preserve non-zero return code
				os.Exit(status.ExitStatus())
			}
		}
		log.Fatalln("scancode return code was not found. Something went seriously wrong.")
	}
}

func scancode(workdir string, jobs int) interface{} {
	tmpfile, err := ioutil.TempFile("", "scancode")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	cmdlineargs := []string{"--full-root", "-l", "-c", "--json", tmpfile.Name()}
	if jobs > 1 {
		cmdlineargs = append(cmdlineargs, "--processes", fmt.Sprintf("%d", jobs))
	}
	cmd := exec.Command("scancode", append(cmdlineargs, workdir)...)
	log.Printf("Calling %s\n", cmd.Args)

	output, err := cmd.CombinedOutput()
	scancodeExitHandler(err, output)

	re := regexp.MustCompile("{.+")
	scanResult, err := ioutil.ReadFile(tmpfile.Name())
	if err != nil {
		log.Fatalf("Could not read from result file %v", err)
	}
	jsonScanResult := re.Find(scanResult)
	var scanData interface{}
	err = json.Unmarshal(jsonScanResult, &scanData)
	if err != nil {
		log.Fatalf("parsing scan data failed %s\n", err)
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

func (scanalyzer *ScancodeAnalyzer) detectCopyrights(srcFilePath string) (*service.InfoNode, error) {
	copyrights := []*service.InfoNode_DataNode{}
	scanDatamap := scanalyzer.scanData.(map[string]interface{})
	for _, file := range scanDatamap["files"].([]interface{}) {
		fileData := file.(map[string]interface{})
		if fileData["path"] == srcFilePath {
			for _, copyright := range fileData["copyrights"].([]interface{}) {
				copyrightData := copyright.(map[string]interface{})
				for _, copyrightHolder := range copyrightData["holders"].([]interface{}) {
					log.Printf("Found copyright holder %s", copyrightHolder)
					copyrights = append(copyrights, &service.InfoNode_DataNode{
						Type: "copyrightHolder",
						Data: copyrightHolder.(string),
					})
				}
				for _, author := range copyrightData["authors"].([]interface{}) {
					log.Printf("Found author %s", author)
					copyrights = append(copyrights, &service.InfoNode_DataNode{
						Type: "author",
						Data: author.(string),
					})
				}
			}
			copyrightInfoNode := &service.InfoNode{
				Type:      "copyright",
				DataNodes: copyrights,
			}
			return copyrightInfoNode, nil
		}
	}
	return nil, fmt.Errorf("No copyright info found for %s", srcFilePath)
}

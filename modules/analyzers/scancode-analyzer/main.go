package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"syscall"

	"github.com/QMSTR/qmstr/lib/go-qmstr/analysis"
	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/module"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

const (
	scancodeErrorCount = "Errors count:\\s*(\\d+)"
)

type ScancodeAnalyzer struct {
	Paths map[string][]*InfoStruct
}

type InfoStruct struct {
	Type            string        `json:"type"`
	ConfidenceScore float64       `json:"confidenceScore,omitempty"`
	DataNodes       []*DataStruct `json:"dataNodes"`
}

type DataStruct struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

func main() {
	analyzer := analysis.NewAnalyzer(&ScancodeAnalyzer{})
	if err := analyzer.RunAnalyzerModule(); err != nil {
		log.Printf("%v failed: %v\n", analyzer.GetModuleName(), err)
		os.Exit(master.ReturnAnalyzerFailed)
	}
}

func (scanalyzer *ScancodeAnalyzer) Configure(configMap map[string]string) error {
	cached := false
	if usecache, ok := configMap["cached"]; ok && usecache == "true" {
		cached = true
	}

	resultfile := ""
	if rf, ok := configMap["resultfile"]; ok {
		resultfile = rf
	}

	var err error
	if cached && resultfile != "" {
		if err = scanalyzer.readScancodeFile(resultfile); err != nil {
			log.Fatal(err)
		}
		return nil
	}

	if workdir, ok := configMap["workdir"]; ok {
		err = scanalyzer.scancode(workdir, runtime.NumCPU(), resultfile)
		if err != nil {
			log.Fatal(err)
		}
		return nil
	}

	return fmt.Errorf("Misconfigured scancode analyzer")
}

func (scanalyzer *ScancodeAnalyzer) Analyze(masterClient *module.MasterClient, token int64) error {
	stream, err := masterClient.AnaSvcClient.GetSourceFileNodes(context.Background(), &service.DummyRequest{})
	if err != nil {
		log.Printf("failed getting source file nodes %v", err)
		return err
	}

	infoNodeMsgs := []*service.InfoNodesMessage{}

	for {
		fileNode, err := stream.Recv()
		if err == io.EOF {
			break
		}

		log.Printf("Analyzing file %s", fileNode.Path)

		infoNodes, _ := scanalyzer.detectLicensesCopyrights(fileNode.GetPath())

		if len(infoNodes) > 0 {
			infoNodeMsgs = append(infoNodeMsgs, &service.InfoNodesMessage{Token: token, Infonodes: infoNodes, Uid: fileNode.FileData.Uid})
		}
	}

	sendStream, err := masterClient.AnaSvcClient.SendInfoNodes(context.Background())
	if err != nil {
		return err
	}
	for _, inodeMsg := range infoNodeMsgs {
		sendStream.Send(inodeMsg)
	}

	reply, err := sendStream.CloseAndRecv()
	if err != nil {
		return err
	}
	if reply.Success {
		log.Println("Scancode Analyzer sent InfoNodes")
	}
	return nil
}

func (scanalyzer *ScancodeAnalyzer) PostAnalyze() error {
	return nil
}

func scancodeExitHandler(err error, output []byte) {
	// print to stdout so master server can log it
	fmt.Print(string(output))
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

func (scanalyzer *ScancodeAnalyzer) readScancodeFile(resultfile string) error {
	scanResult, err := ioutil.ReadFile(resultfile)
	if err != nil {
		return fmt.Errorf("Error while reading the scancode input file: %v", err)
	}
	err = json.Unmarshal(scanResult, &scanalyzer.Paths)
	if err != nil {
		return fmt.Errorf("Failed while unmarshalling json file: %v", err)
	}
	return nil
}

func (scanalyzer *ScancodeAnalyzer) scancode(workdir string, jobs int, resultfilepath string) error {
	if resultfilepath == "" {
		tmpfile, err := ioutil.TempFile("", "qmstr-scancode")
		resultfilepath = tmpfile.Name()
		if err != nil {
			return err
		}
		defer os.Remove(tmpfile.Name())
	}

	cmdlineargs := []string{"--full-root", "-l", "-c", "--custom-output", resultfilepath, "--custom-template", "qmstr.j2"}
	if jobs > 1 {
		cmdlineargs = append(cmdlineargs, "--processes", fmt.Sprintf("%d", jobs))
	}
	cmd := exec.Command("scancode", append(cmdlineargs, workdir)...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Failed running scancode command line: %v", err)
	}
	scancodeExitHandler(err, output)

	return scanalyzer.readScancodeFile(resultfilepath)
}

func (scanalyzer *ScancodeAnalyzer) detectLicensesCopyrights(srcFilePath string) ([]*service.InfoNode, error) {
	if fileData, ok := scanalyzer.Paths[srcFilePath]; ok {
		infoNodes := []*service.InfoNode{}
		for _, infoData := range fileData {
			log.Printf("Found license/copyright information for path: %s\n", srcFilePath)

			tempDataNodes := []*service.InfoNode_DataNode{}
			for _, iData := range infoData.DataNodes {
				tempDataNodes = append(tempDataNodes, &service.InfoNode_DataNode{
					Type: iData.Type,
					Data: iData.Data,
				})
			}
			infoNodes = append(infoNodes, &service.InfoNode{
				Type:            infoData.Type,
				ConfidenceScore: infoData.ConfidenceScore / 100,
				DataNodes:       tempDataNodes,
			})
		}
		return infoNodes, nil
	}
	return nil, fmt.Errorf("No license/copyright information found for %s", srcFilePath)
}

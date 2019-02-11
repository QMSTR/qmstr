package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/QMSTR/go-qmstr/analysis"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/go-qmstr/service"
)

var spdxPattern = regexp.MustCompile(`SPDX-License-Identifier: (.+)\s*`)

type SpdxAnalyzer struct {
	basedir string
}

func main() {
	analyzer := analysis.NewAnalyzer(&SpdxAnalyzer{})
	if err := analyzer.RunAnalyzerModule(); err != nil {
		log.Printf("%v failed: %v\n", analyzer.GetModuleName(), err)
		os.Exit(master.ReturnAnalyzerFailed)
	}
}

func (spdxalizer *SpdxAnalyzer) Configure(configMap map[string]string) error {
	if workdir, ok := configMap["workdir"]; ok {
		spdxalizer.basedir = workdir
	} else {
		return fmt.Errorf("no working directory configured")
	}

	return nil
}

func (spdxalizer *SpdxAnalyzer) Analyze(controlService service.ControlServiceClient, analysisService service.AnalysisServiceClient, token int64) error {
	queryNode := &service.FileNode{FileType: service.FileNode_SOURCE}

	stream, err := controlService.GetFileNode(context.Background(), queryNode)
	if err != nil {
		log.Printf("Could not get file node %v", err)
		return err
	}

	for {
		fileNode, err := stream.Recv()
		if err == io.EOF {
			break
		}

		diagnosticNodeMsg := service.DiagnosticNodeMessage{}

		log.Printf("Analyzing file %s", fileNode.Path)
		spdxIdent, lineNo, columnNo, err := detectSPDXLicense(fileNode.Path)
		if err != nil {
			log.Printf("%v", err)
			// Adding warning node
			diagnosticWarningNode := service.DiagnosticNode{Severity: service.DiagnosticNode_WARNING, Message: fmt.Sprintf("%v", err)}
			diagnosticNodeMsg = service.DiagnosticNodeMessage{Token: token, Diagnosticnode: &diagnosticWarningNode, Uid: fileNode.Uid}
		} else if _, ok := spdxLicenses[spdxIdent]; !ok {
			log.Printf("Found invalid spdx license identifier %v.", spdxIdent)
			// Adding error node
			file, err := filepath.Rel(spdxalizer.basedir, fileNode.Path)
			if err != nil {
				return err
			}
			diagnosticErrorNode := service.DiagnosticNode{Severity: service.DiagnosticNode_ERROR, Message: fmt.Sprintf("%v:%v:%v Invalid SPDX license expression %v", file, lineNo, columnNo, spdxIdent)}
			diagnosticNodeMsg = service.DiagnosticNodeMessage{Token: token, Diagnosticnode: &diagnosticErrorNode, Uid: fileNode.Uid}
		} else {
			diagnosticInfoNode := service.DiagnosticNode{Severity: service.DiagnosticNode_INFO, Message: fmt.Sprintf("SPDX license expression detected: %s", spdxIdent)}
			diagnosticNodeMsg = service.DiagnosticNodeMessage{Token: token, Diagnosticnode: &diagnosticInfoNode, Uid: fileNode.Uid}
		}
		sendStream, err := analysisService.SendDiagnosticNode(context.Background())
		if err != nil {
			return err
		}
		err = sendStream.Send(&diagnosticNodeMsg)
		if err != nil {
			return err
		}
		reply, err := sendStream.CloseAndRecv()
		if reply.Success {
			log.Println("Simple SPDX Analyzer sent DiagnosticNodes")
		}
	}
	return nil
}

func (spdxalizer *SpdxAnalyzer) PostAnalyze() error {
	return nil
}

func detectSPDXLicense(srcFilePath string) (string, int, int, error) {
	f, err := os.Open(srcFilePath)
	if err != nil {
		return "", 0, 0, err
	}

	scanner := bufio.NewScanner(f)
	lineNo := 0

	for scanner.Scan() {
		lineNo++
		if lineNo > 100 {
			break
		}
		line := scanner.Text()

		matches := spdxPattern.FindStringSubmatch(line)
		if matches != nil {
			column := spdxPattern.FindStringSubmatchIndex(line)
			return matches[1], lineNo, column[2], nil
		}
	}
	return "", 0, 0, fmt.Errorf("No SPDX license identifier found")
}

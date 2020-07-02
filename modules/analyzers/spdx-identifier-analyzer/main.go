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

	"github.com/QMSTR/qmstr/lib/go-qmstr/analysis"
	"github.com/QMSTR/qmstr/lib/go-qmstr/cli"
	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/module"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

var spdxPattern = regexp.MustCompile(`SPDX-License-Identifier: (.+)\s*`)

type SpdxAnalyzer struct {
	basedir string
}

func main() {
	analyzer := analysis.NewAnalyzer(&SpdxAnalyzer{})
	go func() {
		<-cli.PingAnalyzer // wait for the analysis phase to start
		log.Printf("Spdx identifier analyzer starts the analysis")
		if err := analyzer.RunAnalyzerModule(); err != nil {
			msg := fmt.Sprintf("Analyzer %v failed: %v\n", analyzer.GetModuleName(), err)
			log.Printf(msg)
			analyzer.CtrlSvcClient.ShutdownModule(context.Background(), &service.ShutdownModuleRequest{
				Message: msg, DB: true})
			os.Exit(master.ReturnAnalyzerFailed)
		}
		analysis.ReduceAnalyzersCounter()
	}()
}

func (spdxalizer *SpdxAnalyzer) Configure(configMap map[string]string) error {
	if workdir, ok := configMap["workdir"]; ok {
		spdxalizer.basedir = workdir
	} else {
		return fmt.Errorf("no working directory configured")
	}

	return nil
}

func (spdxalizer *SpdxAnalyzer) Analyze(masterClient *module.MasterClient, token int64) error {
	stream, err := masterClient.AnaSvcClient.GetSourceFileNodes(context.Background(), &service.DummyRequest{})
	if err != nil {
		log.Printf("failed getting source file nodes %v", err)
		return err
	}

	for {
		fileNode, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		log.Printf("Analyzing file %s", fileNode.Path)
		diagnosticNodeMsg := &service.DiagnosticNodeMessage{}
		spdxIdent, lineNo, columnNo, err := detectSPDXLicense(fileNode.Path)
		if err != nil {
			log.Printf("%v", err)
			// Adding warning node
			diagnosticWarningNode := &service.DiagnosticNode{Severity: service.DiagnosticNode_WARNING, Message: fmt.Sprintf("%v", err)}
			diagnosticNodeMsg = &service.DiagnosticNodeMessage{Token: token, Diagnosticnode: diagnosticWarningNode, Uid: fileNode.FileData.Uid}
		} else if _, ok := spdxLicenses[spdxIdent]; !ok {
			log.Printf("Found invalid spdx license identifier %v.", spdxIdent)
			// Adding error node
			file, err := filepath.Rel(spdxalizer.basedir, fileNode.Path)
			if err != nil {
				return err
			}
			diagnosticErrorNode := &service.DiagnosticNode{Severity: service.DiagnosticNode_ERROR, Message: fmt.Sprintf("%v:%v:%v Invalid SPDX license expression %v", file, lineNo, columnNo, spdxIdent)}
			diagnosticNodeMsg = &service.DiagnosticNodeMessage{Token: token, Diagnosticnode: diagnosticErrorNode, Uid: fileNode.FileData.Uid}
		} else {
			// Create both diagnostic and info node
			diagnosticInfoNode := &service.DiagnosticNode{Severity: service.DiagnosticNode_INFO, Message: fmt.Sprintf("SPDX license expression detected: %s", spdxIdent)}
			diagnosticNodeMsg = &service.DiagnosticNodeMessage{Token: token, Diagnosticnode: diagnosticInfoNode, Uid: fileNode.FileData.Uid}

			dataNodes := []*service.InfoNode_DataNode{
				&service.InfoNode_DataNode{Type: "name", Data: spdxIdent},
				&service.InfoNode_DataNode{Type: "spdxIdentifier", Data: spdxIdent},
			}
			infoNodes := []*service.InfoNode{&service.InfoNode{Type: "license", DataNodes: dataNodes}}
			infoNodeMsg := &service.InfoNodesMessage{Token: token, Infonodes: infoNodes, Uid: fileNode.FileData.Uid}

			sendStream, err := masterClient.AnaSvcClient.SendInfoNodes(context.Background())
			if err != nil {
				return err
			}
			err = sendStream.Send(infoNodeMsg)
			if err != nil {
				return err
			}
			reply, err := sendStream.CloseAndRecv()
			if err != nil {
				return err
			}
			if reply.Success {
				log.Println("Simple SPDX Analyzer sent InfoNodes")
			}
		}
		sendStream, err := masterClient.AnaSvcClient.SendDiagnosticNode(context.Background())
		if err != nil {
			return err
		}
		err = sendStream.Send(diagnosticNodeMsg)
		if err != nil {
			return err
		}
		reply, err := sendStream.CloseAndRecv()
		if err != nil {
			return err
		}
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

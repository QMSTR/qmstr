package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

var spdxPattern = regexp.MustCompile(`SPDX-License-Identifier: (.+)\s*`)

type SpdxAnalyzer struct{}

func main() {
	analyzer := analysis.NewAnalyzer(&SpdxAnalyzer{})
	if err := analyzer.RunAnalyzerModule(); err != nil {
		log.Printf("%v failed: %v\n", analyzer.GetModuleName(), err)
		os.Exit(master.ReturnAnalyzerFailed)
	}
}

func (spdxalizer *SpdxAnalyzer) Configure(configMap map[string]string) error {
	return nil
}

func (spdxalizer *SpdxAnalyzer) Analyze(controlService service.ControlServiceClient, analysisService service.AnalysisServiceClient, token int64, session string) error {
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

		infoNodeMsg := service.InfoNodeMessage{}
		log.Printf("Analyzing file %s", fileNode.Path)
		spdxIdent, err := detectSPDXLicense(fileNode.Path)
		if err != nil {
			log.Printf("%v", err)
			// Adding warning node
			warningNode := service.CreateWarningNode(fmt.Sprintf("%v", err))
			infoNodeMsg = service.InfoNodeMessage{Token: token, Infonode: warningNode, Uid: fileNode.Uid}
		} else if _, ok := analysis.SpdxLicenses[spdxIdent]; !ok {
			log.Printf("Found invalid spdx license identifier %v.", spdxIdent)
			// Adding error node
			errorNode := service.CreateErrorNode(fmt.Sprintf("Invalid SPDX license expression %v", spdxIdent))
			infoNodeMsg = service.InfoNodeMessage{Token: token, Infonode: errorNode, Uid: fileNode.Uid}
		} else {
			licenseNode := service.InfoNode{
				Type: "license",
				DataNodes: []*service.InfoNode_DataNode{
					&service.InfoNode_DataNode{
						Type: "spdxIdentifier",
						Data: spdxIdent,
					},
				},
			}
			infoNodeMsg = service.InfoNodeMessage{Token: token, Infonode: &licenseNode, Uid: fileNode.Uid}
		}
		sendStream, err := analysisService.SendInfoNodes(context.Background())
		if err != nil {
			return err
		}
		err = sendStream.Send(&infoNodeMsg)
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
	return nil
}

func (spdxalizer *SpdxAnalyzer) PostAnalyze() error {
	return nil
}

func detectSPDXLicense(srcFilePath string) (string, error) {
	f, err := os.Open(srcFilePath)
	if err != nil {
		return "", err
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
			return matches[1], nil
		}
	}
	return "", fmt.Errorf("No SPDX license identifier found")
}

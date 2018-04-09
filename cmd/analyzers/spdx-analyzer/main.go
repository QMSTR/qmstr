package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/service"
)

var spdxPattern = regexp.MustCompile(`SPDX-License-Identifier: (.+)\s*`)

type SpdxAnalyzer struct{}

func main() {
	analyzer := analysis.NewAnalyzer("SPDX Analyzer", &SpdxAnalyzer{})
	analyzer.RunAnalyzerModule()
}

func (spdxalizer *SpdxAnalyzer) Configure(configMap map[string]string) error {
	return nil
}

func (spdxalizer *SpdxAnalyzer) Analyze(node *service.FileNode) (*service.InfoNodeSlice, error) {
	log.Printf("Analyzing file %s", node.Path)
	spdxIdent, err := detectSPDXLicense(node.Path)
	if err != nil {
		log.Printf("No SPDX license identifier found")
		return &service.InfoNodeSlice{Inodes: []*service.InfoNode{}}, nil
	}

	licenseNode := service.InfoNode{
		Type: "license",
		DataNodes: []*service.InfoNode_DataNode{
			&service.InfoNode_DataNode{
				Type: "spdxIdentifier",
				Data: spdxIdent,
			},
		},
	}
	return &service.InfoNodeSlice{Inodes: []*service.InfoNode{&licenseNode}}, nil
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
	return "", fmt.Errorf("No spdx identifier found in %s", srcFilePath)
}

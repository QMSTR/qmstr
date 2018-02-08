package analysis

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/QMSTR/qmstr/pkg/database"
)

var spdxPattern = regexp.MustCompile(`SPDX-License-Identifier: (.+)\s*`)

type SpdxAnalyzer struct {
	Config map[string]string
	db     *database.DataBase
}

func NewSpdxAnalyzer(config map[string]string, db *database.DataBase) *SpdxAnalyzer {
	return &SpdxAnalyzer{Config: config, db: db}
}

func (spdx *SpdxAnalyzer) Analyze(node *AnalysisNode) error {
	license, err := detectSPDXLicense(node.GetPath())
	if err != nil {
		return err
	}

	// TODO merge with currently set licenses
	node.SetLicense(license)
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
	return "", fmt.Errorf("No spdx identifier found in %s", srcFilePath)
}

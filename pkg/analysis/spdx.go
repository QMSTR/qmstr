package analysis

import (
	"bufio"
	"fmt"
	"log"
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
		log.Printf("Error detecting license %v", err)
		err2 := node.SetLicense(database.UnknownLicense)
		if err2 != nil {
			return fmt.Errorf("failed to set license %v: %v", license, err2)
		}
		return nil
	}

	// TODO merge with currently set licenses
	node.SetLicense(&database.License{Key: license, SpdxIdentifier: license})
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

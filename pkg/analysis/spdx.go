package analysis

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

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

func (spdx *SpdxAnalyzer) Analyze(node *database.Node) error {
	log.Printf("Spdx Analyzer analyzing %s", node.Name)

	actualPath := node.Path
	if pathsub, ok := spdx.Config["pathSubstitution"]; ok {
		subs := strings.Split(pathsub, ":")
		actualPath = strings.Replace(node.Path, subs[0], subs[1], 1)
	}

	license, err := detectSPDXLicense(actualPath)
	if err != nil {
		return err
	}

	uid, err := spdx.db.GetLicenseUid(license)
	if err != nil {
		return err
	}
	// TODO merge with currently set licenses
	node.License = database.License{Uid: uid}
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

//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

//"github.com/endocode/qmstr/pkg/reporter/htmlreporter"

func main() {
	if _, err := DetectHugoAndVerifyVersion(); err != nil {
		log.Printf("error generating HTML reports: %v", err)
	}
	// htmlreporter.ConnectToMaster()
	// htmlreporter.Temp()
	// defer htmlreporter.DisconnectFromMaster()
}

// DetectHugoAndVerifyVersion runs Hugo to get the version string.
func DetectHugoAndVerifyVersion() (string, error) {
	cmd := exec.Command("hugo", "version", "--quiet")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		log.Printf("INFO - PATH is %v", os.Getenv("PATH"))
		return "", fmt.Errorf("Hugo not detected, aborting")
	}
	return ParseAndCheckVersion(output)
}

// ParseAndCheckVersion returns the version for both released and self-compiled versions
func ParseAndCheckVersion(output []byte) (string, error) {
	// is this a version built from a repository?
	re := regexp.MustCompile("Site Generator v(.+)-(.+) .+/.+ BuildDate")
	match := re.FindSubmatch(output)
	if match != nil {
		version := strings.TrimSpace(string(match[1][:]))
		tag := strings.TrimSpace(string(match[2][:]))
		log.Printf("Detected Hugo %v-%v", version, tag)
		return version, nil
	}

	re = regexp.MustCompile("Site Generator v(.+) .+/.+ BuildDate")
	match = re.FindSubmatch(output)
	if match != nil {
		version := strings.TrimSpace(string(match[1][:]))
		log.Printf("Detected Hugo release %v", version)
		return version, nil
	}
	return "", fmt.Errorf(" Unable to parse Hugo version in \"%v\"", string(output[:]))
}

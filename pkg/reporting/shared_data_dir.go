package reporting

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"
)

// DetectSharedDataDirectory detects the shared data directory for all of QMSTR.
// It looks for /usr/share/qmstr, /usr/local/share/qmstr and /opt/share/qmstr, in that order.
func DetectSharedDataDirectory() (string, error) {
	var sharedDataLocations = []string{"/usr/share/qmstr", "/usr/local/share/qmstr", "/opt/share/qmstr"}
	for _, location := range sharedDataLocations {
		fileInfo, err := os.Stat(location)
		if err != nil {
			continue
		}
		if !fileInfo.IsDir() {
			return "", fmt.Errorf("shared data directory exists at %v, but is not a directory, strange", location)
		}
		log.Printf("shared data directory identified at %v", location) // Debug...
		return location, nil
	}
	return "", fmt.Errorf("no suitable QMSTR shared data location found (candidates are %s)", strings.Join(sharedDataLocations, ", "))
}

// DetectModuleSharedDataDirectory detects the directory where QMSTR's shared data is stored.
func DetectModuleSharedDataDirectory(moduleName string) (string, error) {
	sharedDataLocation, err := DetectSharedDataDirectory()
	if err != nil {
		return "", err
	}
	moduleDataLocation := path.Join(sharedDataLocation, moduleName)
	fileInfo, err := os.Stat(moduleDataLocation)
	if err != nil {
		return "", fmt.Errorf("module shared data directory %v not accessible: %v", moduleDataLocation, err)
	}
	if !fileInfo.IsDir() {
		return "", fmt.Errorf("module shared data directory %v not found in shared data directory at %v", moduleDataLocation, sharedDataLocation)
	}
	log.Printf("module shared data directory identified at %v", moduleDataLocation)
	return moduleDataLocation, nil
}

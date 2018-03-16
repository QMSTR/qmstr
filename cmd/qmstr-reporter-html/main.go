//go:generate protoc -I ../../pkg/buildservice --go_out=plugins=grpc:../../pkg/buildservice ../../pkg/buildservice/buildservice.proto
package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"regexp"
	"strings"
	"text/template"

	version "github.com/hashicorp/go-version"
)

//"github.com/endocode/qmstr/pkg/reporter/htmlreporter"

const minimumHugoVersion = "0.32"
const themeDirectory = "/usr/share/qmstr/reporter-html/theme"

func main() {
	detectedVersion, err := DetectHugoAndVerifyVersion()
	if err != nil {
		log.Fatalf("error generating HTML reports: %v", err)
	}
	log.Printf("detected beautiful Hugo version %v", detectedVersion)

	// htmlreporter.ConnectToMaster()
	wd, err := CreateHugoWorkingDirectory()
	if err != nil {
		log.Fatalf("error preparing temporary Hugo working directory: %v", err)
	}
	log.Printf("created temporary Hugo working directory in %v", wd)
	defer func() {
		//TODO is there a --keep option or environment variable to keep temporary files?
		log.Printf("deleting temporary Hugo working directory in %v", wd)
		if err := os.RemoveAll(wd); err != nil {
			log.Printf("warning - error deleting temporary Hugo working directory in %v: %v", wd, err)
		}
	}()
	// generate the content data (markdown and JSON) from the master data model
	// TODO

	staticHTMLContentDir, err := CreateStaticHTML(wd)
	if err != nil {
		log.Fatalf("error generating reports: %v", err)
	}
	log.Printf("QMSTR reports generated in %v", staticHTMLContentDir)

	// defer htmlreporter.DisconnectFromMaster()
}

// DetectHugoAndVerifyVersion runs Hugo to get the version string.
func DetectHugoAndVerifyVersion() (string, error) {
	cmd := exec.Command("hugo", "version", "--quiet")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Command finished with error: %v", err)
		log.Printf("INFO - PATH is %v", os.Getenv("PATH"))
		return "", fmt.Errorf("Hugo not detected")
	}
	version, err := ParseVersion(output)
	if err != nil {
		return version, err
	}
	return version, CheckMinimumRequiredVersion(version)
}

// ParseVersion returns the version for both released and self-compiled versions
func ParseVersion(output []byte) (string, error) {
	// is this a version built from a repository?
	re := regexp.MustCompile("Site Generator v(.+)-(.+) .+/.+ BuildDate")
	match := re.FindSubmatch(output)
	if match != nil {
		version := strings.TrimSpace(string(match[1][:]))
		// tag := strings.TrimSpace(string(match[2][:]))
		return version, nil
	}

	re = regexp.MustCompile("Site Generator v(.+) .+/.+ BuildDate")
	match = re.FindSubmatch(output)
	if match != nil {
		version := strings.TrimSpace(string(match[1][:]))
		return version, nil
	}
	return "", fmt.Errorf(" Unable to parse Hugo version in \"%v\"", string(output[:]))
}

// CheckMinimumRequiredVersion compares the detected Hugo version with the minimum requirement
func CheckMinimumRequiredVersion(versionstring string) error {
	detectedVersion, err := version.NewVersion(versionstring)
	if err != nil {
		return fmt.Errorf("unable to parse \"%v\" as a Hugo version", versionstring)
	}
	minimumVersion, err := version.NewVersion(minimumHugoVersion)
	if err != nil {
		return fmt.Errorf("unable to parse minimum required version \"%v\" as a version (this should not happen)", minimumVersion)
	}
	if detectedVersion.LessThan(minimumVersion) {
		return fmt.Errorf("the Quartermaster HTML reporter requires at least Hugo version %v (version %v found)", minimumVersion, detectedVersion)
	}
	return nil
}

//CreateHugoWorkingDirectory creates a temporary directory with the directory structure required to run Hugo
func CreateHugoWorkingDirectory() (string, error) {
	tmpWorkDir, err := ioutil.TempDir("", "qmstr-")
	if err != nil {
		return tmpWorkDir, fmt.Errorf("error creating temporary Hugo working directory: %v", err)
	}
	// populate working directory with a site template and the theme
	// TODO: add "incremental mode": use an existing, previously generated site and extend it
	cmd := exec.Command("hugo", "new", "site", tmpWorkDir)
	cmd.Dir = tmpWorkDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("NOTE - output:\n%v", string(output[:]))
		return "", fmt.Errorf("error generating site template: %v", err)
	}
	log.Printf("generated Hugo site template in %v", tmpWorkDir)
	// Link the theme directory (see themeDirectory) into the working directory:
	cmd = exec.Command("ln", "-s", themeDirectory, "themes/qmstr-theme")
	cmd.Dir = tmpWorkDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("NOTE - output:\n%v", string(output[:]))
		return "", fmt.Errorf("error linking theme into site template: %v", err)
	}
	// Copy the exampleSite page skeleton:
	// ... The syntax of the copy command is "particular": It copies the *content* of the exampleSite directory.
	// ... Unfortunately, path.Join() strips a trailing dot.
	cmd = exec.Command("cp", "-Rfp", path.Join(themeDirectory, "exampleSite")+"/.", "./.")
	cmd.Dir = tmpWorkDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("NOTE - output:\n%v", string(output[:]))
		return "", fmt.Errorf("error populating site skeleton: %v", err)
	}
	log.Printf("populated site with default content in %v", tmpWorkDir)
	// Generate the configuration file:
	// TODO: Export this type, populate the structure from the reporter configuration:
	type Configuration struct {
		Title, TitleEn, BaseURL string
	}
	configuration := Configuration{
		"Quartermaster Compliance Report",
		"Quartermaster Compliance Report",
		"localhost:1313/",
	}
	configTomlInPath := path.Join(themeDirectory, "config.toml.in")
	configTomlIn, err := ioutil.ReadFile(configTomlInPath)
	if err != nil {
		return "", fmt.Errorf("unable to read configuration template file \"%s\"", configTomlInPath)
	}
	t := template.Must(template.New("site configuration file").Parse(string(configTomlIn)))

	configFilePath := path.Join(tmpWorkDir, "config.toml")
	configFile, err := os.Create(configFilePath)
	if err != nil {
		return "", fmt.Errorf("unable to create configuration file \"%s\v\"", configFilePath)
	}
	defer configFile.Close()
	writer := bufio.NewWriter(configFile)

	if err := t.Execute(writer, configuration); err != nil {
		return "", fmt.Errorf("error applying variables to site configuration: %v", err)
	}
	writer.Flush()
	log.Printf("generated configuration file %v", configFilePath)
	return tmpWorkDir, nil
}

// CreateStaticHTML executes Hugo to generate the static HTML page with the QMSTR reports.
// It returns the path to the generated content and/or an error.
func CreateStaticHTML(contentDir string) (string, error) {
	outputDir := path.Join(contentDir, "qmstr-reports")
	cmd := exec.Command("hugo", "-v", "-d", outputDir)
	cmd.Dir = contentDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("NOTE - output:\n%v", string(output[:]))
		return "", fmt.Errorf("error generating static HTML reports: %v", err)
	}
	log.Printf("generated static HTML reports in %v", outputDir)
	return outputDir, nil

}

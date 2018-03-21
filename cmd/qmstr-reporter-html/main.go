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

	"github.com/QMSTR/qmstr/pkg/reporting"
	"github.com/QMSTR/qmstr/pkg/service"

	version "github.com/hashicorp/go-version"
)

const minimumHugoVersion = "0.32"
const moduleName = "html-reporter"
const themeDirectoryName = "theme"

type HTMLReporter struct {
	workingDir string
	keep       bool
}

func main() {
	reporter := reporting.NewReporter(&HTMLReporter{keep: false})
	reporter.RunReporterPlugin()
}

func (r *HTMLReporter) Configure(config map[string]string) error {
	if val, ok := config["keep"]; ok && val == "true" {
		r.keep = true
	}

	detectedVersion, err := DetectHugoAndVerifyVersion()
	if err != nil {
		return fmt.Errorf("error generating HTML reports: %v", err)
	}
	log.Printf("detected beautiful Hugo version %v", detectedVersion)

	// htmlreporter.ConnectToMaster()
	sharedDataDir, err := DetectSharedDataDirectory()
	if err != nil {
		log.Fatalf("cannot identify QMSTR shared data directory: %v", err)
	}
	r.workingDir, err = CreateHugoWorkingDirectory(sharedDataDir)
	if err != nil {
		return fmt.Errorf("error preparing temporary Hugo working directory: %v", err)
	}
	log.Printf("created temporary Hugo working directory in %v", r.workingDir)

	return nil
}

func (r *HTMLReporter) Report(filenode *service.FileNode) error {
	return nil
}

func (r *HTMLReporter) cleanup() {
	log.Printf("deleting temporary Hugo working directory in %v", r.workingDir)
	if err := os.RemoveAll(r.workingDir); err != nil {
		log.Printf("warning - error deleting temporary Hugo working directory in %v: %v", r.workingDir, err)
	}
}

func (r *HTMLReporter) PostReport() error {
	if !r.keep {
		defer r.cleanup()
	}

	staticHTMLContentDir, err := CreateStaticHTML(r.workingDir)
	if err != nil {
		return fmt.Errorf("error generating reports: %v", err)
	}
	log.Printf("QMSTR reports generated in %v", staticHTMLContentDir)

	// Create the reports package (for publishing, etc).
	// ... TODO: configure default target directory for every reporter (uses CWD at the moment)
	packagePath, _ := os.Getwd()
	if err := CreateReportsPackage(r.workingDir, staticHTMLContentDir, packagePath); err != nil {
		return fmt.Errorf("error packaging report to %v: %v", packagePath, err)
	}
	log.Printf("QMSTR reports package generated in %v", packagePath)
	return nil
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

// DetectSharedDataDirectory detects the directory where QMSTR's shared data is stored.
// It looks for /usr/share/qmstr, /usr/local/share/qmstr and /opt/share/qmstr, in that order.
func DetectSharedDataDirectory() (string, error) {
	var sharedDataLocations = []string{"/usr/share/qmstr", "/usr/local/share/qmstr", "/opt/share/qmstr"}
	for _, location := range sharedDataLocations {
		fileInfo, err := os.Stat(location)
		if err != nil {
			continue
		}
		if fileInfo.IsDir() {
			log.Printf("shared data directory identified at %v", location)
			return location, nil
		}
		return "", fmt.Errorf("shared data directory exists at %v, but is not a directory, strange", location)
	}
	return "", fmt.Errorf("no suitable QMSTR shared data location found (candidates are %s)", strings.Join(sharedDataLocations, ", "))
}

// ParseVersion returns the version for both released and self-compiled versions
func ParseVersion(output []byte) (string, error) {
	// is this a version built from a repository?
	re := regexp.MustCompile("Site Generator v(.+)-(.+) .+/.+ BuildDate")
	match := re.FindSubmatch(output)
	if match != nil {
		version := strings.TrimSpace(string(match[1][:]))
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
func CreateHugoWorkingDirectory(sharedDataDir string) (string, error) {
	tmpWorkDir, err := ioutil.TempDir("", "qmstr-")
	if err != nil {
		return tmpWorkDir, fmt.Errorf("error creating temporary Hugo working directory: %v", err)
	}
	themeDirectory := path.Join(sharedDataDir, themeDirectoryName)
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

// CreateReportsPackage creates a tarball of the static HTML reports in the packagePath directory.
func CreateReportsPackage(contentDir string, staticHTMLContentDir string, packagePath string) error {
	outputFile := path.Join(packagePath, "qmstr-reports.tar.bz2")
	cmd := exec.Command("tar", "cfj", outputFile, staticHTMLContentDir)
	cmd.Dir = contentDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("NOTE - output:\n%v", string(output[:]))
		return fmt.Errorf("error creating package of QMSTR reports: %v", err)
	}
	log.Printf("generated package of QMSTR reports at %v", outputFile)
	return nil
}

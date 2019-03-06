package htmlreporter

import (
	"bufio"
	"context"
	"errors"
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

const (
	minimumHugoVersion = "0.32"
	// ModuleName is used across QMSTR to reference this module
	ModuleName         = "reporter-html"
	themeDirectoryName = "theme"
	cacheVersion       = "0.2"
)

// HTMLReporter is the context of the HTML reporter module
type HTMLReporter struct {
	workingDir     string
	sharedDataDir  string
	Keep           bool
	baseURL        string
	siteData       *reporting.SiteData
	packageDir     string
	cacheDir       string
	enableWarnings bool
	enableErrors   bool
}

// Configure sets up the working directory for this reporting run.
// It is part of the ReporterModule interface.
func (r *HTMLReporter) Configure(config map[string]string) error {
	if outDir, ok := config["outputdir"]; ok {
		r.packageDir = outDir
	} else {
		return fmt.Errorf("no output directory configured")
	}

	if val, ok := config["keep"]; ok && val == "true" {
		r.Keep = true
	}

	if val, ok := config["baseurl"]; ok {
		r.baseURL = val
	} else {
		r.baseURL = "file:///var/lib/qmstr/reports"
	}

	if siteData, err := reporting.GetSiteDataFromConfiguration(config); err == nil {
		r.siteData = siteData
	} else {
		return fmt.Errorf("missing or incomplete site data in configuration: %v", err)
	}

	if cacheDir, ok := config["cachedir"]; ok {
		r.cacheDir = cacheDir
	} else {
		return fmt.Errorf("no cache directory specified - it is required")
	}

	if enable, ok := config["warnings"]; ok && enable == "true" {
		r.enableWarnings = true
	}

	if enable, ok := config["errors"]; ok && enable == "true" {
		r.enableErrors = true
	}

	detectedVersion, err := DetectHugoAndVerifyVersion()
	if err != nil {
		return fmt.Errorf("error generating HTML reports: %v", err)
	}
	log.Printf("detected beautiful Hugo version %v", detectedVersion)

	r.sharedDataDir, err = reporting.DetectModuleSharedDataDirectory(ModuleName)
	if err != nil {
		return fmt.Errorf("cannot identify QMSTR shared data directory: %v", err)
	}
	r.workingDir, err = r.CreateHugoWorkingDirectory(r.sharedDataDir, r.baseURL)
	if err != nil {
		return fmt.Errorf("error preparing temporary Hugo working directory: %v", err)
	}
	log.Printf("created temporary Hugo working directory in %v", r.workingDir)

	return nil
}

// Report generates the actual reports.
// It is part of the ReporterModule interface.
func (r *HTMLReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient) error {
	packageNode, err := cserv.GetPackageNode(context.Background(), &service.PackageNode{})
	if err != nil {
		return fmt.Errorf("could not get package node: %v", err)
	}

	bom, err := rserv.GetBOM(context.Background(), &service.BOMRequest{Warnings: r.enableWarnings, Errors: r.enableErrors})
	if err != nil {
		return err
	}
	log.Printf("%v", bom)

	if err := r.CreatePackageLevelReports(bom, cserv, rserv); err != nil {
		return fmt.Errorf("error generating package level report: %v", err)
	}
	log.Printf("HTML reporter: created report for %v", packageNode.Name)
	return nil
}

// PostReport is called after the report has been generated.
// It is part of the ReporterModule interface.
func (r *HTMLReporter) PostReport() error {
	if !r.Keep {
		defer r.cleanup()
	}

	staticHTMLContentDir, err := CreateStaticHTML(r.workingDir)
	if err != nil {
		return fmt.Errorf("error generating reports: %v", err)
	}
	log.Printf("QMSTR reports generated in %v", staticHTMLContentDir)

	// Create the reports package (for publishing, etc).
	if err := CreateReportsPackage(r.workingDir, staticHTMLContentDir, r.packageDir); err != nil {
		return fmt.Errorf("error packaging report to %v: %v", r.packageDir, err)
	}
	log.Printf("QMSTR reports package generated in %v", r.packageDir)
	return nil
}

func (r *HTMLReporter) cleanup() {
	log.Printf("deleting temporary Hugo working directory in %v", r.workingDir)
	if err := os.RemoveAll(r.workingDir); err != nil {
		log.Printf("warning - error deleting temporary Hugo working directory in %v: %v", r.workingDir, err)
	}
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
		return version, nil
	}

	re = regexp.MustCompile("Site Generator v(.+) .+/.+")
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

func checkCache(cacheDir string, cacheversion string) error {
	cacheVersionFileName := path.Join(cacheDir, ".cachever")
	version, err := ioutil.ReadFile(cacheVersionFileName)
	if err != nil {
		return fmt.Errorf("could not read cache version file %s: %v", cacheVersionFileName, err)
	}
	if string(version) != cacheversion {
		return errors.New("cache version mismatch")
	}
	return nil
}

func prepareCache(cacheDir string, workDir string) error {
	err := checkCache(cacheDir, cacheVersion)
	if err != nil {
		log.Printf("invalid cache found: %v", err)
		err = os.RemoveAll(cacheDir)
		if err != nil {
			return err
		}
	}

	for _, folder := range []string{"content", "data"} {
		cachePath := path.Join(cacheDir, folder)
		// create the folder in the cache directory (it may exist, that is not an error):
		if err := os.MkdirAll(cachePath, 0700); err != nil {
			return fmt.Errorf("error creating folder \"%s\" in cache directory: %v", folder, err)
		}
		// symlink it to the temp working directory:
		link := path.Join(workDir, folder)
		if err := os.Symlink(cachePath, link); err != nil {
			return fmt.Errorf("unable to link folder \"%s\" from the cache directory to the temporary working directory: %v", folder, err)
		}
	}

	// write cache version file
	cacheVersionFileName := path.Join(cacheDir, ".cachever")
	err = ioutil.WriteFile(cacheVersionFileName, []byte(cacheVersion), 0644)
	if err != nil {
		return fmt.Errorf("failed to write cache version file: %v", err)
	}

	return nil
}

//CreateHugoWorkingDirectory creates a temporary directory with the directory structure required to run Hugo
func (r *HTMLReporter) CreateHugoWorkingDirectory(sharedDataDir string, baseURL string) (string, error) {
	tmpWorkDir, err := ioutil.TempDir("", "qmstr-")
	if err != nil {
		return tmpWorkDir, fmt.Errorf("error creating temporary Hugo working directory: %v", err)
	}
	themeDirectory := path.Join(sharedDataDir, themeDirectoryName)
	// populate working directory with a site template and the theme
	skeletonDir := "skeleton"
	templateDir := "templates"
	linksfromTo := make(map[string]string)

	if err := os.MkdirAll(path.Join(tmpWorkDir, "themes"), 0700); err != nil {
		return "", fmt.Errorf("error creating folder in site skeleton: %v", err)
	}

	// save content and data in the cache directory:
	if err := prepareCache(r.cacheDir, tmpWorkDir); err != nil {
		return "", fmt.Errorf("Failed to prepare cache : %v", err)
	}

	linksfromTo[path.Join(sharedDataDir, skeletonDir, "archetypes")] = "archetypes"
	linksfromTo[path.Join(sharedDataDir, skeletonDir, "layouts")] = "layouts"
	linksfromTo[path.Join(sharedDataDir, skeletonDir, "static")] = "static"
	linksfromTo[themeDirectory] = "themes/qmstr-theme"

	for from, to := range linksfromTo {
		cmd := exec.Command("ln", "-s", from, to)
		cmd.Dir = tmpWorkDir
		if output, err := cmd.CombinedOutput(); err != nil {
			log.Printf("NOTE - output:\n%v", string(output[:]))
			return "", fmt.Errorf("error generating links to site skeleton and theme: %v", err)
		}
	}
	log.Printf("generated Hugo site template in %v", tmpWorkDir)
	// TODO: Export this type, populate the structure from the reporter configuration:
	type Configuration struct {
		Title, TitleEn, BaseURL string
	}
	configuration := Configuration{
		"Quartermaster Compliance Report",
		"Quartermaster Compliance Report",
		baseURL,
	}
	configTomlInPath := path.Join(sharedDataDir, templateDir, "config.toml.in")
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
// It returns the directory with the generated content (relative to contentDir) and/or an error.
func CreateStaticHTML(workingdir string) (string, error) {
	outputDir := "reports"
	outputPath := path.Join(workingdir, outputDir)
	cmd := exec.Command("hugo", "-v", "-d", outputPath)
	cmd.Dir = workingdir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("NOTE - output:\n%v", string(output[:]))
		return "", fmt.Errorf("error generating static HTML reports: %v", err)
	}
	log.Printf("generated static HTML reports in %v", outputPath)
	return outputDir, nil
}

// CreateReportsPackage creates a tarball of the static HTML reports in the packagePath directory.
func CreateReportsPackage(workingDir string, contentDir string, packagePath string) error {
	outputFile := path.Join(packagePath, "qmstr-reports.tar.bz2")
	cmd := exec.Command("tar", "cfj", outputFile, "-C", workingDir, contentDir)
	cmd.Dir = workingDir
	if output, err := cmd.CombinedOutput(); err != nil {
		log.Printf("NOTE - output:\n%v", string(output[:]))
		return fmt.Errorf("error creating package of QMSTR reports: %v", err)
	}
	log.Printf("generated package of QMSTR reports at %v", outputFile)
	return nil
}

func applyTemplate(templatePath string, data interface{}, target string) error {
	funcMap := template.FuncMap{
		"summary":   reporting.CommitMessageSummary,
		"shortenId": reporting.ShortenedVersionIdentifier,
	}

	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("unable to read template file \"%s\"", templatePath)
	}
	t := template.Must(template.New(templatePath).Funcs(funcMap).Delims("{{{", "}}}").Parse(string(templateData)))

	targetFile, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("unable to create target file \"%s\v\"", target)
	}
	defer targetFile.Close()
	writer := bufio.NewWriter(targetFile)

	if err := t.Execute(writer, data); err != nil {
		return fmt.Errorf("error applying data to template: %v", err)
	}
	writer.Flush()
	log.Printf("generated configuration file %v", target)

	return nil
}

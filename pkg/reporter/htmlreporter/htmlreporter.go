package htmlreporter

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/QMSTR/qmstr/pkg/service"
	version "github.com/hashicorp/go-version"
)

const (
	minimumHugoVersion = "0.32"
	// ModuleName is used across QMSTR to reference this module
	ModuleName         = "reporter-html"
	themeDirectoryName = "theme"
)

// HTMLReporter is the context of the HTML reporter module
type HTMLReporter struct {
	workingDir    string
	sharedDataDir string
	Keep          bool
	baseURL       string
	siteData      SiteData
	packageDir    string
	cacheDir      string
}

// Configure sets up the working directory for this reporting run.
// It is part of the ReporterModule interface.
func (r *HTMLReporter) Configure(config map[string]string) error {
	if outDir, ok := config["outputdir"]; ok {
		r.packageDir = outDir
	} else {
		return errors.New("no output directory configured")
	}

	if val, ok := config["keep"]; ok && val == "true" {
		r.Keep = true
	}

	if val, ok := config["baseurl"]; ok {
		r.baseURL = val
	} else {
		r.baseURL = "file:///var/lib/qmstr/reports"
	}

	if sitePro, ok := config["siteprovider"]; ok {
		r.siteData = SiteData{Provider: sitePro}
	} else {
		r.siteData = SiteData{Provider: "The Site Provider"}
	}

	if cacheDir, ok := config["cachedir"]; ok {
		r.cacheDir = cacheDir
	}

	detectedVersion, err := DetectHugoAndVerifyVersion()
	if err != nil {
		return fmt.Errorf("error generating HTML reports: %v", err)
	}
	log.Printf("detected beautiful Hugo version %v", detectedVersion)

	r.sharedDataDir, err = DetectModuleSharedDataDirectory(ModuleName)
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

// TEMP: until Report is called with the Package node:
var once = false

// Report generates the actual reports.
// It is part of the ReporterModule interface.
func (r *HTMLReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient, session string) error {
	packageNode, err := cserv.GetPackageNode(context.Background(), &service.PackageRequest{Session: session})
	if err != nil {
		return fmt.Errorf("could not get package node: %v", err)
	}

	licenses, err := rserv.GetInfoData(context.Background(), &service.InfoDataRequest{RootID: packageNode.Targets[0].Uid, Infotype: "license", Datatype: "spdxIdentifier"})
	if err != nil {
		return err
	}
	log.Printf("Licenses: %v", licenses.Data)

	authors, err := rserv.GetInfoData(context.Background(), &service.InfoDataRequest{RootID: packageNode.Targets[0].Uid, Infotype: "copyright", Datatype: "author"})
	if err != nil {
		return err
	}
	log.Printf("Authors: %v", authors.Data)

	if !once {
		once = true
		if err := r.CreatePackageLevelReports(packageNode); err != nil {
			return fmt.Errorf("error generating package level report: %v", err)
		}
	}
	log.Printf("(r *HTMLReporter) Report: %v", packageNode.Name)
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

// DetectSharedDataDirectory detects the shared data directory for all of QMSTR.
// It looks for /usr/share/qmstr, /usr/local/share/qmstr and /opt/share/qmstr, in that order.
// TODO this function should be refactored to be used across all modules
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
// TODO this function should be refactored to be used across all modules
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
	if len(r.cacheDir) == 0 {
		return "", fmt.Errorf("no cache directory specified - it is required")
	}

	for _, folder := range []string{"content", "data"} {
		cachePath := path.Join(r.cacheDir, folder)
		// create the folder in the cache directory (it may exist, that is not an error):
		if err := os.MkdirAll(cachePath, 0700); err != nil {
			return "", fmt.Errorf("error creating folder \"%s\" in cache directory: %v", folder, err)
		}
		// symlink it to the temp working directory:
		link := path.Join(tmpWorkDir, folder)
		if err := os.Symlink(cachePath, link); err != nil {
			return "", fmt.Errorf("unable to link folder \"%s\" from the cache directory to the temporary working directory: %v", folder, err)
		}
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

// CreateReporCreateReportsPackagetsPackage creates a tarball of the static HTML reports in the packagePath directory.
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

// SiteData contains information about this Quartermaster site.
type SiteData struct {
	Provider string // the responsible entity running the site
}

// PackageData is the package metadata that the report will visualize.
// PackageData is expected to stay more or less constant across versions of the package.
// oc... refers to OpenChain related fields
type PackageData struct {
	PackageName         string   // The package name, e.g. "CURL" or "Linux"
	Vendor              string   // Name of the entity distributing this package
	OcFossLiaison       string   // Name of the FOSS liaison function
	OcComplianceContact string   // Email address acting as the general FOSS compliance contact for the vendor
	Site                SiteData // The site this page is associated with
}

// RevisionData contains metadata about a specific revision.
type RevisionData struct {
	VersionIdentifier string      // Usually a Git hash, but any string can be used
	ChangeDateTime    string      // The change timestamp
	Author            string      // The author of the change
	Package           PackageData // The package this version is associated with.
}

// CreatePackageLevelReports creates the top level report about the package.
func (r *HTMLReporter) CreatePackageLevelReports(packageNode *service.PackageNode) error {
	packageData := PackageData{packageNode.Name, "Vendor", "FossLiaison", "Compliance contact email", r.siteData}
	revisionData := RevisionData{"a3ca6e98ab6ca4be5d74052efa97b2d3f710dd39", "2017-11-06 14:35", "Jonas Oberg", packageData}

	ps := reflect.ValueOf(&packageData)
	s := ps.Elem()
	for _, inode := range packageNode.AdditionalInfo {
		if inode.Type == "metadata" {
			for _, dnode := range inode.DataNodes {
				f := s.FieldByName(dnode.Type)
				if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
					f.SetString(dnode.Data)
				}
			}
		}
		if inode.Type == "Revision" {
			for _, dnode := range inode.DataNodes {
				switch dnode.Type {
				case "AuthorName":
					revisionData.Author = dnode.Data
				case "Description":
					revisionData.VersionIdentifier = dnode.Data
				case "CommitterDate":
					revisionData.ChangeDateTime = dnode.Data
				}
			}
		}
	}

	log.Printf("Using revision %v", revisionData)

	dataDirectory := path.Join(r.workingDir, "data")
	contentDirectory := path.Join(r.workingDir, "content")
	packageContentDirectory := path.Join(contentDirectory, packageData.PackageName)
	versionContentDirectory := path.Join(packageContentDirectory, revisionData.VersionIdentifier)
	packageDirectory := path.Join(dataDirectory, packageData.PackageName)
	versionDirectory := path.Join(packageDirectory, revisionData.VersionIdentifier)
	if err := os.MkdirAll(versionDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("error creating package metadata directory: %v", err)
	}
	packageJSON, err := json.Marshal(packageData)
	if err != nil {
		return fmt.Errorf("error generating JSON representation of package metadata: %v", err)
	}
	packageDataFile := path.Join(packageDirectory, "data.json")
	if err := ioutil.WriteFile(packageDataFile, packageJSON, 0644); err != nil {
		return fmt.Errorf("error creating JSON package metadata file: %v", err)
	}

	// create content directories for package and version:
	if err := os.MkdirAll(versionContentDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("error creating content directories: %v", err)
	}
	// generate top-level site data:
	{
		templatePath := path.Join(r.sharedDataDir, "templates", "site-index.md")
		outputPath := path.Join(contentDirectory, "_index.md")
		if err := applyTemplate(templatePath, r.siteData, outputPath); err != nil {
			return fmt.Errorf("error creating site index: %v", err)
		}
	}
	// generate content/<package>/_index.md
	{
		templatePath := path.Join(r.sharedDataDir, "templates", "package-index.md")
		outputPath := path.Join(packageContentDirectory, "_index.md")
		if err := applyTemplate(templatePath, packageData, outputPath); err != nil {
			return fmt.Errorf("error creating package content: %v", err)
		}
	}
	// generate content/<package>/<version>/_index.md
	{
		templatePath := path.Join(r.sharedDataDir, "templates", "version-index.md")
		outputPath := path.Join(versionContentDirectory, "_index.md")
		if err := applyTemplate(templatePath, revisionData, outputPath); err != nil {
			return fmt.Errorf("error creating version index page: %v", err)
		}
	}
	revisionJSON, err := json.Marshal(revisionData)
	if err != nil {
		return fmt.Errorf("error generating JSON representation of revision metadata: %v", err)
	}
	versionDataFile := path.Join(versionDirectory, "data.json")
	if err := ioutil.WriteFile(versionDataFile, revisionJSON, 0644); err != nil {
		return fmt.Errorf("error creating JSON version data file: %v", err)
	}
	return nil
}

func applyTemplate(templatePath string, data interface{}, target string) error {

	templateData, err := ioutil.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("unable to read template file \"%s\"", templatePath)
	}
	t := template.Must(template.New(templatePath).Delims("{{{", "}}}").Parse(string(templateData)))

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

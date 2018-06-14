package htmlreporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/QMSTR/qmstr/pkg/service"
)

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
	Message           string      // The commit message
	Summary           string      // The short form of the commit message
	Package           PackageData // The package this version is associated with.
}

// CreatePackageLevelReports creates the top level report about the package.
func (r *HTMLReporter) CreatePackageLevelReports(packageNode *service.PackageNode) error {
	packageData := PackageData{packageNode.Name, "Vendor", "FossLiaison", "Compliance contact email", r.siteData}
	revisionData := RevisionData{"(SHA)", "(commit datetime)", "(author)", "(commit message)", "(commit summary)", packageData}

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
				case "CommitMessage":
					revisionData.Message = dnode.Data
				case "CommitID":
					log.Printf("WARN: using CommitID instead of description this can be misleading as it does not cover not commited changes")
					revisionData.VersionIdentifier = dnode.Data
				case "CommitterDate":
					revisionData.ChangeDateTime = dnode.Data
				}
			}
		}
	}

	revisionData.Summary = commitMessageSummary(revisionData.Message)
	log.Printf("Using revision %v: %s", revisionData.VersionIdentifier[:8], revisionData.Summary)

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

	if len(packageNode.Targets) == 0 {
		log.Printf("Note: package node contains no targets, please verify that targets got built")
	}

	for _, node := range packageNode.Targets {
		if err := r.CreateTargetLevelReports(node); err != nil {
			return fmt.Errorf("error creating target report for %s: %v", node.Name, err)
		}
	}

	return nil
}

// commitMessageSummary returns the summary of the commit message according to the usual guidelines
// (see https://chris.beams.io/posts/git-commit/, "Limit the subject line to 50 characters")
func commitMessageSummary(message string) string {
	lines := strings.Split(message, "\n")
	if len(lines) == 0 {
		return ""
	}
	summary := strings.TrimSpace(lines[0])
	if len(summary) > 50 {
		summary = fmt.Sprintf("%s...", summary[47:])
	}
	return summary
}

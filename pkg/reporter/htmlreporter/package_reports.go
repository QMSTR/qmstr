package htmlreporter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/QMSTR/qmstr/pkg/reporting"
	"github.com/QMSTR/qmstr/pkg/service"
)

// CreatePackageLevelReports creates the top level report about the package.
func (r *HTMLReporter) CreatePackageLevelReports(packageNode *service.PackageNode, cserv service.ControlServiceClient, rserv service.ReportServiceClient) error {
	packageData := reporting.PackageData{packageNode.Name, "Vendor", "FossLiaison", "Compliance contact email", r.siteData}
	revisionData, err := reporting.GetRevisionData(packageNode, packageData)
	log.Printf("Using revision %v: %s", revisionData.VersionIdentifierShort, revisionData.Summary)

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
		if err := r.CreateTargetLevelReports(node, cserv, rserv); err != nil {
			return fmt.Errorf("error creating target report for %s: %v", node.Name, err)
		}
	}

	return nil
}

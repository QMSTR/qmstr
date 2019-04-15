package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

// CreatePackageLevelReports creates the top level report about the package.
func (r *HTMLReporter) CreatePackageLevelReports(bom *service.BOM, cserv service.ControlServiceClient, rserv service.ReportServiceClient) error {
	packageData := reporting.GetPackageData(bom, r.siteData)
	revisionData := reporting.GetRevisionData(bom, packageData)

	// remove remotes/<remotename>
	revisionData.VersionIdentifier = filepath.Base(revisionData.VersionIdentifier)

	log.Printf("Using revision %v: %s", reporting.ShortenedVersionIdentifier(revisionData.VersionIdentifier), reporting.CommitMessageSummary(revisionData.Message))

	contentDirectory := path.Join(r.workingDir, "content")
	packageContentDirectory := path.Join(contentDirectory, packageData.PackageName)
	versionContentDirectory := path.Join(packageContentDirectory, revisionData.VersionIdentifier)
	buildconfigContentDirectory := path.Join(versionContentDirectory, packageData.BuildConfig)

	dataDirectory := path.Join(r.workingDir, "data")
	packageDirectory := path.Join(dataDirectory, packageData.PackageName)
	versionDirectory := path.Join(packageDirectory, revisionData.VersionIdentifier)
	buildconfigDirectory := path.Join(versionDirectory, packageData.BuildConfig)

	if err := os.MkdirAll(buildconfigDirectory, os.ModePerm); err != nil {
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

	// create content directories for package, version and build config:
	if err := os.MkdirAll(buildconfigContentDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("error creating content directories: %v", err)
	}
	// generate top-level site data:
	// TODO this needs to be refactored out of the *package* level reports
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

	type VersionData struct {
		reporting.RevisionData
		BuildConfig string
	}

	versionData := VersionData{*revisionData, packageData.BuildConfig}
	// generate content/<package>/<version>/_index.md
	{
		templatePath := path.Join(r.sharedDataDir, "templates", "version-index.md")
		outputPath := path.Join(versionContentDirectory, "_index.md")
		if err := applyTemplate(templatePath, versionData, outputPath); err != nil {
			return fmt.Errorf("error creating version index page: %v", err)
		}
	}
	// generate content/<package>/<version>/<buildconfig>/_index.md
	{
		templatePath := path.Join(r.sharedDataDir, "templates", "buildconfig-index.md")
		outputPath := path.Join(buildconfigContentDirectory, "_index.md")
		if err := applyTemplate(templatePath, versionData, outputPath); err != nil {
			return fmt.Errorf("error creating buildconfig index page: %v", err)
		}
	}

	revisionJSON, err := json.Marshal(versionData)
	if err != nil {
		return fmt.Errorf("error generating JSON representation of revision metadata: %v", err)
	}
	versionDataFile := path.Join(versionDirectory, "data.json")
	if err := ioutil.WriteFile(versionDataFile, revisionJSON, 0644); err != nil {
		return fmt.Errorf("error creating JSON version data file: %v", err)
	}

	buildconfigDataFile := path.Join(buildconfigDirectory, "data.json")
	if err := ioutil.WriteFile(buildconfigDataFile, revisionJSON, 0644); err != nil {
		return fmt.Errorf("error creating JSON version data file: %v", err)
	}

	if len(bom.Targets) == 0 {
		log.Printf("Note: package node contains no targets, please verify that targets got built")
	}

	return nil
}

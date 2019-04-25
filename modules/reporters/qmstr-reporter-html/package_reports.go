package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"

	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

// CreatePackageLevelReports creates the top level report about the package.
func (r *HTMLReporter) CreatePackageLevelReports(proj *service.ProjectNode, pkg *service.PackageNode, cserv service.ControlServiceClient, rserv service.ReportServiceClient) error {
	projectData := reporting.GetProjectData(proj, r.siteData)
	packageData := reporting.GetPackageData(pkg, projectData)

	contentDirectory := path.Join(r.workingDir, "content")
	projectContentDirectory := path.Join(contentDirectory, proj.Name)
	packageContentDirectory := path.Join(projectContentDirectory, pkg.Name)
	versionContentDirectory := path.Join(packageContentDirectory, pkg.Version)

	dataDirectory := path.Join(r.workingDir, "data")
	projectDirectory := path.Join(dataDirectory, proj.Name)
	packageDirectory := path.Join(projectDirectory, pkg.Name)
	versionDirectory := path.Join(packageDirectory, pkg.Version)

	if err := os.MkdirAll(versionDirectory, os.ModePerm); err != nil {
		return fmt.Errorf("error creating package metadata directory: %v", err)
	}

	projectJSON, err := json.Marshal(projectData)
	if err != nil {
		return fmt.Errorf("error generating JSON representation of package metadata: %v", err)
	}
	projectDataFile := path.Join(projectDirectory, "data.json")
	if err := ioutil.WriteFile(projectDataFile, projectJSON, 0644); err != nil {
		return fmt.Errorf("error creating JSON package metadata file: %v", err)
	}

	// create content directories for project, package and version:
	if err := os.MkdirAll(versionContentDirectory, os.ModePerm); err != nil {
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
	// generate content/<project>/_index.md
	{
		templatePath := path.Join(r.sharedDataDir, "templates", "project-index.md")
		outputPath := path.Join(projectContentDirectory, "_index.md")
		if err := applyTemplate(templatePath, projectData, outputPath); err != nil {
			return fmt.Errorf("error creating buildconfig index page: %v", err)
		}
	}
	// generate content/<project>/<package>/_index.md
	{
		templatePath := path.Join(r.sharedDataDir, "templates", "package-index.md")
		outputPath := path.Join(packageContentDirectory, "_index.md")
		if err := applyTemplate(templatePath, packageData, outputPath); err != nil {
			return fmt.Errorf("error creating package content: %v", err)
		}
	}

	// generate content/<project>/<package>/<version>/_index.md
	{
		templatePath := path.Join(r.sharedDataDir, "templates", "version-index.md")
		outputPath := path.Join(versionContentDirectory, "_index.md")
		if err := applyTemplate(templatePath, packageData, outputPath); err != nil {
			return fmt.Errorf("error creating version index page: %v", err)
		}
	}

	packageJSON, err := json.Marshal(packageData)
	if err != nil {
		return fmt.Errorf("error generating JSON representation of revision metadata: %v", err)
	}

	packageDataFile := path.Join(packageDirectory, "data.json")
	if err := ioutil.WriteFile(packageDataFile, packageJSON, 0644); err != nil {
		return fmt.Errorf("error creating JSON version data file: %v", err)
	}

	versionDataFile := path.Join(versionDirectory, "data.json")
	if err := ioutil.WriteFile(versionDataFile, packageJSON, 0644); err != nil {
		return fmt.Errorf("error creating JSON version data file: %v", err)
	}

	if len(pkg.Targets) == 0 {
		log.Printf("Note: package node contains no targets, please verify that targets got built")
	}

	return nil
}

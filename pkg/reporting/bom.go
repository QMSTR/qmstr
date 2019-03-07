package reporting

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/QMSTR/qmstr/pkg/service"
)

func GetBOM(proj *service.ProjectNode, pkg *service.PackageNode, enableWarnings bool, enableErrors bool) (*service.BOM, error) {
	packageInfo, err := getPackageInfo(pkg, proj)
	if err != nil {
		return nil, fmt.Errorf("Failed to get package information : %v", err)
	}
	revisionInfo, err := getRevisionInfo(pkg)
	if err != nil {
		return nil, fmt.Errorf("Failed to get revision information : %v", err)
	}
	warnings, errors := []string{}, []string{}
	if enableWarnings {
		warnings = getDataByInfoType(pkg, "warning")
	}
	if enableErrors {
		errors = getDataByInfoType(pkg, "error")
	}
	bom := service.BOM{
		PackageInfo: packageInfo,
		VersionInfo: revisionInfo,
		Warnings:    warnings,
		Errors:      errors,
		Targets:     getTargetsInfo(pkg),
	}
	return &bom, nil
}

func getDataByInfoType(packageNode *service.PackageNode, infoType string) []string {
	data := []string{}
	for _, target := range packageNode.Targets {
		for _, depNode := range target.DerivedFrom {
			for _, info := range depNode.AdditionalInfo {
				if info.Type == infoType {
					for _, d := range info.DataNodes {
						data = append(data, d.Data)
					}
				}
			}
		}

	}
	return data
}

func getPackageInfo(pkg *service.PackageNode, proj *service.ProjectNode) (*service.PackageInformation, error) {
	packageInfo := service.PackageInformation{
		Name:                pkg.GetName(),
		BuildConfig:         pkg.GetBuildConfig(),
		Vendor:              "Vendor",
		OcFossLiaison:       "FossLiaison",
		OcComplianceContact: "Compliance contact email"}
	ps := reflect.ValueOf(&packageInfo)
	s := ps.Elem()
	for _, inode := range proj.AdditionalInfo {
		if inode.Type == "metadata" {
			for _, dnode := range inode.DataNodes {
				f := s.FieldByName(getFieldName(dnode.Type))
				if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
					f.SetString(dnode.Data)
				}
			}
		}
	}
	return &packageInfo, nil
}

func getFieldName(name string) string {
	fields := strings.Split(name, "_")
	for idx := range fields {
		fields[idx] = strings.Title(fields[idx])
	}
	return strings.Join(fields, "")
}

func getEmptyRevision() *service.Revision {
	return &service.Revision{
		Author:     &service.Person{Name: "Unknown Author", Email: "unknown@unknown.local"},
		AuthorDate: "",
		CommitDate: "",
		Committer:  &service.Person{Name: "Unknown committer", Email: "unknown@unknown.local"},
		Id:         "unknown id",
		Message:    "",
	}
}

func getRevisionInfo(packageNode *service.PackageNode) (*service.Revision, error) {
	revisionData := getEmptyRevision()

	for _, inode := range packageNode.AdditionalInfo {
		if inode.Type == "Revision" {
			for _, dnode := range inode.DataNodes {
				switch dnode.Type {
				case "AuthorName":
					revisionData.Author.Name = dnode.Data
				case "AuthorEmail":
					revisionData.Author.Email = dnode.Data
				case "CommitMessage":
					revisionData.Message = dnode.Data
				case "Description":
					if dnode.Data != "" {
						revisionData.Id = dnode.Data
					}
				case "CommitID":
					if revisionData.Id == "" {
						log.Printf("WARN: using CommitID instead of description this can be misleading as it does not cover uncommited changes")
						revisionData.Id = dnode.Data
					}
				case "CommitterDate":
					revisionData.CommitDate = dnode.Data
				}
			}
		}
	}
	return revisionData, nil
}

func getTargetsInfo(packageNode *service.PackageNode) []*service.Target {
	targetsResult := []*service.Target{}
	for _, target := range packageNode.Targets {
		tmpTarget := &service.Target{
			Name:         target.Name,
			Sha1:         target.Hash,
			Path:         target.Path,
			Sources:      getSources(target),
			Dependencies: getDependencies(target),
		}
		targetsResult = append(targetsResult, tmpTarget)
	}
	return targetsResult
}

func getSources(filenode *service.FileNode) []*service.Source {
	retSources := []*service.Source{}
	for _, depNode := range filenode.DerivedFrom {
		if depNode.FileType == service.FileNode_SOURCE {
			tmpSource := &service.Source{
				File:    depNode.Path,
				License: getLicense(depNode),
				Authors: getAuthors(depNode),
			}
			retSources = append(retSources, tmpSource)
		} else {
			retSources = append(retSources, getSources(depNode)...)
		}
	}
	return retSources
}

func getDependencies(filenode *service.FileNode) []*service.Dependency {
	retDeps := []*service.Dependency{}
	for _, depNode := range filenode.DerivedFrom {
		switch depNode.FileType {
		case service.FileNode_TARGET:
			tmpDep := &service.Dependency{
				Filepath: depNode.Path,
				Name:     depNode.Name,
			}
			retDeps = append(retDeps, tmpDep)
		default:
			retDeps = append(retDeps, getDependencies(depNode)...)
		}
	}
	return retDeps
}

func getLicense(fileNode *service.FileNode) *service.License {
	for _, info := range fileNode.AdditionalInfo {
		if info.Type == "license" {
			lic := service.License{}
			for _, d := range info.DataNodes {
				switch d.Type {
				case "spdxIdentifier":
					lic.SpdxIdentifier = d.Data
				case "name":
					lic.Name = d.Data
				}
			}
			return &lic
		}
	}
	return &service.License{Name: "Unknown", SpdxIdentifier: "Unknown"}
}

func getAuthors(fileNode *service.FileNode) []*service.Person {
	authors := []*service.Person{}
	for _, info := range fileNode.AdditionalInfo {
		if info.Type == "copyright" {
			author := service.Person{}
			for _, d := range info.DataNodes {
				switch d.Type {
				case "author":
					author.Name = d.Data
				}
			}
			authors = append(authors, &author)
		}
	}
	return authors
}

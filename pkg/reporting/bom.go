package reporting

import (
	"fmt"
	"log"
	"reflect"

	"github.com/QMSTR/qmstr/pkg/service"
)

func GetBOM(pkg *service.PackageNode) (*service.BOM, error) {
	packageInfo, err := getPackageInfo(pkg)
	if err != nil {
		return nil, fmt.Errorf("Failed to get package information : %v", err)
	}
	revisionInfo, err := getRevisionInfo(pkg)
	if err != nil {
		return nil, fmt.Errorf("Failed to get revision information : %v", err)
	}
	bom := service.BOM{
		PackageInfo: packageInfo,
		VersionInfo: revisionInfo,
		Targets:     getTargetsInfo(pkg),
	}
	return &bom, nil
}

func getPackageInfo(pkg *service.PackageNode) (*service.PackageInformation, error) {
	packageInfo := service.PackageInformation{
		Name:                pkg.GetName(),
		Vendor:              "Vendor",
		OcFossLiaison:       "FossLiaison",
		OcComplianceContact: "Compliance contact email"}
	ps := reflect.ValueOf(&packageInfo)
	s := ps.Elem()
	for _, inode := range pkg.AdditionalInfo {
		if inode.Type == "metadata" {
			for _, dnode := range inode.DataNodes {
				f := s.FieldByName(dnode.Type)
				if f.IsValid() && f.CanSet() && f.Kind() == reflect.String {
					f.SetString(dnode.Data)
				}
			}
		}
	}
	return &packageInfo, nil
}

func getRevisionInfo(packageNode *service.PackageNode) (*service.Revision, error) {
	revisionData := service.Revision{Author: &service.Person{}}
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
				case "CommitID":
					log.Printf("WARN: using CommitID instead of description this can be misleading as it does not cover not commited changes")
					revisionData.Id = dnode.Data
				case "CommitterDate":
					revisionData.CommitDate = dnode.Data
				}
			}
		}
	}
	revisionData.Summary = CommitMessageSummary(revisionData.Message)
	revisionData.ShortId = ShortenedVersionIdentifier(revisionData.Id)
	return &revisionData, nil
}

func getTargetsInfo(packageNode *service.PackageNode) []*service.Target {
	targetsResult := []*service.Target{}
	for _, target := range packageNode.Targets {
		tmpTarget := &service.Target{
			Name:    target.Name,
			Sha1:    target.Hash,
			Sources: getSources(target),
		}
		targetsResult = append(targetsResult, tmpTarget)
	}
	return targetsResult
}

func getSources(filenode *service.FileNode) []*service.Source {
	retSources := []*service.Source{}
	for _, depNode := range filenode.DerivedFrom {
		if depNode.Type == "sourcecode" {
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

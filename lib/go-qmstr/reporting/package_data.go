package reporting

import (
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

// PackageData contains metadata about a specific package.
type PackageData struct {
	Name            string
	Version         string // Usually a Git hash, but any string can be used
	LicenseDeclared string
	Targets         []*service.FileNode
	Authors         []string
	Project         string // The project this package is associated with.
}

func GetPackageData(pkg *service.PackageNode, projData *ProjectData) *PackageData {
	packageData := PackageData{
		Name:            pkg.GetName(),
		Version:         pkg.GetVersion(),
		LicenseDeclared: "NO License declared",
		Project:         projData.Name,
	}
	authors := []string{}
	targets := []*service.FileNode{}
	for _, fileNode := range pkg.GetTargets() {
		targets = append(targets, fileNode)
		for _, info := range fileNode.AdditionalInfo {
			if info.Type == "copyright" {
				for _, d := range info.DataNodes {
					switch d.Type {
					case "author":
						authors = append(authors, d.Data)
					}
				}
			}
		}
	}
	packageData.Authors = authors
	packageData.Targets = targets

	return &packageData
}

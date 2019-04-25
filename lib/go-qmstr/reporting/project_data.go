package reporting

import (
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"reflect"
	"strings"
)

// ProjectData is the project metadata that the report will visualize.
// ProjectData is expected to stay more or less constant across versions of the package.
// oc... refers to OpenChain related fields
type ProjectData struct {
	Name                string    // The project name, e.g. "CURL" or "Linux"
	Vendor              string    // Name of the entity distributing this package
	OcFossLiaison       string    // Name of the FOSS liaison function
	OcComplianceContact string    // Email address acting as the general FOSS compliance contact for the vendor
	Site                *SiteData // The site this page is associated with
}

func GetProjectData(proj *service.ProjectNode, siteData *SiteData) *ProjectData {
	projectData := ProjectData{
		Name:                proj.GetName(),
		Vendor:              "Vendor",
		OcFossLiaison:       "FossLiaison",
		OcComplianceContact: "Compliance contact email",
		Site:                siteData,
	}

	ps := reflect.ValueOf(&projectData)
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
	return &projectData
}

func getFieldName(name string) string {
	fields := strings.Split(name, "_")
	for idx := range fields {
		fields[idx] = strings.Title(fields[idx])
	}
	return strings.Join(fields, "")
}

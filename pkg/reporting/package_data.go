package reporting

import (
	"reflect"

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

// GetPackageData extracts the package data from the given PackageNode
func GetPackageData(packageNode *service.PackageNode, siteData SiteData) PackageData {
	packageData := PackageData{packageNode.Name, "Vendor", "FossLiaison", "Compliance contact email", siteData}
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
	}
	return packageData
}

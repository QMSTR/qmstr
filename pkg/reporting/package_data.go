package reporting

import (
	"github.com/QMSTR/qmstr/pkg/service"
)

// PackageData is the package metadata that the report will visualize.
// PackageData is expected to stay more or less constant across versions of the package.
// oc... refers to OpenChain related fields
type PackageData struct {
	PackageName         string // The package name, e.g. "CURL" or "Linux"
	BuildConfig         string // The name of the build configuration e.g. "amd_64_something"
	Vendor              string // Name of the entity distributing this package
	OcFossLiaison       string // Name of the FOSS liaison function
	OcComplianceContact string // Email address acting as the general FOSS compliance contact for the vendor
	LicenseDeclared     string
	Targets             []*service.Target
	Site                *SiteData // The site this page is associated with
}

// GetPackageData extracts the package data from the given BOM
func GetPackageData(bom *service.BOM, siteData *SiteData) *PackageData {
	return &PackageData{bom.PackageInfo.Name, bom.PackageInfo.BuildConfig, bom.PackageInfo.Vendor, bom.PackageInfo.OcFossLiaison, bom.PackageInfo.OcComplianceContact, bom.PackageInfo.LicenseDeclared, bom.Targets, siteData}
}

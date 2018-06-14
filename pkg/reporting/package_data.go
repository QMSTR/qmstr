package reporting

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

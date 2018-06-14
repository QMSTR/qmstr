package reporting

import (
	"fmt"
)

// SiteData contains information about this Quartermaster site.
type SiteData struct {
	Provider string // the responsible entity running the site
}

// GetSiteDataFromConfiguration extracts the site provider informatiom from the configuration
func GetSiteDataFromConfiguration(config map[string]string) (SiteData, error) {
	var siteData SiteData
	if sitePro, ok := config["siteprovider"]; ok {
		siteData = SiteData{Provider: sitePro}
	} else {
		siteData = SiteData{Provider: "(Site Provider)"}
		return siteData, fmt.Errorf("missing required site provider configuration (\"config/siteprovider\")")
	}

	return SiteData{}, nil
}

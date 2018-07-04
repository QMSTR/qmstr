package reporting

import (
	"fmt"
)

// SiteData contains information about this Quartermaster site.
type SiteData struct {
	Provider string // the responsible entity running the site
}

// GetSiteDataFromConfiguration extracts the site provider informatiom from the configuration
func GetSiteDataFromConfiguration(config map[string]string) (*SiteData, error) {
	var siteData *SiteData
	const key = "siteprovider"
	if sitePro, ok := config[key]; ok {
		siteData = &SiteData{Provider: sitePro}
	} else {
		return nil, fmt.Errorf("missing required site provider configuration (key \"%s\")", key)
	}
	return siteData, nil
}

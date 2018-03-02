package report

import (
	"fmt"

	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type LicenseReporter struct {
}

func NewLicenseReporter() *LicenseReporter {
	return &LicenseReporter{}
}

func (lr *LicenseReporter) Generate(nodes []*database.Node) (*service.ReportResponse, error) {
	licenses := map[string][]string{}
	for _, node := range nodes {
		for _, lic := range getLicense(node) {
			licenses[node.Path] = append(licenses[node.Path], lic.Key)
		}
	}

	result := ""
	for artifact, licenses := range licenses {
		result = fmt.Sprintf("%s\n%s\t%s", result, artifact, licenses)
	}

	ret := service.ReportResponse{Success: true, FileNodes: nil}
	return &ret, nil
}

func getLicense(node *database.Node) []*database.License {
	if len(node.License) > 0 {
		return node.License
	}
	licenseSet := map[string]*database.License{}

	for _, node := range node.DerivedFrom {
		for _, lic := range getLicense(node) {
			licenseSet[lic.Uid] = lic
		}
	}

	licenses := []*database.License{}
	for _, license := range licenseSet {
		licenses = append(licenses, license)
	}
	return licenses
}

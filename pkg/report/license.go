package report

import (
	"fmt"
	"log"

	"github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/database"
)

type LicenseReporter struct {
}

func NewLicenseReporter() *LicenseReporter {
	return &LicenseReporter{}
}

func (lr *LicenseReporter) Generate(nodes []ReportNode) (*buildservice.ReportResponse, error) {
	licenses := map[string][]string{}
	for _, node := range nodes {
		for _, lic := range getLicense(&node) {
			licenses[node.GetPath()] = append(licenses[node.GetPath()], lic.SpdxIdentifier)
		}
	}

	ret := buildservice.ReportResponse{Success: true, ResponseMessage: fmt.Sprintf("Found licenses: %v", licenses)}
	return &ret, nil
}

func getLicense(rnode *ReportNode) []database.License {
	if len(rnode.actualNode.License) > 0 {
		return rnode.actualNode.License
	}
	licenseSet := map[string]database.License{}

	derivedFromNodes := rnode.actualNode.DerivedFrom

	// circumvent empty DerivedFrom
	if len(derivedFromNodes) == 0 {
		freshNode, err := rnode.db.GetNodeByHash(rnode.actualNode.Hash, true)
		if err != nil {
			log.Fatalf("Not able to get Node for hash %s from DB", rnode.actualNode.Hash)
		}
		derivedFromNodes = freshNode.DerivedFrom
	}

	for _, node := range derivedFromNodes {
		newRNode := NewReportNode(node, rnode.db)
		for _, lic := range getLicense(&newRNode) {
			licenseSet[lic.Uid] = lic
		}
	}

	licenses := []database.License{}
	for _, license := range licenseSet {
		licenses = append(licenses, license)
	}
	return licenses
}

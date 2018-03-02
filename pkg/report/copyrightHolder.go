package report

import (
	"fmt"

	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type CopyrightHolderReporter struct {
}

func NewCopyrightHolderReporter() *CopyrightHolderReporter {
	return &CopyrightHolderReporter{}
}

func (chr *CopyrightHolderReporter) Generate(nodes []*database.Node) (*service.ReportResponse, error) {
	copyrightHolders := map[string][]string{}
	for _, node := range nodes {
		for _, cprHolder := range getCopyrightHolders(node) {
			copyrightHolders[node.Path] = append(copyrightHolders[node.Path], cprHolder.Name)
		}
	}

	result := ""
	for artifact, copyrightHolders := range copyrightHolders {
		result = fmt.Sprintf("%s\n%s\t%s", result, artifact, copyrightHolders)
	}

	ret := service.ReportResponse{Success: true, FileNodes: nil}
	return &ret, nil
}

func getCopyrightHolders(node *database.Node) []*database.CopyrightHolder {
	if len(node.CopyrightHolder) > 0 {
		return node.CopyrightHolder
	}
	copyrightHolderSet := map[string]*database.CopyrightHolder{}

	for _, node := range node.DerivedFrom {
		for _, cprHolder := range getCopyrightHolders(node) {
			copyrightHolderSet[cprHolder.Uid] = cprHolder
		}
	}

	copyrightHolders := []*database.CopyrightHolder{}
	for _, copyrightHolder := range copyrightHolderSet {
		copyrightHolders = append(copyrightHolders, copyrightHolder)
	}
	return copyrightHolders
}

package htmlreporter

import (
	"fmt"

	"github.com/QMSTR/qmstr/pkg/service"
)

// CreateTargetLevelReports creates the report for a link target
func (r *HTMLReporter) CreateTargetLevelReports(targetNode *service.FileNode) error {

	if targetNode == nil {
		return fmt.Errorf("package node contains no targets, please verify that targets got built")
	}

	// licenses, err := rserv.GetInfoData(context.Background(), &service.InfoDataRequest{RootID: packageNode.Targets[0].Uid, Infotype: "license", Datatype: "spdxIdentifier"})
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Licenses: %v", licenses.Data)

	// authors, err := rserv.GetInfoData(context.Background(), &service.InfoDataRequest{RootID: packageNode.Targets[0].Uid, Infotype: "copyright", Datatype: "author"})
	// if err != nil {
	// 	return err
	// }
	// log.Printf("Authors: %v", authors.Data)

	return fmt.Errorf("NI")
}

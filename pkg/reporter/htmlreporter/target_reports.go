package htmlreporter

import (
	"context"
	"fmt"
	"log"

	"github.com/QMSTR/qmstr/pkg/service"
)

// CreateTargetLevelReports creates the report for a link target
func (r *HTMLReporter) CreateTargetLevelReports(targetNode *service.FileNode, cserv service.ControlServiceClient, rserv service.ReportServiceClient) error {

	if targetNode == nil {
		return fmt.Errorf("package node contains no targets, please verify that targets got built")
	}

	licenses, err := rserv.GetInfoData(context.Background(), &service.InfoDataRequest{RootID: targetNode.Uid, Infotype: "license", Datatype: "spdxIdentifier"})
	if err != nil {
		return err
	}
	log.Printf("Licenses: %v", licenses.Data)

	authors, err := rserv.GetInfoData(context.Background(), &service.InfoDataRequest{RootID: targetNode.Uid, Infotype: "copyright", Datatype: "author"})
	if err != nil {
		return err
	}
	log.Printf("Authors: %v", authors.Data)

	log.Printf("NI")
	return nil
}

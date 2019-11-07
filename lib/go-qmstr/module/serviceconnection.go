package module

import (
	"log"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"google.golang.org/grpc"
)

type MasterClient struct {
	CtrlSvcClient service.ControlServiceClient
	AnaSvcClient  service.AnalysisServiceClient
	RptSvcClient  service.ReportServiceClient
}

func NewMasterClient(serviceAddress string) MasterClient {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}

	ctrlSvcClient := service.NewControlServiceClient(conn)
	anaSvcClient := service.NewAnalysisServiceClient(conn)
	rptSvcClient := service.NewReportServiceClient(conn)
	return MasterClient{CtrlSvcClient: ctrlSvcClient, AnaSvcClient: anaSvcClient, RptSvcClient: rptSvcClient}
}

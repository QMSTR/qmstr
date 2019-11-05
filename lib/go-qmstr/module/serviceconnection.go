package module

import (
	"log"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"google.golang.org/grpc"
)

type MasterClient struct {
	ctrlSvcClient service.ControlServiceClient
	anaSvcClient  service.AnalysisServiceClient
	rptSvcClient  service.ReportServiceClient
}

func NewMasterClient(serviceAddress string) *MasterClient {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}

	ctrlSvcClient := service.NewControlServiceClient(conn)
	anaSvcClient := service.NewAnalysisServiceClient(conn)
	rptSvcClient := service.NewReportServiceClient(conn)
	return &MasterClient{ctrlSvcClient: ctrlSvcClient, anaSvcClient: anaSvcClient, rptSvcClient: rptSvcClient}
}

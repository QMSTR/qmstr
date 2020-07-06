package module

import (
	"log"
	"os"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"google.golang.org/grpc"
)

type MasterClient struct {
	CtrlSvcClient service.ControlServiceClient
	AnaSvcClient  service.AnalysisServiceClient
	RptSvcClient  service.ReportServiceClient
}

func NewMasterClient(serviceAddress string) MasterClient {
	// Connect to qmstr master instance
	if len(serviceAddress) == 0 {
		serviceAddress = os.Getenv("QMSTR_MASTER")
	}
	if len(serviceAddress) == 0 {
		log.Fatalf("Error: master address not specified")
	}
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to master: %v", err)
	}

	ctrlSvcClient := service.NewControlServiceClient(conn)
	anaSvcClient := service.NewAnalysisServiceClient(conn)
	rptSvcClient := service.NewReportServiceClient(conn)
	return MasterClient{CtrlSvcClient: ctrlSvcClient, AnaSvcClient: anaSvcClient, RptSvcClient: rptSvcClient}
}

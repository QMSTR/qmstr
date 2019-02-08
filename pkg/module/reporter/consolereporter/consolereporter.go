package consolereporter

import (
	"context"
	"fmt"
	"io"

	"github.com/QMSTR/go-qmstr/service"
)

const (
	ModuleName = "reporter-console"
)

type ConsoleReporter struct {
	enableWarnings bool
}

func (r *ConsoleReporter) Configure(config map[string]string) error {
	return nil
}

func (r *ConsoleReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient) error {
	stream, err := cserv.GetDiagnosticNode(context.Background(), &service.DiagnosticNode{Type: service.DiagnosticNode_ERROR})
	if err != nil {
		return fmt.Errorf("could not get diagnostic nodes: %v", err)
	}

	for {
		diagnosticNode, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		fmt.Printf("ERROR: %s\n", diagnosticNode.Message)
	}
	return nil
}

func (r *ConsoleReporter) PostReport() error {
	return nil
}

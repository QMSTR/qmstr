package main

import (
	"context"
	"fmt"
	"io"

	"github.com/QMSTR/qmstr/lib/go-qmstr/module"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
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

func (r *ConsoleReporter) Report(masterClient *module.MasterClient) error {
	stream, err := masterClient.CtrlSvcClient.GetDiagnosticNode(context.Background(), &service.DiagnosticNode{Severity: service.DiagnosticNode_ERROR})
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

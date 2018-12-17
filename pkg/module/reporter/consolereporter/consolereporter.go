package consolereporter

import (
	"context"
	"fmt"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
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

func (r *ConsoleReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient, session string) error {
	packageNode, err := cserv.GetPackageNode(context.Background(), &service.PackageRequest{Session: session})
	if err != nil {
		return fmt.Errorf("could not get package node: %v", err)
	}

	for _, target := range packageNode.Targets {
		for _, depNode := range target.DerivedFrom {
			for _, info := range depNode.AdditionalInfo {
				if info.Type == "error" {
					for _, d := range info.DataNodes {
						fmt.Printf("ERROR: %s\n", d.Data)
					}
				}
			}
		}
	}
	return nil
}

func (r *ConsoleReporter) PostReport() error {
	return nil
}

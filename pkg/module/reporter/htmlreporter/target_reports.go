package htmlreporter

import (
	"fmt"
	"log"

	"github.com/QMSTR/qmstr/pkg/service"
)

// CreateTargetLevelReports creates the report for a link target
func (r *HTMLReporter) CreateTargetLevelReports(target *service.Target) error {

	if target == nil {
		return fmt.Errorf("package node contains no targets, please verify that targets got built")
	}

	log.Printf("Not yet implemented")

	return nil
}

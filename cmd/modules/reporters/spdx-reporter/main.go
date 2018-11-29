package main

import (
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/module/reporter/spdxreporter"
	"github.com/QMSTR/qmstr/pkg/reporting"
)

func main() {
	reporter := reporting.NewReporter(&spdxreporter.SPDXReporter{})
	if err := reporter.RunReporterModule(); err != nil {
		log.Printf("%v failed: %v\n", spdxreporter.ModuleName, err)
		os.Exit(master.ReturnReporterFailed)
	}
	log.Printf("%v completed successfully\n", spdxreporter.ModuleName)
}

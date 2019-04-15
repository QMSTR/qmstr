package main

import (
	"log"
	"os"

	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/module/reporter/htmlreporter"
	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
)

func main() {
	reporter := reporting.NewReporter(&htmlreporter.HTMLReporter{Keep: false})
	if err := reporter.RunReporterModule(); err != nil {
		log.Printf("%v failed: %v\n", htmlreporter.ModuleName, err)
		os.Exit(master.ReturnReporterFailed)
	}
	log.Printf("%v completed successfully\n", htmlreporter.ModuleName)
}

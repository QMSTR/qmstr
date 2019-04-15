package main

import (
	"log"
	"os"

	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/module/reporter/consolereporter"
	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
)

func main() {
	reporter := reporting.NewReporter(&consolereporter.ConsoleReporter{})
	if err := reporter.RunReporterModule(); err != nil {
		log.Printf("%v failed: %v\n", consolereporter.ModuleName, err)
		os.Exit(master.ReturnReporterFailed)
	}
	log.Printf("%v completed successfully\n", consolereporter.ModuleName)
}

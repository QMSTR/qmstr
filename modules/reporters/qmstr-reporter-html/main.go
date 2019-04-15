package main

import (
	"log"
	"os"

	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
)

func main() {
	reporter := reporting.NewReporter(&HTMLReporter{Keep: false})
	if err := reporter.RunReporterModule(); err != nil {
		log.Printf("%v failed: %v\n", ModuleName, err)
		os.Exit(master.ReturnReporterFailed)
	}
	log.Printf("%v completed successfully\n", ModuleName)
}

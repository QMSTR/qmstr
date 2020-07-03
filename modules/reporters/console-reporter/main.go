package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/QMSTR/qmstr/lib/go-qmstr/cli"
	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

func main() {
	reporter := reporting.NewReporter(&ConsoleReporter{})
	go func() {
		<-cli.PingReporter // wait for the reporting phase to start
		log.Printf("Console reporter starts the reporting\n")
		if err := reporter.RunReporterModule(); err != nil {
			msg := fmt.Sprintf("%v failed: %v\n", reporter.GetModuleName(), err)
			log.Printf(msg)
			reporter.CtrlSvcClient.ShutdownModule(context.Background(), &service.ShutdownModuleRequest{
				Message: msg, DB: false})
			os.Exit(master.ReturnReporterFailed)
		}
		reporting.ReduceReportersCounter()
	}()
}

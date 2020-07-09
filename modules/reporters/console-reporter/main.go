package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/QMSTR/qmstr/lib/go-qmstr/cli"
	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

var wg sync.WaitGroup

func main() {
	reporter := reporting.NewReporter(&ConsoleReporter{})
	log.Printf("Console reporter was initialized")
	wg.Add(1)
	go func() {
		defer wg.Done()
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
	wg.Wait() // Waits until the goroutine is done
	log.Printf("Console reporter  finished")
}

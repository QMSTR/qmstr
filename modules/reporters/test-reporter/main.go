package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/lib/go-qmstr/cli"
	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/module"
	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/QMSTR/qmstr/lib/go-qmstr/tester"
)

type TestReporter struct{}

var testprojectNode *service.ProjectNode

func main() {
	reporter := reporting.NewReporter(&TestReporter{})
	go func() {
		<-cli.PingReporter // wait for the reporting phase to start
		log.Printf("Test reporter starts the reporting\n")
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

// Configure sets up the working directory for this reporting run.
func (testRporter *TestReporter) Configure(config map[string]string) error {
	return nil
}

// Report generates the actual reports.
func (testRporter *TestReporter) Report(masterClient *module.MasterClient) error {
	var err error
	testprojectNode, err = masterClient.RptSvcClient.GetProjectNode(context.Background(), &service.ProjectNode{})
	if err != nil {
		return fmt.Errorf("could not get project node: %v", err)
	}
	testSuite := []testing.InternalTest{
		{
			Name: "TestGraphLAdditionalInfo",
			F:    TestGraphLAdditionalInfo,
		},
	}
	t := &tester.DummyTestDeps{}
	testing.MainStart(t, testSuite, nil, nil).Run()

	return nil
}

// PostReport is called after the report has bee generated.
func (testRporter *TestReporter) PostReport() error {
	return nil
}

func TestGraphLAdditionalInfo(t *testing.T) {
	if len(testprojectNode.AdditionalInfo) == 0 {
		t.Logf("The graph doesn't contain any information nodes. Project name: %s", testprojectNode.Name)
		t.Fail()
	}
}

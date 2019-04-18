package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/lib/go-qmstr/master"
	"github.com/QMSTR/qmstr/lib/go-qmstr/reporting"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/QMSTR/qmstr/lib/go-qmstr/tester"
)

type TestReporter struct{}

var testprojectNode *service.ProjectNode

func main() {
	reporter := reporting.NewReporter(&TestReporter{})
	if err := reporter.RunReporterModule(); err != nil {
		log.Printf("%v failed: %v\n", reporter.GetModuleName(), err)
		os.Exit(master.ReturnReporterFailed)
	}
	log.Printf("%v completed successfully\n", reporter.GetModuleName())

}

// Configure sets up the working directory for this reporting run.
func (testRporter *TestReporter) Configure(config map[string]string) error {
	return nil
}

// Report generates the actual reports.
func (testRporter *TestReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient) error {
	var err error
	testprojectNode, err = rserv.GetProjectNode(context.Background(), &service.ProjectNode{})
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

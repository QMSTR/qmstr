package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/reporting"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/QMSTR/qmstr/pkg/tester"
)

type TestReporter struct{}

var testpackageNode *service.PackageNode

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
func (testRporter *TestReporter) Report(cserv service.ControlServiceClient, rserv service.ReportServiceClient, session string) error {
	var err error
	testpackageNode, err = cserv.GetPackageNode(context.Background(), &service.PackageRequest{Session: session})
	if err != nil {
		return fmt.Errorf("could not get package node: %v", err)
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
	if len(testpackageNode.AdditionalInfo) == 0 {
		t.Logf("The graph doesn't contain any information nodes. Package name: %s", testpackageNode.Name)
		t.Fail()
	}
}

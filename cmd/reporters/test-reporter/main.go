package main

import (
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/reporting"
	"github.com/QMSTR/qmstr/pkg/service"
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
func (testRporter *TestReporter) Report(packageNode *service.PackageNode) error {
	testpackageNode = packageNode
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

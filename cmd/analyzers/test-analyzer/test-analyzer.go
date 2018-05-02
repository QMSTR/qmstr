package main

import (
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/service"
)

var testnode *service.FileNode

type TestAnalyzer struct{}

func main() {
	analyzer := analysis.NewAnalyzer(&TestAnalyzer{})
	if err := analyzer.RunAnalyzerModule(); err != nil {
		log.Printf("%v failed: %v\n", analyzer.GetModuleName(), err)
		os.Exit(master.ReturnAnalyzerFailed)
	}
}

func (testanalyzer *TestAnalyzer) Configure(configMap map[string]string) error {
	return nil
}

func (testanalyzer *TestAnalyzer) Analyze(node *service.FileNode) (*service.InfoNodeSlice, error) {
	testnode = node
	testSuite := []testing.InternalTest{
		{
			Name: "TestGraphIntegrity",
			F:    TestGraphIntegrity,
		},
	}
	t := &DummyTestDeps{}
	testing.MainStart(t, testSuite, nil, nil).Run()
	return &service.InfoNodeSlice{}, nil
}

func TestGraphIntegrity(t *testing.T) {
	if testnode.Type == "linkedtarget" {
		if len(testnode.DerivedFrom) == 0 {
			t.Logf("Broken linked target %s. There are no derived nodes.", testnode.Name)
			t.Fail()
		}
		for _, node := range testnode.DerivedFrom {
			// TODO: better log messages
			if node.Type == "library" {
				for _, objectfile := range node.DerivedFrom {
					if objectfile.Type != "objectfile" || objectfile.DerivedFrom[0].Type != "sourcefile" {
						t.Logf("Broken library node %s.", node.Name)
						t.Fail()
					}
				}
			} else if node.Type == "objectfile" && node.DerivedFrom[0].Type != "sourcefile" {
				t.Logf("Broken object file %s .There is no source file connected to it.", node.Name)
				t.Fail()
			}
		}
	} else if testnode.Type == "library" {
		if len(testnode.DerivedFrom) == 0 {
			t.Logf("Broken library %s. There are no derived nodes.", testnode.Name)
			t.Fail()
		}
		for _, objectfile := range testnode.DerivedFrom {
			if objectfile.Type != "objectfile" || objectfile.DerivedFrom[0].Type != "sourcefile" {
				t.Logf("Broken library node %s.", testnode.Name)
				t.Fail()
			}
		}
	} else if testnode.Type == "objectfile" && testnode.DerivedFrom[0].Type != "sourcefile" {
		t.Logf("Broken object file %s .There is no source file connected to it.", testnode.Name)
		t.Fail()
	}
}

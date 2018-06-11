package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/QMSTR/qmstr/pkg/tester"
)

const (
	queryType = "sourcecode"
)

var (
	testnode     *service.FileNode
	pkgNode      *service.PackageNode
	tests        []string
	testfunction func(*testing.T)
)

type TestAnalyzer struct{}

func main() {
	fmt.Println("This is the testalyzer")
	analyzer := analysis.NewAnalyzer(&TestAnalyzer{})
	if err := analyzer.RunAnalyzerModule(); err != nil {
		log.Printf("%v failed: %v\n", analyzer.GetModuleName(), err)
		os.Exit(master.ReturnAnalyzerFailed)
	}
}

func (testanalyzer *TestAnalyzer) Configure(configMap map[string]string) error {
	if _, ok := configMap["tests"]; !ok {
		log.Println("No build graph tests provided. Running default test.")
		return nil
	}
	tests = strings.Split(configMap["tests"], ";")
	return nil
}

func (testanalyzer *TestAnalyzer) Analyze(controlService service.ControlServiceClient, session string) error {
	queryNode := &service.FileNode{Type: queryType}

	pkgNodeResp, err := controlService.GetPackageNode(context.Background(), &service.PackageRequest{Session: session})
	pkgNode = pkgNodeResp.PackageNode

	stream, err := controlService.GetFileNode(context.Background(), queryNode)
	if err != nil {
		log.Printf("Could not get file node %v", err)
		return err
	}

	// Run tests for file nodes
	for {
		fileNode, err := stream.Recv()
		if err == io.EOF {
			break
		}

		fmt.Println("Testalyzer running tests")
		testnode = fileNode
		testSuite := []testing.InternalTest{
			{
				Name: "TestGraphIntegrity",
				F:    TestGraphIntegrity,
			},
		}
		t := &tester.DummyTestDeps{}
		testing.MainStart(t, testSuite, nil, nil).Run()
	}

	// Run tests for package node
	for _, test := range tests {
		if test == "TestPackageNode" {
			testfunction = TestPackageNode
		} else if test == "TestCalcBuildGraph" {
			testfunction = TestCalcBuildGraph
		} else if test == "TestCurlBuildGraph" {
			testfunction = TestCurlBuildGraph
		} else {
			log.Printf("Unknown test. Please check the test name provided in the configuration.")
			os.Exit(master.ReturnAnalyzerFailed)
		}
		testSuite := []testing.InternalTest{
			{
				Name: test,
				F:    testfunction,
			},
		}
		t := &tester.DummyTestDeps{}
		testing.MainStart(t, testSuite, nil, nil).Run()
	}

	return nil
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

func TestPackageNode(t *testing.T) {
	if len(pkgNode.Targets) < 1 {
		t.Logf("Package node '%v' is not connected to any linked targets", pkgNode.Name)
		t.Fail()
	}
}

func TestCurlBuildGraph(t *testing.T) {
	if pkgNode.Targets[0].Path != "/buildroot/curl/build/src/curl" && pkgNode.Targets[1].Path != "/buildroot/curl/build/src/curl" {
		t.Logf("Package node %v is not connected to curl linked target", pkgNode.Name)
		t.Fail()
	}
	if pkgNode.Targets[0].Path != "/buildroot/curl/build/lib/libcurl.so" && pkgNode.Targets[1].Path != "/buildroot/curl/build/lib/libcurl.so" {
		t.Logf("Package node %v is not connected to libcurl linked target", pkgNode.Name)
		t.Fail()
	}
}

func TestCalcBuildGraph(t *testing.T) {
	if pkgNode.Targets[0].Path != "/buildroot/Calculator/calc" {
		t.Logf("Package node %v is not connected to curl linked target", pkgNode.Name)
		t.Fail()
	}
}

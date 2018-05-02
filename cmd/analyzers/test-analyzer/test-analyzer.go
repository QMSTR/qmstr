package main

import (
	"fmt"
	"testing"
)

func main() {
	testSuite := []testing.InternalTest{
		{
			Name: "TestCaseA",
			F:    TestCaseA,
		},
	}
	t := &DummyTestDeps{}
	testing.MainStart(t, testSuite, nil, nil).Run()
}

func TestCaseA(t *testing.T) {
	fmt.Println("Hello tester")
	t.Fail()
}

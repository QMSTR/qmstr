package main

import (
	"fmt"
	"testing"
)

func main() {
	testSuite := []testing.InternalTest{
		{
			Name: "TestGraphIntegrity",
			F:    TestGraphIntegrity,
		},
	}
	t := &DummyTestDeps{}
	testing.MainStart(t, testSuite, nil, nil).Run()
}

func TestGraphIntegrity(t *testing.T) {
	fmt.Println("Hello tester")
	t.Fail()
}

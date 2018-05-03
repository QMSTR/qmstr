package tester

import "io"

// DummyTestDeps is a dummy implementation of the testDeps interface used to simulate `go test` from qmstr module code
type DummyTestDeps struct{}

func (t *DummyTestDeps) MatchString(pat, str string) (bool, error) {
	return true, nil
}

func (t *DummyTestDeps) StartCPUProfile(io.Writer) error {
	return nil
}

func (t *DummyTestDeps) StopCPUProfile() {}

func (t *DummyTestDeps) WriteHeapProfile(io.Writer) error {
	return nil
}
func (t *DummyTestDeps) StartTestLog(io.Writer) {
}

func (t *DummyTestDeps) StopTestLog() error {
	return nil
}

func (t *DummyTestDeps) WriteProfileTo(string, io.Writer, int) error {
	return nil
}
func (t *DummyTestDeps) ImportPath() string {
	return ""
}

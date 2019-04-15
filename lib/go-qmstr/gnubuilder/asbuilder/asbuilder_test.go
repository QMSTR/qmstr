package asbuilder_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	builder "github.com/QMSTR/qmstr/lib/go-qmstr/gnubuilder/asbuilder"
)

func getTestCompiler() *builder.AsBuilder {
	return builder.NewAsBuilder("/tmp", log.New(os.Stdout, "TESTING", log.LstdFlags), false)
}

func TestOutputFlag(t *testing.T) {
	as := getTestCompiler()
	as.Analyze([]string{"as", "a.s", "-o", "a.o"})
	if as.Input != "a.s" && as.Output != "a.o" {
		t.Fail()
	}
}

func TestNoOutputFlag(t *testing.T) {
	as := getTestCompiler()
	as.Analyze([]string{"as", "a.s"})
	if as.Input != "a.s" && as.Output != "a.o" {
		t.Fail()
	}
}

func TestCleanCmdBoolArgs(t *testing.T) {
	as := getTestCompiler()
	as.Analyze([]string{"as", "--32", "a.s", "-o", "a.o"})
	if fmt.Sprintf("%v", as.Args) != "[a.s -o a.o]" {
		t.Fail()
	}
}

func TestCmdFlagArgs(t *testing.T) {
	as := getTestCompiler()
	as.Analyze([]string{"as", "-I", "..", "a.s", "-o", "a.o"})
	if len(as.Args) != 5 {
		t.Logf("Arguments: %v", as.Args)
		t.Fail()
	}
	if as.Input != "a.s" && as.Output != "a.o" {
		t.Logf("Input: %s and output: %s ", as.Input, as.Output)
		t.Fail()
	}
}

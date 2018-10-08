package objcopybuilder_test

import (
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/pkg/objcopybuilder"
)

func getTestBuilder() *objcopybuilder.ObjcopyBuilder {
	return objcopybuilder.NewObjcopyBuilder("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags), false)
}

func TestTargets(t *testing.T) {
	o := getTestBuilder()
	o.Analyze([]string{"objcopy", "-S", "a.so.dbg", "a.so"})
	if len(o.Targets) != 2 {
		t.Logf("%v", o.Targets)
		t.Fail()
	}
	if o.Targets[0] != "a.so.dbg" && o.Targets[1] != "a.so" {
		t.Fail()
	}
}

func TestSimpleArgs(t *testing.T) {
	o := getTestBuilder()
	o.Analyze([]string{"objcopy", "-S", "a", "a.bin"})
	if len(o.Args) != 2 {
		t.Logf("%v", o.Args)
		t.Fail()
	}
	if o.Input != "a" && o.Output != "a.bin" {
		t.Fail()
	}
}

func TestOutputTargetFlag(t *testing.T) {
	o := getTestBuilder()
	o.Analyze([]string{"objcopy", "-U", "-O", "binary", "a.elf"})
	if len(o.Args) != 3 {
		t.Logf("%v", o.Args)
		t.Fail()
	}
	if len(o.Targets) != 1 {
		t.Logf("%v", o.Targets)
		t.Fail()
	}
	if o.Input != "a.elf" && o.Output != "a.bin" {
		t.Fail()
	}

	o = getTestBuilder()
	o.Analyze([]string{"objcopy", "-U", "--output-target", "binary", "a.elf"})
	if len(o.Args) != 3 {
		t.Logf("%v", o.Args)
		t.Fail()
	}
	if len(o.Targets) != 1 {
		t.Logf("%v", o.Targets)
		t.Fail()
	}
	if o.Input != "a.elf" && o.Output != "a.bin" {
		t.Fail()
	}
}

func TestTargetTypeFlag(t *testing.T) {
	o := getTestBuilder()
	o.Analyze([]string{"objcopy", "-F", "binary", "a.bin"})
	if len(o.Args) != 3 {
		t.Logf("%v", o.Args)
		t.Fail()
	}
	if len(o.Targets) != 1 {
		t.Logf("%v", o.Targets)
		t.Fail()
	}
	if o.Input != "a.bin" && o.Output != "a.bin" {
		t.Fail()
	}
}

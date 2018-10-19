package ldbuilder_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/QMSTR/qmstr/pkg/gnubuilder"
	"github.com/QMSTR/qmstr/pkg/gnubuilder/ldbuilder"
	"github.com/spf13/afero"
)

func getTestBuilder() *ldbuilder.LdBuilder {
	builder := ldbuilder.NewLdBuilder("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags), false)
	builder.Afs = afero.NewMemMapFs()
	return builder
}

func fakeLibFile(builder *ldbuilder.LdBuilder, name string, static bool) error {
	libDir := gnubuilder.GetSysLibPath()[0]
	pre, dsuf, ssuf, err := gnubuilder.GetOsLibFixes()
	if err != nil {
		return err
	}
	var suf string
	if static {
		suf = ssuf[0]
	} else {
		suf = dsuf[0]
	}
	libpath := filepath.Join(libDir, fmt.Sprintf("%s%s%s", pre, name, suf))
	_, err = builder.Afs.Create(libpath)
	return err
}

func TestDefaultOutput(t *testing.T) {
	ld := getTestBuilder()
	ld.Analyze([]string{"ld", "a.o"})
	if ld.Output[0] != "a.out" {
		t.Fail()
	}
}

func TestSingleInput(t *testing.T) {
	ld := getTestBuilder()
	ld.Analyze([]string{"ld", "a.o"})
	if ld.Input[0] != "a.o" {
		t.Fail()
	}
}

func TestMultiInput(t *testing.T) {
	ld := getTestBuilder()
	ld.Analyze([]string{"ld", "a.o", "b.o"})
	if ld.Input[0] != "a.o" || ld.Input[1] != "b.o" {
		t.Fail()
	}
}

func TestDefinedOutput(t *testing.T) {
	ld := getTestBuilder()
	ld.Analyze([]string{"ld", "-o", "output", "a.o", "b.o"})
	if ld.Output[0] != "output" {
		t.Fail()
	}
	if ld.Input[0] != "a.o" || ld.Input[1] != "b.o" {
		t.Fail()
	}
}

func TestCleanCmdStringArgs(t *testing.T) {
	ld := getTestBuilder()
	ld.Analyze([]string{"ld", "-m", "elf_x86_64", "-o", "out.o", "b.a"})
	if fmt.Sprintf("%v", ld.Args) != "[-o out.o b.a]" {
		t.Fail()
	}
}

func TestCleanCmdBoolArgs(t *testing.T) {
	ld := getTestBuilder()
	ld.Analyze([]string{"ld", "--whole-archive", "a.a", "--no-whole-archive", "b.a", "c.a"})
	if fmt.Sprintf("%v", ld.Args) != "[a.a b.a c.a]" {
		t.Fail()
	}
}

func TestForcedStaticLib(t *testing.T) {
	ld := getTestBuilder()
	err := fakeLibFile(ld, "gcc", true)
	if err != nil {
		t.Fail()
	}
	ld.Analyze([]string{"ld", "-static-libgcc", "-o", "out", "a.c", "-lgcc"})
	if ld.Output[0] != "out" {
		t.Fail()
	}
	if _, ok := ld.StaticLibs["gcc"]; !ok {
		t.Fail()
	}
}

func TestNoOutput(t *testing.T) {
	ld := getTestBuilder()
	ld.Analyze([]string{"ld", "-m", "elf_x86_64", "-v"})
	if ld.Mode != gnubuilder.ModePrintOnly {
		t.Fail()
	}
}

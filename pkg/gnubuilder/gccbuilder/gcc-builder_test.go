package gccbuilder_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/QMSTR/qmstr/pkg/gnubuilder"

	"github.com/spf13/afero"

	builder "github.com/QMSTR/qmstr/pkg/gnubuilder/gccbuilder"
)

func getTestCompiler() *builder.GccBuilder {
	builder := builder.NewGccBuilder("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags), false)
	builder.Afs = afero.NewMemMapFs()
	return builder
}

func fakeLibFile(builder *builder.GccBuilder, name string, static bool) error {
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

func TestAssembleOnly(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "-c", "a.c"})
	if gcc.Mode != gnubuilder.ModeAssemble {
		t.Fail()
	}
}

func TestCompileOnly(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "-S", "a.c"})
	if gcc.Mode != gnubuilder.ModeCompile {
		t.Fail()
	}
}

func TestPreProcessorOnly(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "-E", "a.c"})
	if gcc.Mode != gnubuilder.ModePreproc {
		t.Fail()
	}
}

func TestSingleInput(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "a.c"})
	if gcc.Input[0] != "a.c" {
		t.Fail()
	}
}

func TestMultiInput(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "a.c", "b.c"})
	if gcc.Input[0] != "a.c" || gcc.Input[1] != "b.c" {
		t.Fail()
	}
}

func TestDefaultLinkOutput(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "a.c"})
	if gcc.Output[0] != "a.out" {
		t.Fail()
	}
}

func TestDefaultAssembleOutput(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "-c", "a.c"})
	if gcc.Output[0] != "a.o" {
		t.Fail()
	}
}

func TestDefaultMultiAssembleOutput(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "-c", "a.c", "b.c"})
	if gcc.Output[0] != "a.o" || gcc.Output[1] != "b.o" {
		t.Fail()
	}
}

func TestDefinedOutput(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "-o", "outProg", "a.c", "b.c"})
	if gcc.Output[0] != "outProg" {
		t.Fail()
	}
	if gcc.Input[0] != "a.c" || gcc.Input[1] != "b.c" {
		t.Fail()
	}
}

func TestCleanCommandlineStringArgs(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "-DMACRO", "a.c", "b.c"})
	if fmt.Sprintf("%v", gcc.Args) != "[a.c b.c]" {
		t.Fail()
	}
	gcc.Analyze([]string{"gcc", "-D", "MACRO", "a.c", "b.c"})
	if fmt.Sprintf("%v", gcc.Args) != "[a.c b.c]" {
		t.Fail()
	}
	gcc.Analyze([]string{"gcc", "-D", "MACRO=test", "a.c", "b.c"})
	if fmt.Sprintf("%v", gcc.Args) != "[a.c b.c]" {
		t.Fail()
	}
}

func TestLastArgFlag(t *testing.T) {
	gcc := getTestCompiler()
	gcc.Analyze([]string{"gcc", "-DMACRO", "a.c", "b.c", "-DMACRO"})
	if fmt.Sprintf("%v", gcc.Args) != "[a.c b.c]" {
		t.Fail()
	}
	gcc.Analyze([]string{"gcc", "-D", "MACRO", "a.c", "b.c", "-D", "MACRO"})
	if fmt.Sprintf("%v", gcc.Args) != "[a.c b.c]" {
		t.Fail()
	}
	gcc.Analyze([]string{"gcc", "-D", "MACRO=test", "a.c", "b.c", "-D", "MACRO=test"})
	if fmt.Sprintf("%v", gcc.Args) != "[a.c b.c]" {
		t.Fail()
	}
}

func TestForcedStaticLib(t *testing.T) {
	gcc := getTestCompiler()
	err := fakeLibFile(gcc, "gcc", true)
	if err != nil {
		t.Fail()
	}
	gcc.Analyze([]string{"gcc", "-static-libgcc", "-o", "out", "a.c", "-lgcc"})
	if gcc.Output[0] != "out" {
		t.Fail()
	}
	if _, ok := gcc.StaticLibs["gcc"]; !ok {
		t.Fail()
	}
}

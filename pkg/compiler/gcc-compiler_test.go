package compiler_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/pkg/compiler"
)

func TestAssembleOnly(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "-c", "a.c"})
	if gcc.Mode != compiler.Assemble {
		t.Fail()
	}
}

func TestCompileOnly(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "-S", "a.c"})
	if gcc.Mode != compiler.Compile {
		t.Fail()
	}
}

func TestPreProcessorOnly(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "-E", "a.c"})
	if gcc.Mode != compiler.Preproc {
		t.Fail()
	}
}

func TestSingleInput(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "a.c"})
	if gcc.Input[0] != "a.c" {
		t.Fail()
	}
}

func TestMultiInput(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "a.c", "b.c"})
	if gcc.Input[0] != "a.c" || gcc.Input[1] != "b.c" {
		t.Fail()
	}
}

func TestDefaultLinkOutput(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "a.c"})
	if gcc.Output[0] != "a.out" {
		t.Fail()
	}
}

func TestDefaultAssembleOutput(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "-c", "a.c"})
	if gcc.Output[0] != "a.o" {
		t.Fail()
	}
}

func TestDefaultMultiAssembleOutput(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "-c", "a.c", "b.c"})
	if gcc.Output[0] != "a.o" || gcc.Output[1] != "b.o" {
		t.Fail()
	}
}

func TestDefinedOutput(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
	gcc.Analyze([]string{"gcc", "-o", "outProg", "a.c", "b.c"})
	if gcc.Output[0] != "outProg" {
		t.Fail()
	}
	if gcc.Input[0] != "a.c" || gcc.Input[1] != "b.c" {
		t.Fail()
	}
}

func TestCleanCommandlineStringArgs(t *testing.T) {
	gcc := compiler.NewGccCompiler("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags))
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

package compiler_test

import "testing"
import "github.com/QMSTR/qmstr/pkg/compiler"

func TestAssembleOnly(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "-c", "a.c"})
	if gcc.Mode != compiler.Assemble {
		t.Fail()
	}
}

func TestCompileOnly(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "-S", "a.c"})
	if gcc.Mode != compiler.Compile {
		t.Fail()
	}
}

func TestPreProcessorOnly(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "-E", "a.c"})
	if gcc.Mode != compiler.Preproc {
		t.Fail()
	}
}

func TestSingleInput(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "a.c"})
	if gcc.Input[0] != "a.c" {
		t.Fail()
	}
}

func TestMultiInput(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "a.c", "b.c"})
	if gcc.Input[0] != "a.c" || gcc.Input[1] != "b.c" {
		t.Fail()
	}
}

func TestDefaultLinkOutput(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "a.c"})
	if gcc.Output[0] != "a.out" {
		t.Fail()
	}
}

func TestDefaultAssembleOutput(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "-c", "a.c"})
	if gcc.Output[0] != "a.o" {
		t.Fail()
	}
}

func TestDefaultMultiAssembleOutput(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "-c", "a.c", "b.c"})
	if gcc.Output[0] != "a.o" || gcc.Output[1] != "b.o" {
		t.Fail()
	}
}

func TestDefinedOutput(t *testing.T) {
	gcc := compiler.GccCompiler{}
	gcc.Analyze([]string{"gcc", "-o", "outProg", "a.c", "b.c"})
	if gcc.Output[0] != "outProg" {
		t.Fail()
	}
	if gcc.Input[0] != "a.c" || gcc.Input[1] != "b.c" {
		t.Fail()
	}
}

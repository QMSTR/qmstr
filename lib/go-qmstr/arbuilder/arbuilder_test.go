package arbuilder_test

import (
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/lib/go-qmstr/arbuilder"
)

func getTestBuilder() *arbuilder.ArBuilder {
	return arbuilder.NewArBuilder("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags), false)
}

func TestModeReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "r", "libtest.a", "test.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
}

func TestQuickAppendReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "q", "libtest.a", "test.o"})
	if ar.Command != arbuilder.QuickAppend {
		t.Fail()
	}
}

func TestNoopFlagsReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "--help"})
	if ar.Command != arbuilder.Undef {
		t.Fail()
	}
	ar.Analyze([]string{"ar", "--version"})
	if ar.Command != arbuilder.Undef {
		t.Fail()
	}
}

func TestTargetFlagReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "r", "--target=targetformat", "lib.a", "obj.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 2 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" {
		t.Fail()
	}

	ar = getTestBuilder()
	ar.Analyze([]string{"ar", "r", "--target", "targetformat", "lib.a", "obj.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 2 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" {
		t.Fail()
	}
}

func TestPluginFlagReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "r", "--plugin=plugin", "lib.a", "obj.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 2 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" {
		t.Fail()
	}

	ar = getTestBuilder()
	ar.Analyze([]string{"ar", "r", "--plugin", "plugin", "lib.a", "obj.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 2 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" {
		t.Fail()
	}
}

func TestEmuFlagReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "r", "-X32_64", "lib.a", "obj.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 2 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" {
		t.Fail()
	}
}

func TestBeforeModifierReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "rb", "relpos", "lib.a", "obj.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 2 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" {
		t.Fail()
	}
}

func TestAfterModifierReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "ra", "relpos", "lib.a", "obj.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 2 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" {
		t.Fail()
	}
}

func TestCountModifierReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "rN", "4", "lib.a", "obj.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 2 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" {
		t.Fail()
	}
}

func TestAllReplace(t *testing.T) {
	ar := getTestBuilder()
	ar.Analyze([]string{"ar", "-X32_64", "rN", "4", "--plugin", "plugin", "--target=target", "lib.a", "obj.o", "obj2.o"})
	if ar.Command != arbuilder.Replace {
		t.Fail()
	}
	if len(ar.CommandLineArgs) != 3 {
		t.Logf("%v", ar.CommandLineArgs)
		t.Fail()
	}
	if ar.CommandLineArgs[0] != "lib.a" && ar.CommandLineArgs[1] != "obj.o" && ar.CommandLineArgs[2] != "obj2.o" {
		t.Fail()
	}
}

package stripbuilder_test

import (
	"log"
	"os"
	"testing"

	"github.com/QMSTR/qmstr/lib/go-qmstr/stripbuilder"
)

func getTestBuilder() *stripbuilder.StripBuilder {
	return stripbuilder.NewStripBuilder("/tmp", log.New(os.Stdout, "TESTING ", log.LstdFlags), false)
}

func TestArgs(t *testing.T) {
	s := getTestBuilder()
	s.Analyze([]string{"strip", "-S", "a.out"})
	if len(s.Input) != 1 && len(s.Output) != 1 {
		t.Fail()
	}
	if s.Input[0] != "a.out" && s.Output[0] != "a.out" {
		t.Fail()
	}
}

func TestMultiInputs(t *testing.T) {
	s := getTestBuilder()
	s.Analyze([]string{"strip", "--strip-debug", "a.out", "foo"})
	if len(s.Input) != 2 && len(s.Output) != 2 {
		t.Fail()
	}
	if s.Input[0] != "a.out" && s.Input[1] != "foo" && s.Output[0] != "a.out" && s.Output[1] != "foo" {
		t.Fail()
	}
}

func TestStringFlag(t *testing.T) {
	s := getTestBuilder()
	s.Analyze([]string{"strip", "--remove-section=.comment", "foo"})
	if len(s.Input) != 1 && len(s.Output) != 1 {
		t.Logf("%v", s.Input)
		t.Fail()
	}
	if s.Input[0] != "foo" && s.Output[0] != "foo" {
		t.Logf("%v", s.Input)
		t.Fail()
	}
	s = getTestBuilder()
	s.Analyze([]string{"strip", "--remove-section", ".comment", "foo"})
	if len(s.Input) != 1 && len(s.Output) != 1 {
		t.Logf("%v", s.Input)
		t.Fail()
	}
	if s.Input[0] != "foo" && s.Output[0] != "foo" {
		t.Logf("%v", s.Input)
		t.Fail()
	}
}

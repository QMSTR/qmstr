package main

import "testing"

func TestCmdFix(t *testing.T) {
	fixedCmdLine := fixCmdLine(&[]string{"qmstrctl", "create", "file:hash:12345"})
	if (*fixedCmdLine)[2] != "file" {
		t.Fail()
	}
}

func TestCmdFixWithArgs(t *testing.T) {
	fixedCmdLine := fixCmdLine(&[]string{"qmstrctl", "create", "file:hash:12345", "--Name", "foobar"})
	if (*fixedCmdLine)[2] != "file" {
		t.Fail()
	}
}

func TestCmdFixNothingToFix(t *testing.T) {
	originalCmdLine := []string{"qmstrctl", "create", "file", "file:hash:12345"}
	fixedCmdLine := fixCmdLine(&originalCmdLine)
	if fixedCmdLine != &originalCmdLine {
		t.Fail()
	}
}

func TestCmdFixNothingToFixOther(t *testing.T) {
	originalCmdLine := []string{"qmstrctl", "start", "--wait"}
	fixedCmdLine := fixCmdLine(&originalCmdLine)
	if fixedCmdLine != &originalCmdLine {
		t.Fail()
	}
}

func TestCmdFixNothingToFixInvalid(t *testing.T) {
	originalCmdLine := []string{"qmstrctl", "create"}
	fixedCmdLine := fixCmdLine(&originalCmdLine)
	if fixedCmdLine != &originalCmdLine {
		t.Fail()
	}
}

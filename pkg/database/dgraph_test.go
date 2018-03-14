package database

import "testing"

func TestVarNameCalculation(t *testing.T) {
	if getVarName(0) != "A" {
		t.Fail()
	}
	if getVarName(25) != "Z" {
		t.Fail()
	}
	if getVarName(26) != "AA" {
		t.Fail()
	}
	if getVarName(52) != "AAA" {
		t.Fail()
	}
}

package database

import (
	"testing"

	"github.com/QMSTR/qmstr/pkg/service"
)

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

func TestFillType(t *testing.T) {
	pkgNode := service.PackageNode{Targets: []*service.FileNode{&service.FileNode{Name: "Test"}}}
	fillTypeField(&pkgNode)
	if pkgNode.PackageNodeType != "_" {
		t.Fail()
	}
	if pkgNode.Targets[0].FileNodeType != "_" {
		t.Fail()
	}

}

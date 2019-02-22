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

func TestCheckSchema(t *testing.T) {
	const schema = `data:string @index(hash) .
hash:string @index(exact) .
name:string @index(hash) .
path:string @index(hash,trigram) .
type:string @index(hash) .
fileType:int @index(int) .
phase:int .
dataNodes:uid @reverse .
buildConfig:string .
projectNodeType:string @index(hash) .
dataNodeType:string @index(hash) .
fileNodeType:string @index(hash) .
infoNodeType:string @index(hash) .
diagnosticNodeType:string @index(hash) .
severity:string @index(hash) .
additionalInfo:uid .
packageNodeType:string @index(hash) .
analyzerNodeType:string @index(hash) .
qmstrStateNodeType:string .
	`
	if !CheckSchema(schema) {
		t.Fail()
	}
}

package cli

import (
	"log"
	"testing"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/pflag"
)

func TestStringFlagNodeCreation(t *testing.T) {
	cmdFlags = &pflag.FlagSet{}
	err := generateFlags(&service.FileNode{}, cmdFlags)
	if err != nil {
		log.Printf("%v", err)
	}
	cmdFlags.Parse([]string{"--name", "foobar", "--path", "/dev/null"})
	if err = createNode("file:hash:12345", false); err != nil {
		t.FailNow()
	}
	if currentNode.(*service.FileNode).Hash != "12345" || currentNode.(*service.FileNode).Name != "foobar" || currentNode.(*service.FileNode).Path != "/dev/null" {
		t.Fail()
	}
}

func TestBoolFlagNodeCreation(t *testing.T) {
	cmdFlags = &pflag.FlagSet{}
	err := generateFlags(&service.FileNode{}, cmdFlags)
	if err != nil {
		log.Printf("%v", err)
	}
	cmdFlags.Parse([]string{"--broken"})
	if err = createNode("file:hash:12345", false); err != nil {
		t.FailNow()
	}
	if currentNode.(*service.FileNode).Hash != "12345" || !currentNode.(*service.FileNode).Broken {
		t.Fail()
	}
}

func TestIntFlagNodeCreation(t *testing.T) {
	cmdFlags = &pflag.FlagSet{}
	err := generateFlags(&service.FileNode{}, cmdFlags)
	if err != nil {
		log.Printf("%v", err)
	}
	cmdFlags.Parse([]string{"--filetype", "2"})
	if err = createNode("file:hash:12345", false); err != nil {
		t.FailNow()
	}
	if currentNode.(*service.FileNode).Hash != "12345" || currentNode.(*service.FileNode).FileType != service.FileNode_INTERMEDIATE {
		t.Fail()
	}
}

func TestStringPackageNodeCreation(t *testing.T) {
	cmdFlags = &pflag.FlagSet{}
	err := generateFlags(&service.PackageNode{}, cmdFlags)
	if err != nil {
		log.Printf("%v", err)
	}
	cmdFlags.Parse([]string{"--buildconfig", "foobar"})
	if err = createNode("package:TestPackage", false); err != nil {
		t.FailNow()
	}
	if currentNode.(*service.PackageNode).Name != "TestPackage" || currentNode.(*service.PackageNode).BuildConfig != "foobar" {
		t.Fail()
	}
}

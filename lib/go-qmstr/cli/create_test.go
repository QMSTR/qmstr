package cli

import (
	"log"
	"testing"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
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
	if currentNode.(*service.FileNode).FileData.GetHash() != "12345" || currentNode.(*service.FileNode).Name != "foobar" || currentNode.(*service.FileNode).Path != "/dev/null" {
		t.Fail()
	}
}

func TestBoolFlagNodeCreation(t *testing.T) {
	cmdFlags = &pflag.FlagSet{}
	err := generateFlags(&service.FileNode{}, cmdFlags)
	if err != nil {
		log.Printf("%v", err)
	}
	if err = createNode("file:hash:12345", false); err != nil {
		t.FailNow()
	}
	if currentNode.(*service.FileNode).FileData.GetHash() != "12345" {
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
	if currentNode.(*service.FileNode).FileData.GetHash() != "12345" {
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

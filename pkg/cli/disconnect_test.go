package cli

import (
	"testing"

	"github.com/QMSTR/qmstr/pkg/service"
)

func TestRemoveFileNode(t *testing.T) {
	list := []*service.FileNode{
		&service.FileNode{Uid: "aaaa"},
		&service.FileNode{Uid: "bbbb"},
		&service.FileNode{Uid: "cccc"},
		&service.FileNode{Uid: "dddd"},
	}
	node := &service.FileNode{Uid: "cccc"}
	out, err := removeFileNodeFromList(list, node)
	if err != nil {
		t.Fatal(err)
	}
	if out[2].Uid != "dddd" {
		t.Fatalf("expecting out[2] == dddd")
	}
	node = &service.FileNode{Uid: "dddd"}
	out, err = removeFileNodeFromList(list, node)
	if err != nil {
		t.Fatal(err)
	}
	if out[2].Uid != "cccc" {
		t.Fatalf("expecting out[2] == cccc")
	}
}

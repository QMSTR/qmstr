package cli

import (
	"reflect"
	"strconv"
	"testing"

	"github.com/QMSTR/qmstr/pkg/service"
)

func TestInvalidIdentifier(t *testing.T) {
	_, err := ParseNodeID("")
	if err != ErrEmptyNodeIdent {
		t.Fail()
	}
}

func TestInvalidNodeType(t *testing.T) {
	_, err := ParseNodeID("foo:invalid")
	if err.Error() != "Unsupported node type foo" {
		t.Fail()
	}
}

type Foo struct {
	Bar     string
	BarInt  int64
	BarBool bool
	BarFoo  *Foo
}

func TestSetFieldValue(t *testing.T) {
	foo := Foo{}
	err := setFieldValue(&foo, "Bar", "foo")
	if err != nil {
		t.Fail()
	}
	if foo.Bar != "foo" {
		t.Fail()
	}
}

func TestInvalidAttribute(t *testing.T) {
	err := setFieldValue(&Foo{}, "baz", "foo")
	if err != ErrInvalidAttribute {
		t.Fail()
	}
}

func TestInt64Attribute(t *testing.T) {
	foo := Foo{}
	err := setFieldValue(&foo, "BarInt", "5")
	if err != nil {
		t.Fail()
	}
	err = setFieldValue(&foo, "BarInt", "test")
	if err == nil || reflect.TypeOf(err) != reflect.TypeOf((*strconv.NumError)(nil)) {
		t.Fail()
	}
	if foo.BarInt != 5 {
		t.Fail()
	}
}

func TestBoolAttribute(t *testing.T) {
	foo := Foo{BarBool: false}
	err := setFieldValue(&foo, "BarBool", "true")
	if err != nil {
		t.Fail()
	}
	err = setFieldValue(&foo, "BarBool", "test")
	if err == nil || reflect.TypeOf(err) != reflect.TypeOf((*strconv.NumError)(nil)) {
		t.Fail()
	}
	if !foo.BarBool {
		t.Fail()
	}
}

func TestUnsupportedAttribute(t *testing.T) {
	foo := Foo{}
	err := setFieldValue(&foo, "BarFoo", "barfoo")
	if err == nil || err.Error() != "Unsupported type ptr" {
		t.Fail()
	}
}

func TestNotStruct(t *testing.T) {
	foo := 5
	err := setFieldValue(&foo, "BarFoo", "barfoo")
	if err == nil || err.Error() != "Not a struct: int" {
		t.Fail()
	}
}

func TestCallByValue(t *testing.T) {
	err := setFieldValue(Foo{}, "BarFoo", "barfoo")
	if err == nil || err != ErrCallByValue {
		t.Fail()
	}
}

func TestFileNodeParsing(t *testing.T) {
	fileNode, err := ParseNodeID("file:/dev/null")
	if err != nil {
		t.FailNow()
	}
	if fileNode.(*service.FileNode).Path != "/dev/null" {
		t.Fail()
	}
	fileNode, err = ParseNodeID("file:hash:deadbeef")
	if err != nil {
		t.FailNow()
	}
	if fileNode.(*service.FileNode).Hash != "deadbeef" {
		t.Fail()
	}
	fileNode, err = ParseNodeID("file")
	if err != nil {
		t.FailNow()
	}
	if fileNode.(*service.FileNode).Hash != "" {
		t.Fail()
	}
}

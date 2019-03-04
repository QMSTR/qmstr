package cli

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/QMSTR/qmstr/pkg/service"
)

var ErrEmptyNodeIdent = errors.New("Empty node identifier")
var ErrInvalidAttribute = errors.New("Invalid attribute")
var ErrCallByValue = errors.New("you shall not call setFieldValue by value")

func TokenizeNodeID(nodeid string) (string, []string, error) {
	if nodeid == "" {
		return "", nil, ErrEmptyNodeIdent
	}
	nodeIDTokens := strings.Split(nodeid, ":")
	nodetype := nodeIDTokens[0]
	tokens := []string{}
	if len(nodeIDTokens) > 1 {
		tokens = nodeIDTokens[1:]
	}
	return nodetype, tokens, nil
}

func ParseNodeID(nodeid string) (interface{}, error) {
	nodeType, nodeIDTokens, err := TokenizeNodeID(nodeid)
	if err != nil {
		return nil, err
	}
	switch nodeType {
	case "file":
		return createResult(&service.FileNode{}, "Path", nodeIDTokens)
	case "package":
		return createResult(&service.PackageNode{}, "Name", nodeIDTokens)
	case "project":
		return createResult(&service.ProjectNode{}, "Name", nodeIDTokens)
	case "info":
		return nil, fmt.Errorf("%s not yet supported", nodeType)
	case "data":
		return nil, fmt.Errorf("%s not yet supported", nodeType)
	default:
		return nil, fmt.Errorf("Unsupported node type %s", nodeType)
	}
}

func createResult(node interface{}, defaultAttribute string, args []string) (interface{}, error) {
	var attr string
	var value string

	// empty node
	if len(args) < 1 {
		return node, nil
	}
	// set default attribute
	if len(args) < 2 {
		attr = defaultAttribute
		value = args[0]
	} else {
		attr = strings.Title(args[0])
		value = args[1]
	}

	err := setFieldValue(node, attr, value)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func setFieldValue(nodeStruct interface{}, attribute string, value string) error {
	val := reflect.ValueOf(nodeStruct)
	if val.Kind() != reflect.Ptr {
		return ErrCallByValue
	}
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	if kind := val.Kind(); kind != reflect.Struct {
		return fmt.Errorf("Not a struct: %v", kind)
	}

	field := val.FieldByName(attribute)
	if !field.IsValid() {
		return ErrInvalidAttribute
	}

	switch field.Kind() {
	case reflect.String:
		// no need to test if string is a string
		field.SetString(value)
		return nil
	case reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		field.SetInt(i)
		return nil
	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		field.SetBool(b)
		return nil
	default:
		return fmt.Errorf("Unsupported type %v", field.Kind())
	}
}

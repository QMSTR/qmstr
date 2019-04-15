package cli

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/spf13/pflag"
)

func generateFlags(structure interface{}, targetFlagSet *pflag.FlagSet) error {
	val := reflect.ValueOf(structure)
	// need ptr to check for canSet()
	if val.Kind() != reflect.Ptr {
		return ErrCallByValue
	}
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	if kind := val.Kind(); kind != reflect.Struct {
		return fmt.Errorf("Not a struct: %v", kind)
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		fieldName := val.Type().Field(i).Name
		structName := val.Type().Name()

		if strings.HasSuffix(fieldName, "NodeType") || strings.HasPrefix(fieldName, "XXX_") || fieldName == "Uid" {
			continue
		}

		// flags should be lower case
		fieldName = strings.ToLower(fieldName)

		switch field.Kind() {
		case reflect.String:
			targetFlagSet.String(fieldName, "", fmt.Sprintf("Set %s's %s", structName, fieldName))
		case reflect.Int64:
			targetFlagSet.Int64(fieldName, 0, fmt.Sprintf("Set %s's %s", structName, fieldName))
		case reflect.Int32:
			targetFlagSet.Int32(fieldName, 0, fmt.Sprintf("Set %s's %s", structName, fieldName))
		case reflect.Bool:
			targetFlagSet.Bool(fieldName, false, fmt.Sprintf("Set %s's %s", structName, fieldName))
		}
	}
	return nil
}

func setField(nodeStruct interface{}, attribute string, value interface{}) error {
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

	// find right field
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if !field.IsValid() || !field.CanSet() {
			continue
		}
		if strings.ToLower(val.Type().Field(i).Name) == attribute {
			attribute = val.Type().Field(i).Name
		}
	}

	field := val.FieldByName(attribute)
	if !field.IsValid() || !field.CanSet() {
		return ErrInvalidAttribute
	}

	switch field.Kind() {
	case reflect.String:
		field.SetString(value.(string))
	case reflect.Bool:
		field.SetBool(value.(bool))
	case reflect.Int32:
		field.SetInt(int64(value.(int32)))
	case reflect.Int64:
		field.SetInt(value.(int64))
	}
	return nil
}

func visitNodeFlag(flag *pflag.Flag) {
	fieldName := flag.Name
	var value interface{}
	var err error

	switch flag.Value.Type() {
	case "bool":
		value, err = cmdFlags.GetBool(flag.Name)
		if err != nil {
			log.Fatalf("Failed to evaluate node flags: %v", err)
		}
	case "string":
		value, err = cmdFlags.GetString(flag.Name)
		if err != nil {
			log.Fatalf("Failed to evaluate node flags: %v", err)
		}
	case "int64":
		value, err = cmdFlags.GetInt64(flag.Name)
		if err != nil {
			log.Fatalf("Failed to evaluate node flags: %v", err)
		}
	case "int32":
		value, err = cmdFlags.GetInt32(flag.Name)
		if err != nil {
			log.Fatalf("Failed to evaluate node flags: %v", err)
		}
	}
	setField(currentNode, fieldName, value)
}

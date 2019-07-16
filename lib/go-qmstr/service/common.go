package service

import (
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"strings"
)

func NewFileNode(path string, hash string) FileNode {
	node := FileNode{
		Path: path,
		Hash: hash,
		Name: filepath.Base(path),
	}
	return node
}

func CreateInfoNode(infoType string, dataNodes ...*InfoNode_DataNode) *InfoNode {
	return &InfoNode{
		Type:      infoType,
		DataNodes: dataNodes,
	}
}

func CreateWarningNode(warning string) *InfoNode {
	return CreateInfoNode("warning", &InfoNode_DataNode{
		Type: "warning_message",
		Data: warning,
	})
}

func CreateErrorNode(errorMes string) *InfoNode {
	return CreateInfoNode("error", &InfoNode_DataNode{
		Type: "error_message",
		Data: errorMes,
	})
}

func (in *InfoNode) Describe(indent string) string {
	describe := []string{fmt.Sprintf("%s|- Type: %s, Confidence score: %v", indent, in.Type, in.ConfidenceScore)}
	indent = indent + "\t"
	for _, anode := range in.Analyzer {
		describe = append(describe, fmt.Sprintf("%s|- Name: %s", indent, anode.Name))
	}
	for _, dnode := range in.DataNodes {
		describe = append(describe, fmt.Sprintf("%s|- Type: %s, Data: %s", indent, dnode.Type, dnode.Data))
	}
	return strings.Join(describe, "\n")
}

func (fn *FileNode) Describe(less bool, indent string) string {
	describe := []string{fmt.Sprintf("%s|- Name: %s, Path: %s, Hash: %s, Timestamp: %v", indent, fn.Name, fn.Path, fn.Hash, fn.Timestamp)}
	indent = indent + "\t"
	if !less {
		for _, inode := range fn.AdditionalInfo {
			describe = append(describe, inode.Describe(indent))
		}
	}
	for _, fnode := range fn.DerivedFrom {
		describe = append(describe, fnode.Describe(less, indent))
	}
	for _, dep := range fn.Dependencies {
		describe = append(describe, dep.Describe(less, indent))
	}
	return strings.Join(describe, "\n")
}

func (pkg *PackageNode) Describe(less bool, indent string) string {
	describe := []string{fmt.Sprintf("|- Name: %s", pkg.Name)}
	if !less {
		for _, inode := range pkg.AdditionalInfo {
			describe = append(describe, inode.Describe("\t"))
		}
	}
	for _, tnode := range pkg.Targets {
		describe = append(describe, tnode.Describe(less, "\t"))
	}
	return strings.Join(describe, "\n")
}

func (pn *ProjectNode) Describe(less bool) string {
	describe := []string{fmt.Sprintf("|- Name: %s", pn.Name)}
	if !less {
		for _, inode := range pn.AdditionalInfo {
			describe = append(describe, inode.Describe("\t"))
		}
	}
	for _, pkgnode := range pn.Packages {
		describe = append(describe, pkgnode.Describe(less, "\t"))
	}
	return strings.Join(describe, "\n")
}

func (fn *FileNode) IsValid() bool {
	return fn.Path != ""
}

func (pn *PackageNode) IsValid() bool {
	return pn.Name != "" && pn.Version != ""
}

func checkEmpty(structure interface{}) error {
	val := reflect.ValueOf(structure)
	for val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return errors.New("provided non-struct")
	}

	if !reflect.DeepEqual(val.Interface(), reflect.Zero(val.Type()).Interface()) {
		return errors.New("non empty struct")
	}

	return nil
}

func (pn *PackageNode) IsEmpty() bool {
	err := checkEmpty(pn)
	return err == nil
}

func (prn *ProjectNode) IsValid() bool {
	return prn.Name != ""
}

func (prn *ProjectNode) GetMetaData(key string, defaultValue string) string {
	value, err := getMetaData(key, prn.GetAdditionalInfo())
	if err != nil {
		return defaultValue
	}
	return value
}

func (pn *PackageNode) GetMetaData(key string, defaultValue string) string {
	value, err := getMetaData(key, pn.GetAdditionalInfo())
	if err != nil {
		return defaultValue
	}
	return value
}

func getMetaData(key string, info []*InfoNode) (string, error) {
	for _, inode := range info {
		if inode.Type == "metadata" {
			for _, dnode := range inode.DataNodes {
				if key == dnode.Type {
					return dnode.Data, nil
				}
			}
		}
	}
	return "", fmt.Errorf("No metadata found for key %s", key)
}

func RemoveSlash(value string) string {
	newvalue := strings.Replace(value, "/", "_", -1)
	return newvalue
}

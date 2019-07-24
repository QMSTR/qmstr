package service

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"reflect"
	"strings"
)

func NewFileNode(path string, hash string) FileNode {
	p := []*PathInfo{&PathInfo{Phase: PathInfo_Build, Path: path}}
	node := FileNode{
		Paths: p,
		Hash:  hash,
		Name:  filepath.Base(path),
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
	path, err := GetFilePath(fn)
	if err != nil {
		log.Fatal(err)
	}
	describe := []string{fmt.Sprintf("%s|- Name: %s, Path: %s, Hash: %s, Timestamp: %v", indent, fn.Name, path, fn.Hash, fn.Timestamp)}
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
	if len(fn.Paths) == 0 {
		return false
	}
	for _, pathInfo := range fn.Paths {
		if pathInfo.Path == "" {
			return false
		}
	}
	return true
}

func (pn *PackageNode) IsValid() bool {
	return pn.Name != "" && pn.Version != ""
}

func (pi *PathInfo) IsValid() bool {
	return pi.Path != ""
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

//GetFilePath returns the last path added to the node
func GetFilePath(node *FileNode) (string, error) {
	if !node.IsValid() {
		return "", fmt.Errorf("Node %s, %s, %s does not have a path", node.Uid, node.Name, node.Hash)
	}
	i := len(node.Paths)
	return node.Paths[i-1].Path, nil
}

//GetPathInfo returns the last pathInfo added to the node
func GetPathInfo(node *FileNode) (*PathInfo, error) {
	if !node.IsValid() {
		return nil, fmt.Errorf("Node %v does not have a path", node)
	}
	i := len(node.Paths)
	return node.Paths[i-1], nil
}

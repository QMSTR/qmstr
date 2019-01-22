package service

import (
	"fmt"
	"path/filepath"
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
	describe := []string{fmt.Sprintf("%s|- Name: %s, Path: %s, Hash: %s", indent, fn.Name, fn.Path, fn.Hash)}
	indent = indent + "\t"
	if !less {
		for _, inode := range fn.AdditionalInfo {
			describe = append(describe, inode.Describe(indent))
		}
	}
	for _, fnode := range fn.DerivedFrom {
		describe = append(describe, fnode.Describe(less, indent))
	}
	return strings.Join(describe, "\n")
}

func (pn *PackageNode) Describe(less bool) string {
	describe := []string{fmt.Sprintf("|- Name: %s", pn.Name)}
	if !less {
		for _, inode := range pn.AdditionalInfo {
			describe = append(describe, inode.Describe("\t"))
		}
	}
	for _, fnode := range pn.Targets {
		describe = append(describe, fnode.Describe(less, "\t"))
	}
	return strings.Join(describe, "\n")
}

package service

import "path/filepath"

const (
	_ = iota
	NodeTypeFileNode
	NodeTypeInfoNode
	NodeTypeDataNode
	NodeTypeAnalyzerNode
	NodeTypePackageNode
)

func NewFileNode(path string, hash string) FileNode {
	node := FileNode{
		Path: path,
		Hash: hash,
		Name: filepath.Base(path),
	}
	return node
}

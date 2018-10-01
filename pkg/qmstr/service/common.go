package service

import (
	"path/filepath"
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

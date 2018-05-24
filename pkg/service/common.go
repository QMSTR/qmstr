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

func SanitizeFileNode(filenode *FileNode) {
	filenode.NodeType = NodeTypeFileNode
	for _, parent := range filenode.DerivedFrom {
		SanitizeFileNode(parent)
	}
	for _, info := range filenode.AdditionalInfo {
		SanitizeInfoNode(info)
	}
}

func SanitizeInfoNode(infonode *InfoNode) {
	infonode.NodeType = NodeTypeInfoNode
	for _, dataNode := range infonode.DataNodes {
		SanitizeDataNode(dataNode)
	}
	for _, analyzer := range infonode.Analyzer {
		SanitizeAnalyzerNode(analyzer)
	}
}

func SanitizeDataNode(datanode *InfoNode_DataNode) {
	datanode.NodeType = NodeTypeDataNode
}

func SanitizeAnalyzerNode(analyzerNode *Analyzer) {
	analyzerNode.NodeType = NodeTypeAnalyzerNode
}

func SanitizePackageNode(pkgNode *PackageNode, dbpkgNode *PackageNode) {
	pkgNode.Uid = dbpkgNode.Uid
	pkgNode.Session = dbpkgNode.Session
	pkgNode.NodeType = NodeTypePackageNode
}

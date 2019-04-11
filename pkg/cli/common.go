package cli

import (
	"fmt"

	"github.com/QMSTR/qmstr/pkg/service"
	"golang.org/x/net/context"
)

func getNodesFromArgs(args []string) ([]interface{}, error) {
	var these []interface{}
	for _, arg := range args {
		thisID, err := ParseNodeID(arg)
		if err != nil {
			return nil, fmt.Errorf("Failed parsing node %q: %v", arg, err)
		}
		switch thisVal := thisID.(type) {
		case *service.FileNode:
			this, err := getUniqueFileNode(thisVal)
			if err != nil {
				return nil, err
			}
			these = append(these, this)
		case *service.PackageNode:
			this, err := controlServiceClient.GetPackageNode(context.Background(), &service.PackageNode{Name: thisVal.Name})
			if err != nil {
				return nil, err
			}
			these = append(these, this)
		default:
			return nil, fmt.Errorf("unsupported node type %T", thisVal)
		}
	}
	return these, nil
}

func getUniqueFileNode(fnode *service.FileNode) (*service.FileNode, error) {
	stream, err := controlServiceClient.GetFileNode(context.Background(), &service.GetFileNodeMessage{FileNode: fnode, UniqueNode: true})
	if err != nil {
		return nil, err
	}
	fileNode, err := stream.Recv()
	if err != nil {
		return nil, err
	}
	return fileNode, nil
}

// store nodes in a File node array
func createFileNodesArray(these []interface{}) []*service.FileNode {
	var theseFileNodes []*service.FileNode
	for _, fNode := range these {
		theseFileNodes = append(theseFileNodes, fNode.(*service.FileNode))
	}
	return theseFileNodes
}

func createPkgNodesArray(these []interface{}) []*service.PackageNode {
	var thesePkgNodes []*service.PackageNode
	for _, pkgNode := range these {
		thesePkgNodes = append(thesePkgNodes, pkgNode.(*service.PackageNode))
	}
	return thesePkgNodes
}

func sendFileNode(node *service.FileNode) error {
	stream, err := buildServiceClient.Build(context.Background())
	if err != nil {
		return fmt.Errorf("getting stream for build service fail: %v", err)
	}
	err = stream.Send(node)
	if err != nil {
		return fmt.Errorf("sending node fail: %v", err)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("close stream fail: %v", err)
	}
	if !res.Success {
		return fmt.Errorf("sending node fail: %v", err)
	}
	return nil
}

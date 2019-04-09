package cli

import (
	"fmt"

	"github.com/QMSTR/qmstr/pkg/service"
	"golang.org/x/net/context"
)

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

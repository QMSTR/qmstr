package cli

import (
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

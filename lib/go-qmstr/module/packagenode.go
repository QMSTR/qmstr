package module

import (
	"context"
	"fmt"
	"io"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

type PackageNodeProxy struct {
	service.PackageNode
	masterClient *MasterClient
}

func (pnp *PackageNodeProxy) GetTargets() []*FileNodeProxy {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if len(pnp.Targets) == 0 {
		pnp.masterClient.ctrlSvcClient.GetPackageTargets(ctx, &pnp.PackageNode)
	}
	return nil
}

func (m *MasterClient) GetPackageNodes() ([]*PackageNodeProxy, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pnps := []*PackageNodeProxy{}
	pkgStream, err := m.ctrlSvcClient.GetPackageNode(ctx, &service.PackageNode{})
	if err != nil {
		return nil, fmt.Errorf("Couldn't get package node, %v", err)
	}

	for {
		pkg, err := pkgStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error receiving package nodes, %v", err)
		}
		pnps = append(pnps, &PackageNodeProxy{*pkg, m})
	}
	return pnps, nil
}

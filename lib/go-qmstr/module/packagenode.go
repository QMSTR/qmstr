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

func (pnp *PackageNodeProxy) GetMasterClient() *MasterClient {
	return pnp.masterClient
}

func (pnp *PackageNodeProxy) GetTargets() ([]*FileNodeProxy, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	targetStream, err := pnp.masterClient.CtrlSvcClient.GetPackageTargets(ctx, &pnp.PackageNode)
	if err != nil {
		return nil, err
	}
	targets := []*FileNodeProxy{}
	for {
		target, err := targetStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("Error receiving package nodes, %v", err)
		}
		targets = append(targets, &FileNodeProxy{*target, pnp.masterClient})
	}
	return targets, nil
}

func (m *MasterClient) GetPackageNodes() ([]*PackageNodeProxy, error) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pnps := []*PackageNodeProxy{}
	pkgStream, err := m.CtrlSvcClient.GetPackageNode(ctx, &service.PackageNode{})
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

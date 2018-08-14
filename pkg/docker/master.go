package docker

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// GetMasterInfo returns the container id and the internal qmstr port of the master container
// configured via QMSTR_MASTER environment variable
func GetMasterInfo(ctx context.Context, cli *client.Client) (string, uint16, error) {
	qmstrAddr := os.Getenv("QMSTR_MASTER")
	if qmstrAddr == "" {
		return "", 0, errors.New("QMSTR_MASTER not set, can't determine qmstr-master container")
	}
	qmstrAddrS := strings.Split(qmstrAddr, ":")
	qmstrHostPort, err := strconv.ParseUint(qmstrAddrS[len(qmstrAddrS)-1], 10, 64)
	if err != nil {
		return "", 0, err
	}

	args, err := filters.ParseFlag("label=org.qmstr.image", filters.NewArgs())
	if err != nil {
		return "", 0, err
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: args})
	if err != nil {
		return "", 0, err
	}

	for _, container := range containers {
		for _, portCfg := range container.Ports {
			if uint64(portCfg.PublicPort) == qmstrHostPort {
				return container.ID, portCfg.PrivatePort, nil
			}
		}
	}

	return "", 0, errors.New("qmstr-master container not found")
}

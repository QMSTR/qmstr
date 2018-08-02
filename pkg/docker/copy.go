package docker

import (
	"context"

	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func copyResults(ctx context.Context, cli *client.Client, container string, destinationPath string) error {

	data, stat, err := cli.CopyFromContainer(ctx, container, master.ServerOutputDir)
	if err != nil {
		return err
	}
	defer data.Close()

	srcInfo := archive.CopyInfo{
		Path:   master.ServerOutputDir,
		Exists: true,
		IsDir:  stat.Mode.IsDir(),
	}

	return archive.CopyTo(data, srcInfo, destinationPath)
}

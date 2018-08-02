package docker

import (
	"bytes"
	"context"
	"io/ioutil"
	"regexp"

	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/archive"
)

func CopyResults(ctx context.Context, cli *client.Client, container string, destinationPath string) error {
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

func GetMasterConfig(ctx context.Context, cli *client.Client, container string) ([]byte, error) {
	data, _, err := cli.CopyFromContainer(ctx, container, "/qmstr/qmstr.yaml")
	if err != nil {
		return nil, err
	}
	defer data.Close()

	config, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}

	// poor man's untar this might be replaced by a proper untar function
	tarPattern := regexp.MustCompile("^.*package:")
	config = tarPattern.ReplaceAll(config, []byte("package:"))
	config = bytes.Trim(config, "\x00")

	return config, nil
}

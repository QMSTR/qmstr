package docker

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type ClientContainer struct {
	Image             string
	MasterContainerID string
	QmstrInternalPort uint16
	Instdir           string
	Cmd               []string
	Env               []string
	Mount             []mount.Mount
}

func RunClientContainer(ctx context.Context, cli *client.Client, clientConfig *ClientContainer) error {
	log.Printf("connecting to master container %s", clientConfig.MasterContainerID)

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("unable to determine current working directory")
	}
	hostConf := &container.HostConfig{
		Mounts:      append([]mount.Mount{mount.Mount{Source: wd, Target: common.ContainerBuildDir, Type: mount.TypeBind}}, clientConfig.Mount...),
		NetworkMode: container.NetworkMode(fmt.Sprintf("container:%s", clientConfig.MasterContainerID)),
	}

	containerCmd := []string{"qmstr"}
	if clientConfig.Instdir != "" {
		containerCmd = append(containerCmd, fmt.Sprintf("--instdir=%s", clientConfig.Instdir))
	}
	containerCmd = append(containerCmd, append([]string{"--"}, clientConfig.Cmd...)...)

	clientConfig.Env = append([]string{fmt.Sprintf("QMSTR_MASTER=%s:%d", clientConfig.MasterContainerID[:12], clientConfig.QmstrInternalPort)}, clientConfig.Env...)

	containerConf := &container.Config{
		Image: clientConfig.Image,
		Cmd:   containerCmd,
		Tty:   true,
		Env:   clientConfig.Env,
	}

	user, err := user.Current()
	if err == nil {
		containerConf.User = user.Uid
	}

	resp, err := cli.ContainerCreate(ctx, containerConf, hostConf, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	status, err := cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		return err
	}

	log.Printf("Build container returned status %d", status)

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	logmsg, err := ioutil.ReadAll(out)
	if err != nil {
		return err
	}
	log.Printf("Container logs:\n%s", logmsg)

	if status != 0 {
		os.Exit(int(status))
	}

	return nil
}

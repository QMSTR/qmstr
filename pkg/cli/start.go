package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

const retryLimit = 16
const minPort = 6787 // MSTR on the dial pad
const portCount = 1024
const maxPort = minPort + portCount
const containerBuildDir = "/buildroot"
const proto = "tcp"

var hostPortRange = fmt.Sprintf("%d-%d", minPort, maxPort)
var masterImageName = "qmstr/master"
var internalMasterPort string
var wait bool

func startMaster(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		Log.Fatalf("unable to determine current working directory")
	}

	config, err := config.ReadConfigFromFile(filepath.Join(wd, "qmstr.yaml"))
	if err != nil {
		Log.Fatalf("failed to read configuration %v", err)
	}

	configuredImageName := config.Server.ImageName
	if configuredImageName != "" {
		Debug.Printf("using configured image %s", configuredImageName)
		masterImageName = configuredImageName
	}

	internalMasterPort, err = config.GetRPCPort()
	if err != nil {
		Log.Fatalf("failed to get configured rpc port %v", err)
	}

	Debug.Printf("Starting Quartermaster master")

	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		Log.Fatalf("Failed to create docker client %v", err)
	}

	containerID, portBinding, err := startContainer(ctx, cli, wd)
	if err != nil {
		if containerID != "" {
			cleanUpContainer(ctx, cli, containerID, verbose)
		}
		Log.Fatalf("Starting qmstr-master failed: %v", err)
	}

	var address = fmt.Sprintf("%s:%s", portBinding.HostIP, portBinding.HostPort)

	if wait {
		setUpServer()
		defer tearDownServer()
		awaitServer()
	}

	fmt.Printf("export QMSTR_MASTER=%s\n", address)
	Debug.Println("Done.")
}

func startContainer(ctx context.Context, cli *client.Client, workdir string) (string, *nat.PortBinding, error) {

	internalPort, err := nat.NewPort(proto, internalMasterPort)
	if err != nil {
		return "", nil, err
	}

	portsbinds := []nat.PortBinding{nat.PortBinding{HostIP: "0.0.0.0", HostPort: hostPortRange}}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: masterImageName,
		ExposedPorts: nat.PortSet{
			internalPort: struct{}{},
		},
	},
		&container.HostConfig{
			PortBindings: nat.PortMap{internalPort: portsbinds},
			Mounts:       []mount.Mount{mount.Mount{Source: workdir, Target: containerBuildDir, Type: mount.TypeBind}},
		}, nil, "")
	if err != nil {
		return "", nil, err
	}

	Debug.Printf("Start container %v", resp.ID)
	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return resp.ID, nil, err
	}

	// give the container some time to settle
	timer := time.NewTimer(time.Second * 2)
	<-timer.C

	info, err := cli.ContainerInspect(ctx, resp.ID)
	if err != nil {
		return resp.ID, nil, err
	}

	if !info.State.Running {
		return resp.ID, nil, errors.New("container not running")
	}

	portBindings := info.NetworkSettings.Ports[internalPort]

	// there can be only one binding per port and protocol
	if len(portBindings) < 1 {
		return resp.ID, nil, errors.New("no port binding found")
	}
	return resp.ID, &portBindings[0], nil
}

func cleanUpContainer(ctx context.Context, cli *client.Client, containerID string, verbose bool) {
	info, _ := cli.ContainerInspect(ctx, containerID)
	if info.State.Running {
		cli.ContainerStop(ctx, containerID, nil)
	}
	if verbose {
		out, err := cli.ContainerLogs(ctx, containerID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
		if err == nil {
			defer out.Close()
			buf := new(bytes.Buffer)
			buf.ReadFrom(out)
			Debug.Printf("qmstr-master log\n%s\n", buf.String())
		}
	}
	cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the Quartermaster master",
	Long:  fmt.Sprintf("Start the Quartermaster master at a random port in the range between %d and %d.", minPort, maxPort),
	Run:   startMaster,
}

func init() {
	rootCmd.AddCommand(startCmd)
	startCmd.Flags().BoolVar(&wait, "wait", false, "wait for qmstr-master")
}

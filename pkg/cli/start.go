package cli

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/master"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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
const internalConfigPath = "/qmstr/qmstr.yaml"
const globalHostConfigPath = "/etc/qmstr/qmstr.yaml"

var hostPortRange = fmt.Sprintf("%d-%d", minPort, maxPort)
var masterImageName = "qmstr/master"

var internalMasterPort string

var configFile string
var wait bool
var debug bool
var seed string

func getConfig() (*config.MasterConfig, error) {
	var err error

	if _, err := os.Stat(configFile); err != nil {
		return nil, err
	}

	configFile, err = filepath.Abs(configFile)
	if err != nil {
		return nil, err
	}

	config, err := config.ReadConfigFromFiles(globalHostConfigPath, configFile)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func startMaster(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		Log.Fatalf("unable to determine current working directory")
	}

	config, err := getConfig()
	if err != nil {
		Log.Fatalf("failed to read configuration %v", err)
	}

	if config.Server.BuildPath == "" {
		config.Server.BuildPath = containerBuildDir
	}

	configuredImageName := config.Server.ImageName
	if configuredImageName != "" {
		Debug.Printf("using configured image %s", configuredImageName)
		masterImageName = configuredImageName
	}

	debug = config.Server.Debug

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

	containerNetwork, err := createNetwork(ctx, cli)
	if err != nil {
		Log.Fatalf("Failed to create container network %v", err)
	}

	containerID, portBinding, err := startContainer(ctx, cli, wd, containerNetwork, config)
	if err != nil {
		if containerID != "" {
			cleanUpContainer(ctx, cli, containerID, verbose)
		}
		Log.Fatalf("Starting qmstr-master failed: %v", err)
	}

	address = fmt.Sprintf("%s:%s", portBinding.HostIP, portBinding.HostPort)

	if wait {
		setUpControlService()
		defer tearDownServer()
		awaitServer()
	}

	fmt.Printf("export %s=%s\n", common.QMSTRADDRENV, address)
	Debug.Println("Done.")
}

func createNetwork(ctx context.Context, cli *client.Client) (string, error) {
	cleanUpContainerNetworks(ctx, cli)
	// find qmstr networks
	args, err := filters.ParseFlag("label=org.qmstr.network=true", filters.NewArgs())
	if err != nil {
		return "", err
	}
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{Filters: args})
	if err != nil {
		return "", err
	}
	netIDs := []uint64{0}
	for _, net := range networks {
		Debug.Printf("found qmstr net %s", net.Name)
		netSplit := strings.Split(net.Name, "-")
		if len(netSplit) > 1 {
			id, err := strconv.ParseUint(netSplit[1], 10, 64)
			if err != nil {
				return "", err
			}
			netIDs = append(netIDs, id)
		}
	}
	// sort descending
	sort.Slice(netIDs, func(x, y int) bool { return netIDs[x] > netIDs[y] })

	// create new qmstr network
	containerNetName := fmt.Sprintf("qmstrnet-%d", netIDs[0]+1)

	_, err = cli.NetworkCreate(ctx, containerNetName, types.NetworkCreate{CheckDuplicate: true, Labels: map[string]string{"org.qmstr.network": "true"}})
	if err != nil {
		return "", err
	}
	Log.Printf("Container network %s created", containerNetName)
	return containerNetName, nil
}

func getHostConfig(internalPort nat.Port, workdir string, network string, extraMount map[string]string) *container.HostConfig {
	portsbinds := []nat.PortBinding{nat.PortBinding{HostIP: "0.0.0.0", HostPort: hostPortRange}}
	hostConf := &container.HostConfig{
		PortBindings: nat.PortMap{internalPort: portsbinds},
		Mounts: []mount.Mount{
			mount.Mount{Source: workdir, Target: containerBuildDir, Type: mount.TypeBind},
		},
		NetworkMode: container.NetworkMode(network),
	}

	// for debugging container must run qmstr-master via dlv and so the container needs to be allowed to fork
	if debug {
		hostConf.Privileged = true
		hostConf.SecurityOpt = []string{"seccomp=unconfined"}
	}

	for target, source := range extraMount {
		Debug.Printf("Mounting host directory %s to %s inside container", source, target)
		hostConf.Mounts = append(hostConf.Mounts, mount.Mount{Source: source, Target: target, Type: mount.TypeBind})
	}

	return hostConf
}

func getContainerConfig(internalPort nat.Port, workdir string, extraEnv map[string]string) (*container.Config, error) {
	config := &container.Config{
		Image: masterImageName,
		ExposedPorts: nat.PortSet{
			internalPort: struct{}{},
		},
		Env: []string{fmt.Sprintf("PATH_SUB=%s,%s", workdir, containerBuildDir)},
		Tty: true,
	}

	for envk, envv := range extraEnv {
		config.Env = append(config.Env, fmt.Sprintf("%s=%s", envk, envv))
	}

	return config, nil
}

func startContainer(ctx context.Context, cli *client.Client, workdir string, network string,
	masterConfig *config.MasterConfig) (string, *nat.PortBinding, error) {

	extraEnv := masterConfig.Server.ExtraEnv
	extraMount := masterConfig.Server.ExtraMount

	internalPort, err := nat.NewPort(proto, internalMasterPort)
	if err != nil {
		return "", nil, err
	}

	if masterConfig.Server.CacheDir != "" {
		extraMount[master.ServerCacheDir] = masterConfig.Server.CacheDir
	}

	user, err := user.Current()
	if err == nil {
		extraEnv["USERID"] = user.Uid
	}

	if debug {
		gopathTarget := "/go"
		gopath := os.Getenv("GOPATH")
		if gopath == "" {
			return "", nil, errors.New("GOPATH not set")
		}
		extraEnv["GOPATH"] = gopathTarget
		extraEnv["QMSTR_DEBUG"] = "true"
		extraMount[filepath.Join(gopathTarget, "src")] = filepath.Join(gopath, "src")
	}

	if seed != "" {
		extraMount[common.ContainerGraphImportPath] = filepath.Join(workdir, seed)
	}

	containerConfig, err := getContainerConfig(internalPort, workdir, extraEnv)
	if err != nil {
		return "", nil, err
	}
	hostConf := getHostConfig(internalPort, workdir, network, extraMount)

	resp, err := cli.ContainerCreate(ctx, containerConfig, hostConf, nil, "")
	if err != nil {
		return "", nil, err
	}

	// copy config over
	configData, err := config.SerializeConfig(masterConfig)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to serialize configuration in order to pass it to qmstr-master: %v", err)
	}
	err = docker.WriteContainerFile(ctx, cli, configData, resp.ID, internalConfigPath)
	if err != nil {
		return "", nil, fmt.Errorf("Failed to write configuration to qmstr-master container: %v", err)
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

	if debug {
		Log.Println("WARNING: Running qmstr-master container in privileged mode for debugging.")
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
	if !debug {
		cli.ContainerRemove(ctx, containerID, types.ContainerRemoveOptions{Force: true})
	}
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
	startCmd.Flags().IntVarP(&timeout, "timeout", "t", 60, "timeout after the specified time (seconds). Used after the wait flag")
	startCmd.Flags().StringVarP(&configFile, "config", "c", "qmstr.yaml", "Path to qmstr configuration file")
	startCmd.Flags().StringVar(&seed, "seed", "", "Replay dgraph export on init")
}

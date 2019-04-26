package cliqmstr

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/docker"

	"github.com/spf13/cobra"
)

var (
	container string // Image to spawn container from
)

var spawnCmd = &cobra.Command{
	Use:   "spawn",
	Short: "Spawn a container adjacent to the master and execute a command",
	Long:  `spawn starts a container in the same container network as the master, and executes the specified arguments.`,
	Args:  cobra.MinimumNArgs(2),
	Run:   executeSpawn,
}

func init() {
	rootCmd.AddCommand(spawnCmd)
	spawnCmd.Flags().SetInterspersed(false) // this stops flag parsing at the first argument, allowing for e.g. "ls -la" as payload
}

func executeSpawn(cmd *cobra.Command, args []string) {
	container := args[0]
	commands := args[1:]
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		Log.Fatalf("Failed to create docker client %v", err)
	}
	masterContainerID, intPort, err := docker.GetMasterInfo(ctx, cli)
	if err != nil {
		Log.Fatalf("Unable to find qmstr-master container")
	}

	var env []string
	if val, ok := os.LookupEnv(common.QMSTRDEBUGENV); ok {
		env = append(env, fmt.Sprintf("%s=%s", common.QMSTRDEBUGENV, val))
	}

	var mountpoints []mount.Mount
	if val, ok := os.LookupEnv(common.CCACHEDIRENV); ok {
		env = append(env, fmt.Sprintf("%s=%s", common.CCACHEDIRENV, common.ContainerCcacheDir))
		mountpoints = append(mountpoints, mount.Mount{Type: mount.TypeBind, Source: val, Target: common.ContainerCcacheDir})
	}

	Log.Printf("starting build container")
	err = docker.RunClientContainer(ctx, cli, &docker.ClientContainer{
		Image:             container,
		Cmd:               commands,
		MasterContainerID: masterContainerID,
		QmstrInternalPort: intPort,
		Env:               env,
		Mount:             mountpoints,
	})
	if err != nil {
		Log.Fatalf("Build container failed: %v", err)
	}
}

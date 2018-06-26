package cli

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func cleanQmstr(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		Log.Fatalf("Failed to create docker client %v", err)
	}
	cleanUpContainers(ctx, cli)
	cleanUpContainerNetworks(ctx, cli)
}

func cleanUpContainers(ctx context.Context, cli *client.Client) error {
	args, err := filters.ParseFlag("label=org.qmstr.image", filters.NewArgs())
	if err != nil {
		return err
	}

	resp, err := cli.ContainersPrune(ctx, args)
	if err != nil {
		return err
	}
	Debug.Printf("Deleted %v", resp.ContainersDeleted)
	return nil
}

func cleanUpContainerNetworks(ctx context.Context, cli *client.Client) error {
	// find qmstr networks
	args, err := filters.ParseFlag("label=org.qmstr.network=true", filters.NewArgs())
	if err != nil {
		return err
	}
	networks, err := cli.NetworkList(ctx, types.NetworkListOptions{Filters: args})
	if err != nil {
		return err
	}

	for _, net := range networks {
		Debug.Printf("found qmstr net %s", net.Name)
		if len(net.Containers) == 0 {
			Debug.Printf("Remove unused qmstr network %s", net.Name)
			cli.NetworkRemove(ctx, net.ID)
		}

	}
	return nil
}

var cleanCmd = &cobra.Command{
	Use:   "clean",
	Short: "remove dangling qmstr containers",
	Long:  "Remove all stopped qmstr containers",
	Run:   cleanQmstr,
}

func init() {
	rootCmd.AddCommand(cleanCmd)
}

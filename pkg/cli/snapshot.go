package cli

import (
	"fmt"
	"os"
	"path"

	"github.com/QMSTR/qmstr/pkg/docker"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Take snapshot of qmstr-master server",
	Long:  `Create a database export of the graph database`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpServer()
		if err := exportGraph(); err != nil {
			Log.Fatalf("snapshot creation failed: %v", err)
		}
		if err := copyExport(); err != nil {
			Log.Fatalf("copying snapshot failed: %v", err)
		}
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
}

func exportGraph() error {
	_, err := controlServiceClient.ExportGraph(context.Background(), &service.ExportRequest{Wait: true})
	if err != nil {
		return fmt.Errorf("Failed to export graph: %v", err)
	}
	return nil
}

func copyExport() error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("failed to create docker client %v", err)
	}
	mID, _, err := docker.GetMasterInfo(ctx, cli)
	if err != nil {
		return fmt.Errorf("failed to obtain qmstr-master info %v", err)
	}

	outdir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("unable to determine current working directory: %v", err)
	}
	return docker.CopyGraphExport(ctx, cli, mID, path.Join(outdir, "qmstr-snapshot.tar"))
}

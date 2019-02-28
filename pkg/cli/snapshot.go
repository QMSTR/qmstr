package cli

import (
	"fmt"
	"os"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/docker"
	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var snapshotFile string
var forceOverride bool

var snapshotCmd = &cobra.Command{
	Use:   "snapshot",
	Short: "Take snapshot of qmstr-master server",
	Long:  `Create a database export of the graph database`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpControlService()
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
	snapshotCmd.Flags().StringVarP(&snapshotFile, "out", "O", "qmstr-snapshot.tar", "Output filename")
	snapshotCmd.Flags().BoolVarP(&forceOverride, "force", "f", false, "force override snapshot")
}

func exportGraph() error {
	_, err := controlServiceClient.ExportSnapshot(context.Background(), &service.ExportRequest{Wait: true})
	if err != nil {
		return fmt.Errorf("Failed to export snapshot: %v", err)
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

	//file already exists
	if common.IsFileExist(snapshotFile) {
		if !forceOverride {
			return fmt.Errorf("snapshot %s already exists; use -f to overwrite", snapshotFile)
		}
		err := os.Remove(snapshotFile)
		if err != nil {
			return fmt.Errorf("failed to remove snapshot %v", err)
		}
	}
	return docker.CopySnapshot(ctx, cli, mID, snapshotFile)
}

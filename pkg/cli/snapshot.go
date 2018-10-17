package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/QMSTR/qmstr/pkg/docker"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
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
	snapshotCmd.Flags().StringVarP(&snapshotFile, "out", "O", "qmstr-snapshot.tar", "Output filename")
	snapshotCmd.Flags().BoolVarP(&forceOverride, "force", "f", false, "force override snapshot")
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

	//file already exists
	if _, err := os.Stat(snapshotFile); err == nil {
		if forceOverride {
			if err = os.Remove(snapshotFile); err != nil {
				return fmt.Errorf("couldn't remove file: %s, %v", snapshotFile, err)
			}
		} else {
			reader := bufio.NewReader(os.Stdin)
			Log.Printf("File %s already exists. Would you like to override it? [y/N]", snapshotFile)
			answer, _, err := reader.ReadRune()
			if err != nil {
				return err
			}
			switch answer {
			case 'y', 'Y':
				if err = os.Remove(snapshotFile); err != nil {
					return fmt.Errorf("couldn't remove file: %s, %v", snapshotFile, err)
				}
			default:
				reader = bufio.NewReader(os.Stdin)
				Log.Print("Please provide a different snapshot name (followed by the 'tar' extension)")
				newFileName, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				snapshotFile = strings.Replace(newFileName, "\n", "", -1)
				Log.Printf("Creating new snapshot: %s", snapshotFile)
			}
		}
	}
	return docker.CopyGraphExport(ctx, cli, mID, snapshotFile)
}

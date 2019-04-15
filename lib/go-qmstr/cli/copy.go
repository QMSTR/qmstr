package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/lib/go-qmstr/config"
	"github.com/QMSTR/qmstr/lib/go-qmstr/docker"
	"github.com/docker/docker/client"

	"golang.org/x/net/context"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy results from qmstr",
	Long:  `Copy report results from qmstr server.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := copyResults(); err != nil {
			Log.Fatalf("copying results failed: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}

func copyResults() error {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		return fmt.Errorf("failed to create docker client %v", err)
	}
	mID, _, err := docker.GetMasterInfo(ctx, cli)
	if err != nil {
		return fmt.Errorf("failed to obtain qmstr-master info %v", err)
	}

	configdata, err := docker.GetMasterConfig(ctx, cli, mID, internalConfigPath)
	if err != nil {
		return fmt.Errorf("can not load master configuration from container: %v", err)
	}
	Debug.Printf("Got config from master: %s", configdata)
	config, err := config.ReadConfigFromBytes(configdata)
	if err != nil {
		return fmt.Errorf("can not read master configuration from container: %v", err)
	}

	outdir := config.Server.OutputDir
	if outdir == "" {
		var err error
		outdir, err = os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to determine current working directory: %v", err)
		}
		outdir = filepath.Join(outdir, "qmstr")
	}

	// clean outdir if it exists
	if err = os.RemoveAll(outdir); err != nil {
		return fmt.Errorf("failed to clean result directory: %v", err)
	}

	Debug.Printf("writing results to %s", outdir)
	return docker.CopyResults(ctx, cli, mID, outdir)
}

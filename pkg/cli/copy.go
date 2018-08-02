package cli

import (
	"log"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/docker"
	"github.com/docker/docker/client"

	"golang.org/x/net/context"

	"github.com/spf13/cobra"
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy results from qmstr",
	Long:  `Copy report results from qmstr server.`,
	Run: func(cmd *cobra.Command, args []string) {
		copyResults()
	},
}

func init() {
	rootCmd.AddCommand(copyCmd)
}

func copyResults() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		Log.Fatalf("Failed to create docker client %v", err)
	}
	mID, _, err := docker.GetMasterInfo(ctx, cli)
	if err != nil {
		log.Fatal(err)
	}

	configdata, err := docker.GetMasterConfig(ctx, cli, mID)
	if err != nil {
		Log.Fatalf("Can not load master configuration from container : %v", err)
	}
	Debug.Printf("Got config from master: %s", configdata)
	config, err := config.ReadConfig(configdata)
	if err != nil {
		Log.Fatalf("Can not read master configuration from container : %v", err)
	}

	outdir := config.Server.OutputDir
	if outdir == "" {
		var err error
		outdir, err = os.Getwd()
		if err != nil {
			Log.Fatalf("unable to determine current working directory")
		}
		outdir = filepath.Join(outdir, "qmstr")
	}

	Debug.Printf("writing results to %s", outdir)
	docker.CopyResults(ctx, cli, mID, outdir)
}

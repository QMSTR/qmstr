package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/QMSTR/qmstr/pkg/docker"
	"github.com/docker/docker/client"

	"golang.org/x/net/context"

	pb "github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
)

var force, nocopy bool

// quitCmd represents the quit command
var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: "Quit qmstr",
	Long:  `Run quit if you want to quit qmstr.`,
	Run: func(cmd *cobra.Command, args []string) {
		if !nocopy {
			if err := copyResults(); err != nil {
				Log.Fatalf("copying results failed: %v", err)
			}
		}
		setUpControlService()
		quitServer()
		tearDownServer()
		stopMasterContainer()
	},
}

func init() {
	rootCmd.AddCommand(quitCmd)
	quitCmd.Flags().BoolVarP(&force, "force", "f", false, "force quit")
	quitCmd.Flags().BoolVar(&nocopy, "no-copy", false, "Do not copy results")
}

func quitServer() {
	resp, err := controlServiceClient.Quit(context.Background(), &pb.QuitMessage{Kill: force})
	if err != nil {
		Log.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	if !resp.Success {
		Log.Println("Server responded unsuccessful")
		os.Exit(ReturnCodeResponseFalseError)
	}
}

func stopMasterContainer() {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		Log.Fatalf("Failed to create docker client %v", err)
	}
	mID, _, err := docker.GetMasterInfo(ctx, cli)
	if err != nil {
		Log.Fatal(err)
	}
	d := time.Duration(2) * time.Second
	err = cli.ContainerStop(ctx, mID, &d)
	if err != nil {
		err1 := cli.ContainerKill(ctx, mID, "SIGKILL")
		if err1 != nil {
			Log.Fatal(fmt.Errorf("%v : %v", err1, err))
		}
	}
}

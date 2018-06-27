package cli

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/QMSTR/qmstr/pkg/docker"
	"github.com/docker/docker/client"

	"golang.org/x/net/context"

	pb "github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
)

var force bool

// quitCmd represents the quit command
var quitCmd = &cobra.Command{
	Use:   "quit",
	Short: "Quit qmstr",
	Long:  `Run quit if you want to quit qmstr.`,
	Run: func(cmd *cobra.Command, args []string) {

		setUpServer()
		quitServer()
		tearDownServer()
		stopMasterContainer()
	},
}

func init() {
	rootCmd.AddCommand(quitCmd)
	quitCmd.Flags().BoolVarP(&force, "force", "f", false, "force quit")
}

func quitServer() {
	resp, err := controlServiceClient.Quit(context.Background(), &pb.QuitMessage{Kill: force})
	if err != nil {
		fmt.Printf("Failed to communicate with qmstr-master server. %v\n", err)
		os.Exit(ReturnCodeServerCommunicationError)
	}
	if !resp.Success {
		fmt.Println("Server responded unsuccessful")
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
		log.Fatal(err)
	}
	d := time.Duration(2) * time.Second
	err = cli.ContainerStop(ctx, mID, &d)
	if err != nil {
		err1 := cli.ContainerKill(ctx, mID, "SIGKILL")
		if err1 != nil {
			log.Fatal(fmt.Errorf("%v : %v", err1, err))
		}
	}
}

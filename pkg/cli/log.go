package cli

import (
	"bufio"
	"fmt"

	"github.com/QMSTR/qmstr/pkg/docker"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

func logQmstr(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	cli, err := client.NewEnvClient()
	if err != nil {
		Log.Fatalf("Failed to create docker client %v", err)
	}
	err = logMasterContainer(ctx, cli, follow)
	if err != nil {
		Log.Fatalf("Failed print master logs %v", err)
	}
}

func logMasterContainer(ctx context.Context, cli *client.Client, follow bool) error {
	masterID, _, err := docker.GetMasterInfo(ctx, cli)
	if err != nil {
		return err
	}
	logReader, err := cli.ContainerLogs(ctx, masterID,
		types.ContainerLogsOptions{ShowStderr: true, ShowStdout: true, Follow: follow})
	if err != nil {
		return err
	}
	defer logReader.Close()

	bufLogReader := bufio.NewReader(logReader)
	for err == nil {
		var line []byte
		line, _, err = bufLogReader.ReadLine()
		fmt.Printf("%s\n", line)
	}
	if err.Error() == "EOF" {
		err = nil
	}
	return err
}

var logsCmd = &cobra.Command{
	Use:   "logs",
	Short: "print qmstr-master logs",
	Long:  "print qmstr-master container logs",
	Run:   logQmstr,
}

func init() {
	rootCmd.AddCommand(logsCmd)
	logsCmd.Flags().BoolVar(&follow, "follow", false, "follow logs")
}

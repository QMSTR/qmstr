package cli

import (
	"fmt"
	"os"
	"time"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var timeout int

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "waits until the server is ready",
	Long: `waits until the server is ready. If 60 seconds pass 
	then it exits with error message.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpServer()
		awaitServer()
		tearDownServer()
	},
}

func init() {
	rootCmd.AddCommand(waitCmd)
	waitCmd.Flags().IntVarP(&timeout, "timeout", "t", 60, "time is up")

}

func awaitServer() {
	go func() {
		<-time.After(time.Second * time.Duration(timeout))
		fmt.Printf("wait for qmstr-master timed out after %d seconds\n", timeout)
		os.Exit(ReturnCodeTimeout)
	}()
	for {
		res, err := buildServiceClient.Log(context.Background(), &pb.LogMessage{Msg: []byte("Client is waiting for qmstr server to be ready")})
		if err != nil {
			continue
		}
		if res.Success {
			return
		}
	}

}

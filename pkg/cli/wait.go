package cli

import (
	"os"
	"time"

	"github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var timeout int

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "waits until the server is ready",
	Long: `waits until the server is ready. If 60 seconds pass 
	then it exits with an error message.`,
	Run: func(cmd *cobra.Command, args []string) {
		setUpControlService()
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
		Log.Printf("wait for qmstr-master timed out after %d seconds\n", timeout)
		os.Exit(ReturnCodeTimeout)
	}()
	for {
		res, err := controlServiceClient.Status(context.Background(), &service.StatusMessage{})
		if err != nil {
			Debug.Println("Master server not yet available")
			<-time.After(time.Second * time.Duration(1))
			Debug.Println("retrying")
			continue
		}
		if res.PhaseID == service.Phase_FAIL {
			Log.Println("qmstr master is in failure state")
			os.Exit(ReturnCodeServerFailureError)
		}
		if res.PhaseID > service.Phase_INIT {
			return
		}
	}

}

package cli

import (
	"math"
	"os"
	"time"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
	"github.com/spf13/cobra"
	"golang.org/x/net/context"
)

var timeout int

var waitCmd = &cobra.Command{
	Use:   "wait",
	Short: "Wait until the QMSTR master is ready",
	Long:  `Wait until the QMSTR master is ready. Waiting will time out with an error after the timeout period has elapsed.`,
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
	var lastPending uint64 = math.MaxUint64
	timer := time.NewTimer(time.Second * time.Duration(timeout))
	cancel := make(chan interface{})
	defer func() {
		cancel <- nil
	}()
	go func() {
		select {
		case <-timer.C:
			Log.Printf("wait for qmstr-master timed out after %d seconds\n", timeout)
			os.Exit(ReturnCodeTimeout)
		case <-cancel:
			//time out canceled
			timer.Stop()
			return
		}
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
			if lastPending > res.PendingInserts {
				timer.Reset(time.Second * time.Duration(timeout))
				lastPending = res.PendingInserts
				Debug.Printf("Pending inserts: %d", res.PendingInserts)
				<-time.After(time.Second * time.Duration(1))
				continue
			}

			if res.PendingInserts == 0 {
				return
			}

			<-time.After(time.Second * time.Duration(1))
			continue
		}
	}

}

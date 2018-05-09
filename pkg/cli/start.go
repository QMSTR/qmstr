package cli

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

const retryLimit = 16
const minPort = 7686 // MSTR on the dial pad
const portCount = 1024
const maxPort = minPort + portCount

func startMaster(cmd *cobra.Command, args []string) {
	rand.Seed(time.Now().UnixNano())
	for retry := 0; retry < retryLimit; retry++ {
		port := rand.Intn(portCount) + minPort
		wd, err := os.Getwd()
		if err != nil {
			Log.Println("unable to determine current working directory")
			os.Exit(1)
		}
		Debug.Printf("Starting Quartermaster master with port %d...\n", port)
		cmd := exec.Command("docker", "run", "-d", "-p", fmt.Sprintf("%d:50051", port),
			"-v", fmt.Sprintf("%s:/buildroot", wd), "qmstr/master")
		Debug.Printf("master command line: \"%s\".\n", strings.Join(cmd.Args, " "))
		err = cmd.Run()
		switch value := err.(type) {
		case *exec.ExitError:
			ws := value.Sys().(syscall.WaitStatus)
			Log.Printf("Master failed to start: %v (exit status %d).\nRetrying to see if there was a problem allocating the port.\n", err, ws.ExitStatus())
		default:
			fmt.Printf("export QMSTR_MASTER=localhost:%d\n", port)
			Debug.Println("Done.")
			return
		}
	}
	Log.Printf("Error starting the Quartermaster master. I retried %d times.\n", retryLimit)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the Quartermaster master",
	Long:  fmt.Sprintf("Start the Quartermaster master at a random port in the range between %d and %d.", minPort, maxPort),
	Run:   startMaster,
}

func init() {
	AddressOptional = true
	rootCmd.AddCommand(startCmd)
}

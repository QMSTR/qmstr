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
			fmt.Println("unable to determine current working directory")
			os.Exit(1)
		}
		fmt.Printf("Starting Quartermaster master with port %d...\n", port)
		cmd := exec.Command("docker", "run", "-d", "-p", fmt.Sprintf("%d:50051", port),
			"-v", fmt.Sprintf("%s:/buildroot", wd), "qmstr/master")
		err = cmd.Run()
		switch value := err.(type) {
		case *exec.ExitError:
			ws := value.Sys().(syscall.WaitStatus)
			fmt.Printf("Master failed to start: %v (exit status %d).\nRetrying to see if there was a problem allocating the port.\n", err, ws.ExitStatus())
			fmt.Printf("Note: master command line was \"%s\".\n", strings.Join(cmd.Args, " "))
		default:
			fmt.Printf("export QMSTR_MASTER=localhost:%d\n", port)
			fmt.Println("Done.")
			return
		}
	}
	fmt.Printf("Error starting the Quartermaster master. I retried %d times.\n", retryLimit)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the Quartermaster master",
	Long:  fmt.Sprintf("Start the Quartermaster master at a random port in the range between %d and %d.", minPort, maxPort),
	Run:   startMaster,
}

func init() {
	rootCmd.AddCommand(startCmd)
}

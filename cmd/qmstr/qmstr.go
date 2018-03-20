package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

// Options contains the context of a program invocation
type Options struct {
	progName           string //The name the program is called as
	keepTmpDirectories bool   //Keep intermediate files
	verbose            bool   //Enable trace log output
}

var options Options

func main() {
	options.progName = os.Args[0]
	flag.BoolVar(&options.keepTmpDirectories, "keep", false, "Keep the created directories instead of cleaning up.")
	flag.BoolVar(&options.verbose, "verbose", false, "Enable diagnostic log output.")
	flag.Parse()
	arguments := flag.Args()
	if len(arguments) == 0 && !options.keepTmpDirectories {
		usage("No command specified!")
	}
	exitCode := Run(arguments)
	os.Exit(exitCode)
}

func usage(format string, a ...interface{}) {
	message := fmt.Sprintf(format, a...)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(message))
	fmt.Fprintf(os.Stderr, "Usage: %s <flags> [working directory]\n", options.progName)
	flag.PrintDefaults()
	os.Exit(1)
}

// Run does everything
// It also makes sure that even though os.Exit() is called later, all defered functions are properly called.
func Run(payloadCmd []string) int {
	tmpWorkDir, err := ioutil.TempDir("", "qmstr-bin-")
	if err != nil {
		log.Fatalf("error creating temporary Hugo working directory: %v", err)
	}
	defer func() {
		if options.keepTmpDirectories {
			log.Printf("keeping temporary temporary at %v", tmpWorkDir)
		} else {
			log.Printf("deleting temporary temporary instrumentation bin directory in %v", tmpWorkDir)
			if err := os.RemoveAll(tmpWorkDir); err != nil {
				log.Printf("warning - error deleting temporary instrumentation bin directory in %v: %v", tmpWorkDir, err)
			}
		}
	}()
	SetupCompilerInstrumentation(tmpWorkDir)
	if len(payloadCmd) > 0 {
		exitCode, err := RunPayloadCommand(payloadCmd[0], payloadCmd[1:]...)
		if err != nil {
			log.Printf("payload command exited with non-zero exit code: %v", exitCode)
		}
		return exitCode
	}
	return 0
}

// SetupCompilerInstrumentation creates the QMSTR instrumentation symlinks in the given path
func SetupCompilerInstrumentation(tmpWorkDir string) {
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("")
	}
	ownPath, err := filepath.Abs(filepath.Dir(executable))
	if err != nil {
		log.Fatalf("unable to determine path to executable: %v", err)
	}
	const wrapper = "qmstr-wrapper"
	wrapperPath := path.Join(ownPath, wrapper)
	//TODO use exec.LookPath and use the wrapper wherever it is found
	if _, err := os.Stat(wrapperPath); err != nil {
		log.Fatalf("cannot find %s at %s: %v", wrapper, wrapperPath, err)
	}

	// create a "bin" directory in the temporary directory
	binDir := strings.TrimSpace(path.Join(tmpWorkDir, "bin"))
	if err := os.Mkdir(binDir, 0700); err != nil {
		log.Fatalf("unable to create %v: %v", binDir, err)
	}
	//create the symlinks to qmstr-wrapper in there

	//extend the PATH variable to include the created bin/ directory
	paths := filepath.SplitList(os.Getenv("PATH"))
	paths = append([]string{binDir}, paths...)
	separator := string(os.PathListSeparator)
	newPath := strings.Join(paths, separator)
	fmt.Printf("PATH is now %v\n", newPath)
}

// RunPayloadCommand performs the payload command and returns it's exit code and/or an error
func RunPayloadCommand(command string, arguments ...string) (int, error) {
	cmd := exec.Command(command, arguments...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	switch value := err.(type) {
	case *exec.ExitError:
		ws := value.Sys().(syscall.WaitStatus)
		return ws.ExitStatus(), fmt.Errorf("command finished with error: %v", err)
	default:
		return 0, nil
	}
}

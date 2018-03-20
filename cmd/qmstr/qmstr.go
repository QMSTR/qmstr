package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	golog "log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

// Options contains the context of a program invocation
type Options struct {
	progName           string // The name the program is called as
	keepTmpDirectories bool   // Keep intermediate files
	verbose            bool   // Enable trace log output
}

// global variables
var (
	options Options
	// Debug receives log messages in verbose mode
	Debug *golog.Logger
	// Log is the standard logger
	Log *golog.Logger
)

func main() {
	options.progName = os.Args[0]
	flag.BoolVar(&options.keepTmpDirectories, "keep", false, "Keep the created directories instead of cleaning up.")
	flag.BoolVar(&options.verbose, "verbose", false, "Enable diagnostic log output.")
	flag.Parse()
	if options.verbose {
		Debug = golog.New(os.Stderr, "DEBUG: ", golog.Ldate|golog.Ltime)
	} else {
		Debug = golog.New(ioutil.Discard, "", golog.Ldate|golog.Ltime)
	}
	Log = golog.New(os.Stderr, "", golog.Ldate|golog.Ltime)
	exitCode := Run(flag.Args())
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
	if len(payloadCmd) == 0 && !options.keepTmpDirectories {
		usage("No command specified!")
	}
	tmpWorkDir, err := ioutil.TempDir("", "qmstr-bin-")
	if err != nil {
		Log.Fatalf("error creating temporary Hugo working directory: %v", err)
	}
	defer func() {
		if options.keepTmpDirectories {
			Debug.Printf("keeping temporary temporary at %v", tmpWorkDir)
		} else {
			Debug.Printf("deleting temporary temporary instrumentation bin directory in %v", tmpWorkDir)
			if err := os.RemoveAll(tmpWorkDir); err != nil {
				// it is a warning because the program is exiting and we cannot recover anymore
				Log.Printf("warning - error deleting temporary instrumentation bin directory in %v: %v", tmpWorkDir, err)
			}
		}
	}()
	SetupCompilerInstrumentation(tmpWorkDir)
	if len(payloadCmd) > 0 {
		exitCode, err := RunPayloadCommand(payloadCmd[0], payloadCmd[1:]...)
		if err != nil {
			Debug.Printf("payload command exited with non-zero exit code: %v", exitCode)
		}
		return exitCode
	}
	return 0
}

// SetupCompilerInstrumentation creates the QMSTR instrumentation symlinks in the given path
func SetupCompilerInstrumentation(tmpWorkDir string) {
	executable, err := os.Executable()
	if err != nil {
		Log.Fatalf("unable to find myself: %v", err)
	}
	ownPath, err := filepath.Abs(filepath.Dir(executable))
	if err != nil {
		Log.Fatalf("unable to determine path to executable: %v", err)
	}
	const wrapper = "qmstr-wrapper"
	wrapperPath := path.Join(ownPath, wrapper)
	// TODO use exec.LookPath and use the wrapper wherever it is found
	if _, err := os.Stat(wrapperPath); err != nil {
		Log.Fatalf("cannot find %s at %s: %v", wrapper, wrapperPath, err)
	}

	// create a "bin" directory in the temporary directory
	binDir := strings.TrimSpace(path.Join(tmpWorkDir, "bin"))
	if err := os.Mkdir(binDir, 0700); err != nil {
		Log.Fatalf("unable to create %v: %v", binDir, err)
	}
	// create the symlinks to qmstr-wrapper in there
	symlinks := make(map[string]string)
	symlinks[path.Join(binDir, "gcc")] = wrapperPath
	for from, to := range symlinks {
		if err := os.Symlink(to, from); err != nil {
			Log.Fatalf("cannot symlink %s to %s: %v", from, to, err)
		}
	}
	// extend the PATH variable to include the created bin/ directory
	paths := filepath.SplitList(os.Getenv("PATH"))
	paths = append([]string{binDir}, paths...)
	separator := string(os.PathListSeparator)
	newPath := strings.Join(paths, separator)
	os.Setenv("PATH", newPath)
	Debug.Printf("PATH is now %v\n", os.Getenv("PATH"))
	if options.keepTmpDirectories {
		fmt.Printf("export PATH=%v\n", os.Getenv("PATH"))
	}
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

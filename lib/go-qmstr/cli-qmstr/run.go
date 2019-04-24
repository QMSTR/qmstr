package cliqmstr

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	keepTmpDirectories bool                                                  // Keep intermediate files
	instdir            string                                                // create instrumentation in this dir
	wrappedCmds        = []string{"gcc", "g++", "ar", "ld", "as", "objcopy"} // constant (which Go does not do with arrays)
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run a shell command with QMSTR instrumentation",
	Long: `run executes the argument after applying temporary QMSTR instrumentation.
For example, "qmstr make" sets up a QMSTR instrumented environment, then executes  "make", and finally restores the environment and deletes all temporary files.
The environment variable QMSTR_INSTRUMENTATION_HOME will be defined to point to the QMSTR instrumentation directory.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 && !keepTmpDirectories {
			return fmt.Errorf("no command specified")
		}
		return nil
	},
	Run: execute,
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().BoolVarP(&keepTmpDirectories, "keep", "k", false,
		"Keep the created directories instead of cleaning up.")
	runCmd.Flags().StringVarP(&instdir, "instdir", "i", "", "Create instrumentation in this directory (optional)")
}

func execute(cmd *cobra.Command, args []string) {
	exitCode := Run(args)
	os.Exit(exitCode)
}

// Run does everything
// It also makes sure that even though os.Exit() is called later, all defered functions are properly called.
func Run(payloadCmd []string) int {
	// Remind me why Go has no assertions?
	if len(payloadCmd) == 0 && !keepTmpDirectories {
		panic(fmt.Errorf("command validation constraint violated"))
	}

	var tmpWorkDir string
	var err error
	if instdir != "" {
		tmpWorkDir = instdir
		err = os.MkdirAll(tmpWorkDir, os.ModePerm)
	} else {
		tmpWorkDir, err = ioutil.TempDir("", "qmstr-bin-")
	}
	if err != nil {
		Log.Fatalf("error creating temporary working directory: %v", err)
	}

	defer func() {
		if keepTmpDirectories {
			Debug.Printf("keeping temporary directory at %v", tmpWorkDir)
		} else {
			Debug.Printf("deleting temporary instrumentation directory in %v", tmpWorkDir)
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
			fmt.Print(err)
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
	if _, err := os.Stat(wrapperPath); err != nil {
		Debug.Printf("cannot find %s at %s: %v (optional)", wrapper, wrapperPath, err)
		// optionally, search the path and use a qmstr-wrapper if found there
		wrapperPath, err = exec.LookPath(wrapper)
		if err != nil {
			Log.Fatalf("%v not found next in %v or in the PATH", wrapper, executable)
		}
	}

	// create a "bin" directory in the temporary directory
	binDir := strings.TrimSpace(filepath.Join(tmpWorkDir, "bin"))
	if err := link(wrapperPath, binDir, wrappedCmds); err != nil {
		Log.Fatalf("setting up instrumentation failed: %v", err)
	}

	// Setup environment variables
	envvars := make(map[string][]string)
	envvars["gcc"] = []string{"CMAKE_LINKER", "CC"}
	envvars["g++"] = []string{"CXX"}

	files, err := ioutil.ReadDir(binDir)
	for _, file := range files {
		cmd := file.Name()
		if envs, ok := envvars[cmd]; ok {
			for _, envvar := range envs {
				os.Setenv(envvar, filepath.Join(binDir, file.Name()))
			}
		}
	}

	// extend the PATH variable to include the created bin/ directory
	paths := filepath.SplitList(os.Getenv("PATH"))

	hasWhiteSpace, _ := regexp.Compile("\\s+")
	for index, value := range paths {
		if hasWhiteSpace.MatchString(value) {
			Debug.Printf("NOTE - your PATH contains an element with whitespace in it: %v", value)
			paths[index] = fmt.Sprintf("\"%s\"", value)
		}
	}

	paths = append([]string{binDir}, paths...)
	separator := string(os.PathListSeparator)
	newPath := strings.Join(paths, separator)
	os.Setenv("PATH", newPath)
	Debug.Printf("PATH is now %v\n", os.Getenv("PATH"))
	os.Setenv("QMSTR_INSTRUMENTATION_HOME", tmpWorkDir)
	Debug.Printf("QMSTR_INSTRUMENTATION_HOME is now %v\n", os.Getenv("QMSTR_INSTRUMENTATION_HOME"))
	if keepTmpDirectories {
		fmt.Printf("export PATH=%v\n", os.Getenv("PATH"))
		fmt.Printf("export QMSTR_INSTRUMENTATION_HOME=%v\n", os.Getenv("QMSTR_INSTRUMENTATION_HOME"))
	}
}

func link(source string, binDir string, targets []string) error {
	if err := os.Mkdir(binDir, 0700); err != nil {
		return fmt.Errorf("unable to create %v: %v", binDir, err)
	}

	// create the symlinks to qmstr-wrapper in there
	symlinks := make(map[string]string)
	for _, cmd := range wrappedCmds {
		symlink := path.Join(binDir, cmd)
		symlinks[symlink] = source
	}
	for from, to := range symlinks {
		if err := os.Symlink(to, from); err != nil {
			return fmt.Errorf("cannot symlink %s to %s: %v", from, to, err)
		}
	}
	return nil
}

// RunPayloadCommand performs the payload command and returns it's exit code and/or an error
func RunPayloadCommand(command string, arguments ...string) (int, error) {
	cmd := exec.Command(command, arguments...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		switch value := err.(type) {
		case *exec.ExitError:
			ws := value.Sys().(syscall.WaitStatus)
			return ws.ExitStatus(), err
		default:
			return 1, err
		}
	}
	return 0, nil
}

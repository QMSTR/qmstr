package main

import (
	"context"
	"fmt"
	"io/ioutil"
	golog "log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/docker"

	flag "github.com/spf13/pflag"
)

// Options contains the context of a program invocation
type Options struct {
	progName           string // The name the program is called as
	keepTmpDirectories bool   // Keep intermediate files
	verbose            bool   // Enable trace log output
	container          string // Image to spawn container from
	instdir            string // create instrumentation in this dir
}

// global variables
var (
	options Options
	// Debug receives log messages in verbose mode
	Debug *golog.Logger
	// Log is the standard logger
	Log         *golog.Logger
	wrappedCmds = []string{"gcc", "ar", "ld", "as"}
)

func main() {
	options.progName = os.Args[0]
	flag.BoolVar(&options.keepTmpDirectories, "keep", false, "Keep the created directories instead of cleaning up.")
	flag.BoolVar(&options.verbose, "verbose", false, "Enable diagnostic log output.")
	flag.StringVar(&options.container, "container", "", "Run command in a container from this image.")
	flag.StringVar(&options.instdir, "instdir", "", "Create instrumentation in this directory")
	flag.Parse()

	if options.verbose {
		Debug = golog.New(os.Stderr, "DEBUG: ", golog.Ldate|golog.Ltime)
	} else {
		Debug = golog.New(ioutil.Discard, "", golog.Ldate|golog.Ltime)
	}
	Log = golog.New(os.Stderr, "", golog.Ldate|golog.Ltime)

	if options.container != "" {
		if options.instdir == "" {
			options.instdir = "/tmp/qmstr-bin-container"
		}
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			Log.Fatalf("Failed to create docker client %v", err)
		}
		masterContainerID, intPort, err := docker.GetMasterInfo(ctx, cli)
		if err != nil {
			Log.Fatalf("Unable to find qmstr-master container")
		}

		var env []string
		if val, ok := os.LookupEnv(common.QMSTRDEBUGENV); ok {
			env = append(env, fmt.Sprintf("%s=%s", common.QMSTRDEBUGENV, val))
		}

		var mountpoints []mount.Mount
		if val, ok := os.LookupEnv(common.CCACHEDIRENV); ok {
			env = append(env, fmt.Sprintf("%s=%s", common.CCACHEDIRENV, common.ContainerCcacheDir))
			mountpoints = append(mountpoints, mount.Mount{Type: mount.TypeBind, Source: val, Target: common.ContainerCcacheDir})
		}

		Log.Printf("starting build container")
		err = docker.RunClientContainer(ctx, cli, &docker.ClientContainer{
			Image:             options.container,
			Cmd:               flag.Args(),
			MasterContainerID: masterContainerID,
			QmstrInternalPort: intPort,
			Instdir:           options.instdir,
			Env:               env,
			Mount:             mountpoints,
		})
		if err != nil {
			Log.Fatalf("Build container failed: %v", err)
		}
		os.Exit(0)
	}

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

	var tmpWorkDir string
	var err error
	if options.instdir != "" {
		tmpWorkDir = options.instdir
		err = os.MkdirAll(tmpWorkDir, os.ModePerm)
	} else {
		tmpWorkDir, err = ioutil.TempDir("", "qmstr-bin-")
	}
	if err != nil {
		Log.Fatalf("error creating temporary working directory: %v", err)
	}

	defer func() {
		if options.keepTmpDirectories {
			Debug.Printf("keeping temporary directory at %v", tmpWorkDir)
		} else {
			Debug.Printf("deleting temporary instrumentation bin directory in %v", tmpWorkDir)
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
		Debug.Printf("cannot find %s at %s: %v", wrapper, wrapperPath, err)
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
	if options.keepTmpDirectories {
		fmt.Printf("export PATH=%v\n", os.Getenv("PATH"))
		fmt.Printf("export QMSTR_INSTRUMENTATION_HOME=%v\n", os.Getenv("QMSTR_INSTRUMENTATION_HOME"))
	}
}

func link(source string, binDir string, targets []string) error {
	if err := os.Mkdir(binDir, 0700); err != nil {
		//Log.Fatalf("unable to create %v: %v", binDir, err)
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
			fmt.Errorf("cannot symlink %s to %s: %v", from, to, err)
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

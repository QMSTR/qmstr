package main

import (
	"context"
	"fmt"
	"io/ioutil"
	golog "log"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"

	flag "github.com/spf13/pflag"
)

// Options contains the context of a program invocation
type Options struct {
	progName           string // The name the program is called as
	keepTmpDirectories bool   // Keep intermediate files
	verbose            bool   // Enable trace log output
	container          string // Image to spawn container from
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
	fmt.Fprintln(os.Stderr, "WARNING: \"easy mode\" is work in progress! We *think* it works :-)")
	options.progName = os.Args[0]
	flag.BoolVar(&options.keepTmpDirectories, "keep", false, "Keep the created directories instead of cleaning up.")
	flag.BoolVar(&options.verbose, "verbose", false, "Enable diagnostic log output.")
	flag.StringVar(&options.container, "container", "", "Run command in a container from this image.")
	flag.Parse()

	if options.verbose {
		Debug = golog.New(os.Stderr, "DEBUG: ", golog.Ldate|golog.Ltime)
	} else {
		Debug = golog.New(ioutil.Discard, "", golog.Ldate|golog.Ltime)
	}
	Log = golog.New(os.Stderr, "", golog.Ldate|golog.Ltime)

	if options.container != "" {
		ctx := context.Background()
		cli, err := client.NewEnvClient()
		if err != nil {
			Log.Fatalf("Failed to create docker client %v", err)
		}
		masterContainerID, intPort, err := getMasterInfo(ctx, cli)
		if err != nil {
			Log.Fatalf("Unable to find qmstr-master container")
		}
		err = runContainer(ctx, cli, options.container, flag.Args(), masterContainerID, intPort)
		if err != nil {
			Log.Fatalf("Build container failed: %v", err)
		}
		os.Exit(0)
	}

	exitCode := Run(flag.Args())
	os.Exit(exitCode)
}

func getMasterInfo(ctx context.Context, cli *client.Client) (string, uint16, error) {
	qmstrAddr := os.Getenv("QMSTR_MASTER")
	if qmstrAddr == "" {
		return "", 0, errors.New("QMSTR_MASTER not set, can't determine qmstr-master container")
	}
	qmstrAddrS := strings.Split(qmstrAddr, ":")
	qmstrHostPort, err := strconv.ParseUint(qmstrAddrS[len(qmstrAddrS)-1], 10, 64)
	if err != nil {
		return "", 0, err
	}

	args, err := filters.ParseFlag("label=org.qmstr.image=master", filters.NewArgs())
	if err != nil {
		return "", 0, err
	}
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{Filters: args})
	if err != nil {
		return "", 0, err
	}

	for _, container := range containers {
		for _, portCfg := range container.Ports {
			if uint64(portCfg.PublicPort) == qmstrHostPort {
				return container.ID, portCfg.PrivatePort, nil
			}
		}
	}

	return "", 0, errors.New("qmstr-master container not found")
}

func runContainer(ctx context.Context, cli *client.Client, image string, cmd []string, masterContainerID string, qmstrInternalPort uint16) error {
	Debug.Printf("Using master container %s", masterContainerID)

	const containerBuildDir = "/buildroot"
	wd, err := os.Getwd()
	if err != nil {
		Log.Fatalf("unable to determine current working directory")
	}
	hostConf := &container.HostConfig{
		Mounts:      []mount.Mount{mount.Mount{Source: wd, Target: containerBuildDir, Type: mount.TypeBind}},
		NetworkMode: container.NetworkMode(fmt.Sprintf("container:%s", masterContainerID)),
	}

	containerConf := &container.Config{
		Image: image,
		Cmd:   append([]string{"qmstr", "--"}, cmd...),
		Tty:   true,
		Env:   []string{fmt.Sprintf("QMSTR_MASTER=%s:%d", masterContainerID[:12], qmstrInternalPort)},
	}

	user, err := user.Current()
	if err == nil {
		containerConf.User = user.Uid
	}

	resp, err := cli.ContainerCreate(ctx, containerConf, hostConf, nil, "")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	status, err := cli.ContainerWait(ctx, resp.ID)
	if err != nil {
		return err
	}

	Debug.Printf("Build container returned status %d", status)

	out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}

	logmsg, err := ioutil.ReadAll(out)
	if err != nil {
		return err
	}
	Log.Printf("Container logs:\n%s", logmsg)

	if status != 0 {
		os.Exit(int(status))
	}

	return nil
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
	binDir := strings.TrimSpace(path.Join(tmpWorkDir, "bin"))
	if err := os.Mkdir(binDir, 0700); err != nil {
		Log.Fatalf("unable to create %v: %v", binDir, err)
	}

	envvars := make(map[string][]string)
	envvars["gcc"] = []string{"CMAKE_LINKER", "CC"}

	// create the symlinks to qmstr-wrapper in there
	wrappedCmds := []string{"gcc"}
	symlinks := make(map[string]string)
	for _, cmd := range wrappedCmds {
		symlink := path.Join(binDir, cmd)
		symlinks[symlink] = wrapperPath
		if envs, ok := envvars[cmd]; ok {
			for _, envvar := range envs {
				os.Setenv(envvar, symlink)
			}
		}
	}
	for from, to := range symlinks {
		if err := os.Symlink(to, from); err != nil {
			Log.Fatalf("cannot symlink %s to %s: %v", from, to, err)
		}
	}
	// extend the PATH variable to include the created bin/ directory
	paths := filepath.SplitList(os.Getenv("PATH"))

	hasWhiteSpace, _ := regexp.Compile("\\s+")
	for index, value := range paths {
		if hasWhiteSpace.MatchString(value) {
			Debug.Printf("NOTE - your PATH contains a element with whitespace in it: %v", value)
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

// RunPayloadCommand performs the payload command and returns it's exit code and/or an error
func RunPayloadCommand(command string, arguments ...string) (int, error) {
	cmd := exec.Command(command, arguments...)
	cmd.Stdin = os.Stdin
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

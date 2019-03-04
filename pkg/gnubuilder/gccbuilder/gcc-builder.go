package gccbuilder

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"

	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/gnubuilder"
	pb "github.com/QMSTR/qmstr/pkg/service"
	"github.com/spf13/pflag"
)

const undef = "undef"

type GccBuilder struct {
	Builder    string
	Mode       gnubuilder.Mode
	Input      []string
	Output     []string
	WorkDir    string
	LinkLibs   []string
	LibPath    []string
	SysLibPath []string
	Args       []string
	staticLink bool
	StaticLibs map[string]struct{}
	ActualLibs map[string]string
	Ccache     string
	builder.GeneralBuilder
}

func NewGccBuilder(workDir string, logger *log.Logger, debug bool) *GccBuilder {
	var ccacheCmd *exec.Cmd
	// setup ccache if possible
	ccachePath := common.FindExecutablesOnPath("ccache")
	if len(ccachePath) > 0 {
		ccacheCmd = exec.Command(ccachePath[0], "-s")
		ccacheCmd.Start()
	}

	builder := GccBuilder{
		Builder:        "",
		Mode:           gnubuilder.ModeLink,
		Input:          []string{},
		Output:         []string{},
		WorkDir:        workDir,
		LinkLibs:       []string{},
		LibPath:        []string{},
		Args:           []string{},
		SysLibPath:     gnubuilder.GetSysLibPath(),
		staticLink:     false,
		StaticLibs:     map[string]struct{}{},
		ActualLibs:     map[string]string{},
		GeneralBuilder: builder.GeneralBuilder{Logger: logger, Debug: debug, Afs: afero.NewOsFs()},
	}

	if ccacheCmd != nil {
		if err := ccacheCmd.Wait(); err != nil {
			builder.Ccache = ccacheCmd.Path
		}
	}

	return &builder
}

func (g *GccBuilder) GetPrefix() (string, error) {
	if g.Ccache != "" {
		return g.Ccache, nil
	}
	return "", errors.New("Ccache not found")
}

func (g *GccBuilder) GetName() string {
	switch g.Builder {
	case "gcc":
		return "GNU C compiler builder"
	case "g++":
		return "GNU C++ compiler builder"
	default:
		return "unknown C/C++ compiler builder"
	}
}

func (g *GccBuilder) Setup() error {
	return os.Setenv(common.QMSTRWRAPGCC, "")
}

func (g *GccBuilder) TearDown() error {
	// Unset environment variable before we end
	return os.Unsetenv(common.QMSTRWRAPGCC)
}

func (g *GccBuilder) Analyze(commandline []string) ([]*pb.FileNode, error) {

	if g.Debug {
		g.Logger.Printf("Parsing commandline %v", commandline)
	}
	if err := g.parseCommandLine(commandline[1:]); err != nil {
		return nil, fmt.Errorf("Failed to parse commandline: %v", err)
	}

	g.Builder = commandline[0]

	switch g.Mode {
	case gnubuilder.ModeLink:
		if g.staticLink {
			g.Logger.Printf("%s linking statically", g.Builder)
		} else {
			g.Logger.Printf("%s linking", g.Builder)
		}
		fileNodes := []*pb.FileNode{}
		linkedTarget := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, g.Output[0], false), pb.FileNode_TARGET)
		libraries := []*pb.FileNode{}
		dependencies := []*pb.FileNode{}
		for _, inFile := range g.Input {
			inputFileNode := &pb.FileNode{}
			ext := filepath.Ext(inFile)
			if ext == ".o" {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), pb.FileNode_INTERMEDIATE)
				libraries = append(libraries, inputFileNode)
			} else if ext == ".c" || ext == ".cc" || ext == ".cpp" || ext == ".c++" || ext == ".cp" || ext == ".cxx" {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), pb.FileNode_SOURCE)
				libraries = append(libraries, inputFileNode)
			} else if strings.HasSuffix(inFile, ".so") {
				depFileNode := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), pb.FileNode_TARGET)
				dependencies = append(dependencies, depFileNode)
			} else {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), pb.FileNode_TARGET)
				libraries = append(libraries, inputFileNode)
			}
		}
		err := gnubuilder.FindActualLibraries(g.Afs, g.ActualLibs, g.LinkLibs, append(g.LibPath, g.SysLibPath...), g.staticLink, g.StaticLibs)
		if err != nil {
			g.Logger.Fatalf("Failed to collect libraries: %v", err)
		}
		for _, actualLib := range g.ActualLibs {
			if strings.HasSuffix(actualLib, ".so") {
				runtimeDep := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, actualLib, false), pb.FileNode_TARGET)
				dependencies = append(dependencies, runtimeDep)
			} else {
				linkLib := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, actualLib, false), pb.FileNode_TARGET)
				libraries = append(libraries, linkLib)
			}
		}
		linkedTarget.DerivedFrom = libraries
		linkedTarget.Dependencies = dependencies
		fileNodes = append(fileNodes, linkedTarget)
		return fileNodes, nil
	case gnubuilder.ModeAssemble:
		g.Logger.Printf("%s assembling - skipping link", g.Builder)
		fileNodes := []*pb.FileNode{}
		if g.Debug {
			g.Logger.Printf("This is our input %v", g.Input)
			g.Logger.Printf("This is our output %v", g.Output)
		}
		for idx, inFile := range g.Input {
			if g.Debug {
				g.Logger.Printf("This is the source file %s indexed %d", inFile, idx)
			}
			sourceFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), pb.FileNode_SOURCE)
			targetFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, g.Output[idx], false), pb.FileNode_INTERMEDIATE)
			targetFile.DerivedFrom = []*pb.FileNode{sourceFile}
			fileNodes = append(fileNodes, targetFile)
		}
		return fileNodes, nil
	case gnubuilder.ModeCompile:
		g.Logger.Printf("%s compile - skipping assemble and link", g.Builder)
		fileNodes := []*pb.FileNode{}
		if g.Debug {
			g.Logger.Printf("This is our input %v", g.Input)
			g.Logger.Printf("This is our output %v", g.Output)
		}
		for idx, inFile := range g.Input {
			if g.Debug {
				g.Logger.Printf("This is the source file %s indexed %d", inFile, idx)
			}
			sourceFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), pb.FileNode_SOURCE)
			targetFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, g.Output[idx], false), pb.FileNode_SOURCE)
			targetFile.DerivedFrom = []*pb.FileNode{sourceFile}
			fileNodes = append(fileNodes, targetFile)
		}
		return fileNodes, nil
	case gnubuilder.ModePrintOnly:
		g.Logger.Println("print only; nothing produced")
		return nil, nil
	default:
		return nil, builder.ErrBuilderModeNotImplemented
	}
}

func (g *GccBuilder) parseCommandLine(args []string) error {
	if g.Debug {
		g.Logger.Printf("Parsing arguments: %v", args)
	}

	// remove all flags we don't care about but that would break parsing
	g.Args = gnubuilder.CleanCmdLine(args, g.Logger, g.Debug, g.staticLink, g.StaticLibs, gnubuilder.ModeUndef)

	gccFlags := pflag.NewFlagSet("gcc", pflag.ContinueOnError)
	gccFlags.BoolP("assemble", "c", false, "do not link")
	gccFlags.BoolP("compile", "S", false, "do not assemble")
	gccFlags.BoolP("preprocess", "E", false, "do not compile")
	gccFlags.StringP("output", "o", undef, "output")
	gccFlags.StringSliceP("includepath", "I", []string{}, "include path")
	gccFlags.String("isystem", undef, "system include path")
	gccFlags.String("include", undef, "include header file")
	gccFlags.StringSliceVarP(&g.LinkLibs, "linklib", "l", []string{}, "link libs")
	gccFlags.StringSliceVarP(&g.LibPath, "linklibdir", "L", []string{}, "search dir for link libs")

	if g.Debug {
		g.Logger.Printf("Parsing cleaned commandline: %v", g.Args)
	}
	err := gccFlags.Parse(g.Args)
	if err != nil {
		return fmt.Errorf("Unrecoverable commandline parsing error: %v", err)
	}

	g.Input = gccFlags.Args()

	if ok, err := gccFlags.GetBool("assemble"); ok && err == nil {
		g.Mode = gnubuilder.ModeAssemble
	}
	if ok, err := gccFlags.GetBool("compile"); ok && err == nil {
		g.Mode = gnubuilder.ModeCompile
	}
	if ok, err := gccFlags.GetBool("preprocess"); ok && err == nil {
		g.Mode = gnubuilder.ModePreproc
	}
	if g.Debug {
		g.Logger.Printf("Mode set to: %v", g.Mode)
	}

	if output, err := gccFlags.GetString("output"); err == nil && output != undef {
		g.Output = []string{output}
	} else {
		// no output defined
		switch g.Mode {
		case gnubuilder.ModeLink:
			if len(g.Input) == 0 {
				// No input no output
				g.Mode = gnubuilder.ModePrintOnly
				return nil
			}
			g.Output = []string{"a.out"}
		case gnubuilder.ModeAssemble:
			for _, input := range g.Input {
				objectname := strings.TrimSuffix(input, filepath.Ext(input)) + ".o"
				g.Output = append(g.Output, objectname)
			}
		case gnubuilder.ModeCompile:
			for _, input := range g.Input {
				objectname := strings.TrimSuffix(input, filepath.Ext(input)) + ".s"
				g.Output = append(g.Output, objectname)
			}
		}
	}
	return nil
}

func (g *GccBuilder) GetPushFile() (*pb.PushFileMessage, error) {
	// handle piped code
	if g.Input[0] == "-" && g.StdinChannel != nil {
		g.Logger.Println("Reading data from stdin channel")
		data := <-g.StdinChannel
		g.Logger.Printf("Read data %s from stdin channel", data)
		checksum, err := common.Hash(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("failed to capture code from standard input: %v", err)
		}
		return &pb.PushFileMessage{Hash: checksum, Name: checksum, Data: data}, nil
	}
	return g.GeneralBuilder.GetPushFile()
}

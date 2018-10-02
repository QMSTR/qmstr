package gccbuilder

import (
	"errors"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/gnubuilder"
	pb "github.com/QMSTR/qmstr/pkg/qmstr/service"
	"github.com/spf13/pflag"
)

type mode int

const (
	Link mode = iota
	Preproc
	Compile
	Assemble
	PrintOnly
)
const undef = "undef"
const (
	linkedTrg = "linkedtarget"
	obj       = "objectfile"
	src       = "sourcecode"
)

type GccBuilder struct {
	Mode       mode
	Input      []string
	Output     []string
	WorkDir    string
	LinkLibs   []string
	LibPath    []string
	Args       []string
	staticLink bool
	StaticLibs map[string]struct{}
	ActualLibs map[string]string
	builder.GeneralBuilder
}

func NewGccBuilder(workDir string, logger *log.Logger, debug bool) *GccBuilder {
	return &GccBuilder{Link, []string{}, []string{}, workDir, []string{}, []string{}, []string{}, false, map[string]struct{}{}, map[string]string{}, builder.GeneralBuilder{logger, debug}}
}

func (g *GccBuilder) GetPrefix() (string, error) {
	// setup ccache if possible
	ccachePath := common.FindExecutablesOnPath("ccache")
	if len(ccachePath) > 0 {
		return ccachePath[0], nil
	}
	return "", errors.New("Ccache not found")
}

func (g *GccBuilder) GetName() string {
	return "GNU C compiler builder"
}

func (g *GccBuilder) Analyze(commandline []string) (*pb.BuildMessage, error) {
	if g.Debug {
		g.Logger.Printf("Parsing commandline %v", commandline)
	}
	if err := g.parseCommandLine(commandline[1:]); err != nil {
		return nil, fmt.Errorf("Failed to parse commandline: %v", err)
	}

	switch g.Mode {
	case Link:
		if g.staticLink {
			g.Logger.Printf("gcc linking statically")
		} else {
			g.Logger.Printf("gcc linking")
		}
		fileNodes := []*pb.FileNode{}
		linkedTarget := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, g.Output[0], false), linkedTrg)
		dependencies := []*pb.FileNode{}
		for _, inFile := range g.Input {
			inputFileNode := &pb.FileNode{}
			ext := filepath.Ext(inFile)
			if ext == ".o" {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), obj)
			} else if ext == ".c" {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), src)
			} else {
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), linkedTrg)
			}
			dependencies = append(dependencies, inputFileNode)
		}
		err := gnubuilder.FindActualLibraries(g.ActualLibs, g.LinkLibs, g.LibPath, g.staticLink, g.StaticLibs)
		if err != nil {
			g.Logger.Fatalf("Failed to collect dependencies: %v", err)
		}
		for _, actualLib := range g.ActualLibs {
			linkLib := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, actualLib, false), linkedTrg)
			dependencies = append(dependencies, linkLib)
		}
		linkedTarget.DerivedFrom = dependencies
		fileNodes = append(fileNodes, linkedTarget)
		return &pb.BuildMessage{FileNodes: fileNodes}, nil
	case Assemble:
		g.Logger.Printf("gcc assembling - skipping link")
		fileNodes := []*pb.FileNode{}
		if g.Debug {
			g.Logger.Printf("This is our input %v", g.Input)
			g.Logger.Printf("This is our output %v", g.Output)
		}
		for idx, inFile := range g.Input {
			if g.Debug {
				g.Logger.Printf("This is the source file %s indexed %d", inFile, idx)
			}
			sourceFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), src)
			targetFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, g.Output[idx], false), obj)
			targetFile.DerivedFrom = []*pb.FileNode{sourceFile}
			fileNodes = append(fileNodes, targetFile)
		}
		return &pb.BuildMessage{FileNodes: fileNodes}, nil
	case Compile:
		g.Logger.Printf("gcc compile - skipping assemble and link")
		fileNodes := []*pb.FileNode{}
		if g.Debug {
			g.Logger.Printf("This is our input %v", g.Input)
			g.Logger.Printf("This is our output %v", g.Output)
		}
		for idx, inFile := range g.Input {
			if g.Debug {
				g.Logger.Printf("This is the source file %s indexed %d", inFile, idx)
			}
			sourceFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, inFile, false), src)
			targetFile := builder.NewFileNode(common.BuildCleanPath(g.WorkDir, g.Output[idx], false), src)
			targetFile.DerivedFrom = []*pb.FileNode{sourceFile}
			fileNodes = append(fileNodes, targetFile)
		}
		return &pb.BuildMessage{FileNodes: fileNodes}, nil
	case PrintOnly:
		log.Println("print only; nothing produced")
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
	g.Args = gnubuilder.CleanCmdLine(args, g.Logger, g.Debug, g.staticLink, g.StaticLibs, undef)

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
		g.Mode = Assemble
	}
	if ok, err := gccFlags.GetBool("compile"); ok && err == nil {
		g.Mode = Compile
	}
	if ok, err := gccFlags.GetBool("preprocess"); ok && err == nil {
		g.Mode = Preproc
	}
	if g.Debug {
		g.Logger.Printf("Mode set to: %v", g.Mode)
	}

	if output, err := gccFlags.GetString("output"); err == nil && output != undef {
		g.Output = []string{output}
	} else {
		// no output defined
		switch g.Mode {
		case Link:
			if len(g.Input) == 0 {
				// No input no output
				g.Mode = PrintOnly
				return nil
			}
			g.Output = []string{"a.out"}
		case Assemble:
			for _, input := range g.Input {
				objectname := strings.TrimSuffix(input, filepath.Ext(input)) + ".o"
				g.Output = append(g.Output, objectname)
			}
		case Compile:
			for _, input := range g.Input {
				objectname := strings.TrimSuffix(input, filepath.Ext(input)) + ".s"
				g.Output = append(g.Output, objectname)
			}
		}
	}
	return nil
}

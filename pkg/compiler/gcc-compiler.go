package compiler

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/wrapper"
	"github.com/spf13/pflag"
)

type mode int

const (
	Link mode = iota
	Preproc
	Compile
	Assemble
)
const undef = "undef"

type GccCompiler struct {
	Mode    mode
	Input   []string
	Output  []string
	WorkDir string
	logger  *log.Logger
}

func NewGccCompiler(workDir string, logger *log.Logger) *GccCompiler {
	return &GccCompiler{Link, []string{}, []string{}, workDir, logger}
}

func (g *GccCompiler) Analyze(commandline []string) (*pb.BuildMessage, error) {
	g.logger.Printf("Parsing commandline %v", commandline)
	g.parseCommandLine(commandline[1:])

	switch g.Mode {
	case Link:
		g.logger.Printf("gcc linking")
		buildLinkMsg := pb.BuildMessage_Link{Target: &pb.File{Path: g.Output[0]}}
		dependencies := []*pb.File{}
		for _, inFile := range g.Input {
			inputFile := pb.File{Path: wrapper.BuildCleanPath(g.WorkDir, inFile)}
			dependencies = append(dependencies, &inputFile)
		}
		buildLinkMsg.Dependencies = dependencies
		buildMsg := pb.BuildMessage{}
		buildMsg.Binary = []*pb.BuildMessage_Link{&buildLinkMsg}

		return &buildMsg, nil
	case Assemble:
		g.logger.Printf("gcc assembling - skipping link")
		buildMsg := pb.BuildMessage{}
		buildMsg.Compilations = []*pb.BuildMessage_Compile{}

		for idx, inFile := range g.Input {
			sourceFile := pb.File{Path: wrapper.BuildCleanPath(g.WorkDir, inFile)}
			targetFile := pb.File{Path: wrapper.BuildCleanPath(g.WorkDir, g.Output[idx])}
			buildMsg.Compilations = append(buildMsg.Compilations, &pb.BuildMessage_Compile{Source: &sourceFile, Target: &targetFile})
		}
		return &buildMsg, nil
	default:
		return nil, errors.New("Mode not implemented")
	}
}

func (g *GccCompiler) parseCommandLine(args []string) {
	gccFlags := pflag.NewFlagSet("gcc", pflag.ContinueOnError)
	gccFlags.BoolP("assemble", "c", false, "do not link")
	gccFlags.BoolP("compile", "S", false, "do not assemble")
	gccFlags.BoolP("preprocess", "E", false, "do not compile")
	gccFlags.StringP("output", "o", undef, "output")
	gccFlags.StringP("language", "x", undef, "language")
	gccFlags.Parse(args)

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

	if output, err := gccFlags.GetString("output"); err == nil && output != undef {
		g.Output = []string{output}
	} else {
		// no output defined
		switch g.Mode {
		case Link:
			g.Output = []string{"a.out"}
		case Assemble:
			for _, input := range g.Input {
				objectname := strings.TrimSuffix(input, filepath.Ext(input)) + ".o"
				g.Output = append(g.Output, objectname)
			}
		}
	}
}

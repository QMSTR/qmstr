package ldbuilder

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/lib/go-qmstr/builder"
	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/gnubuilder"
	"github.com/spf13/pflag"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

const undef = "undef"

type LdBuilder struct {
	Input       []string
	Output      []string
	WorkDir     string
	LinkLibs    []string
	LibPath     []string
	SysLibsPath []string
	Args        []string
	ActualLibs  map[string]string
	staticLink  bool
	StaticLibs  map[string]struct{}
	Mode        gnubuilder.Mode
	builder.GeneralBuilder
}

func NewLdBuilder(workDir string, logger *log.Logger, debug bool) *LdBuilder {
	return &LdBuilder{
		Input:          []string{},
		Output:         []string{},
		WorkDir:        workDir,
		LinkLibs:       []string{},
		LibPath:        []string{},
		SysLibsPath:    gnubuilder.GetSysLibPath(),
		Args:           []string{},
		ActualLibs:     map[string]string{},
		staticLink:     false,
		StaticLibs:     map[string]struct{}{},
		Mode:           gnubuilder.ModeLink,
		GeneralBuilder: builder.NewGeneralBuilder(logger, debug)}
}

func (ld *LdBuilder) GetPrefix() (string, error) {
	return "", errors.New("ld not prefixed")
}

func (ld *LdBuilder) GetName() string {
	return "GNU ld linker"
}

func (ld *LdBuilder) Analyze(commandline []string) ([]*service.FileNode, error) {
	// skip when wrapping gcc
	if _, gccCalled := os.LookupEnv(common.QMSTRWRAPGCC); gccCalled {
		return []*service.FileNode{}, nil
	}

	if err := ld.parseCommandLine(commandline[1:]); err != nil {
		return nil, fmt.Errorf("Failed to parse commandline: %v", err)
	}

	switch ld.Mode {
	case gnubuilder.ModePrintOnly:
		ld.Logger.Println("print only; nothing produced")
		return nil, nil
	case gnubuilder.ModePreproc, gnubuilder.ModeCompile, gnubuilder.ModeAssemble, gnubuilder.ModeUndef:
		return nil, builder.ErrBuilderModeNotSupported
	}

	if ld.staticLink {
		ld.Logger.Printf("ld linking statically")
	} else {
		ld.Logger.Printf("ld linking")
	}
	fileNodes := []*service.FileNode{}
	linkedTarget := builder.NewFileNode(common.BuildCleanPath(ld.WorkDir, ld.Output[0], false), false)
	libraries := []*service.FileNode{}
	dependencies := []*service.FileNode{}
	for _, inFile := range ld.Input {
		ext := filepath.Ext(inFile)
		inputFileNode := builder.NewFileNode(common.BuildCleanPath(ld.WorkDir, inFile, false), true)
		if _, ok := builder.LibExtension[ext]; ok {
			dependencies = append(dependencies, inputFileNode)
		} else {
			libraries = append(libraries, inputFileNode)
		}
	}
	err := gnubuilder.FindActualLibraries(ld.Afs, ld.ActualLibs, ld.LinkLibs, append(ld.LibPath, ld.SysLibsPath...), ld.staticLink, ld.StaticLibs)
	if err != nil {
		ld.Logger.Fatalf("Failed to collect libraries: %v", err)
	}
	for _, actualLib := range ld.ActualLibs {
		ext := filepath.Ext(actualLib)
		dep := builder.NewFileNode(common.BuildCleanPath(ld.WorkDir, actualLib, false), true)
		if _, ok := builder.LibExtension[ext]; ok {
			dependencies = append(dependencies, dep)
		} else {
			libraries = append(libraries, dep)
		}
	}
	linkedTarget.DerivedFrom = libraries
	linkedTarget.Dependencies = dependencies
	fileNodes = append(fileNodes, linkedTarget)
	return fileNodes, nil
}

func (ld *LdBuilder) parseCommandLine(args []string) error {
	if ld.Debug {
		ld.Logger.Printf("Parsing arguments: %v", args)
	}

	// remove all flags we don't care about but that would break parsing
	ld.Args = gnubuilder.CleanCmdLine(args, ld.Logger, ld.Debug, ld.staticLink, ld.StaticLibs, gnubuilder.ModeLink)

	ldFlags := pflag.NewFlagSet("ld", pflag.ContinueOnError)
	ldFlags.StringP("output", "o", undef, "output")
	ldFlags.StringSliceVarP(&ld.LinkLibs, "linklib", "l", []string{}, "link libs")
	ldFlags.StringSliceVarP(&ld.LibPath, "linklibdir", "L", []string{}, "search dir for link libs")

	if ld.Debug {
		ld.Logger.Printf("Parsing cleaned commandline: %v", ld.Args)
	}
	err := ldFlags.Parse(ld.Args)
	if err != nil {
		return fmt.Errorf("Unrecoverable commandline parsing error: %v", err)
	}

	ld.Input = ldFlags.Args()

	if output, err := ldFlags.GetString("output"); err == nil && output != undef {
		ld.Output = []string{output}
	} else {
		// no output defined
		if len(ld.Input) == 0 {
			// No input no output
			ld.Mode = gnubuilder.ModePrintOnly
		} else {
			ld.Output = []string{"a.out"}
		}
	}
	return nil
}

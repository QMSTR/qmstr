package asbuilder

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"

	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/gnubuilder"
	"github.com/QMSTR/qmstr/pkg/service"
)

const undef = "undef"

type AsBuilder struct {
	Input   string
	Output  string
	WorkDir string
	Args    []string
	builder.GeneralBuilder
}

func NewAsBuilder(workDir string, logger *log.Logger, debug bool) *AsBuilder {
	return &AsBuilder{"", "", workDir, []string{}, builder.NewGeneralBuilder(logger, debug)}
}

func (as *AsBuilder) GetPrefix() (string, error) {
	return "", errors.New("as not prefixed")
}

func (as *AsBuilder) GetName() string {
	return "GNU as assember"
}

func (as *AsBuilder) Analyze(commandline []string) ([]*service.FileNode, error) {
	// skip when wrapping gcc
	if _, gccCalled := os.LookupEnv(common.QMSTRWRAPGCC); gccCalled {
		return []*service.FileNode{}, nil
	}

	if err := as.parseCommandLine(commandline[1:]); err != nil {
		return nil, fmt.Errorf("Failed to parse commandline: %v", err)
	}

	as.Logger.Printf("as assembling")
	fileNodes := []*service.FileNode{}
	if as.Debug {
		as.Logger.Printf("This is our input %v", as.Input)
		as.Logger.Printf("This is our output %v", as.Output)
	}
	if as.Debug {
		as.Logger.Printf("This is the source file %s", as.Input)
	}
	sourceFile := builder.NewFileNode(common.BuildCleanPath(as.WorkDir, as.Input, false), service.FileNode_SOURCE)
	targetFile := builder.NewFileNode(common.BuildCleanPath(as.WorkDir, as.Output, false), service.FileNode_INTERMEDIATE)
	targetFile.DerivedFrom = []*service.FileNode{sourceFile}
	fileNodes = append(fileNodes, targetFile)

	return fileNodes, nil
}

func (as *AsBuilder) parseCommandLine(args []string) error {
	if as.Debug {
		as.Logger.Printf("Parsing arguments: %v", args)
	}

	// remove all flags we don't care about but that would break parsing
	as.Args = gnubuilder.CleanCmdLine(args, as.Logger, as.Debug, false, map[string]struct{}{}, gnubuilder.ModeAssemble)

	asFlags := pflag.NewFlagSet("as", pflag.ContinueOnError)
	asFlags.StringP("output", "o", undef, "output")
	asFlags.StringSliceP("includepath", "I", []string{}, "include path")

	if as.Debug {
		as.Logger.Printf("Parsing cleaned commandline: %v", as.Args)
	}

	err := asFlags.Parse(as.Args)
	if err != nil {
		return fmt.Errorf("Unrecoverable commandline parsing error: %s", err)
	}

	args = asFlags.Args()
	if len(args) != 1 {
		return fmt.Errorf("Commandline should have just 1 input file. Inputs: %v", args)
	}
	as.Input = args[0]

	if output, err := asFlags.GetString("output"); err == nil && output != undef {
		as.Output = output
	} else {
		objectname := strings.TrimSuffix(as.Input, filepath.Ext(as.Input)) + ".o"
		as.Output = objectname
	}
	return nil
}

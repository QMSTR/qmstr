package asbuilder

import (
	"errors"
	"log"
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"

	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/gnubuilder"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

const mode = "Assemble"
const undef = "undef"

const (
	obj = "objectfile"
	src = "sourcecode"
)

type AsBuilder struct {
	Input   []string
	Output  []string
	WorkDir string
	Args    []string
	builder.GeneralBuilder
}

func NewAsBuilder(workDir string, logger *log.Logger, debug bool) *AsBuilder {
	return &AsBuilder{[]string{}, []string{}, workDir, []string{}, builder.GeneralBuilder{logger, debug}}
}

func (as *AsBuilder) GetPrefix() (string, error) {
	return "", errors.New("as not prefixed")
}

func (as *AsBuilder) GetName() string {
	return "GNU as assember"
}

func (as *AsBuilder) Analyze(commandline []string) (*service.BuildMessage, error) {
	as.parseCommandLine(commandline[1:])

	as.Logger.Printf("as assembling")
	fileNodes := []*service.FileNode{}
	if as.Debug {
		as.Logger.Printf("This is our input %v", as.Input)
		as.Logger.Printf("This is our output %v", as.Output)
	}
	for idx, inFile := range as.Input {
		if as.Debug {
			as.Logger.Printf("This is the source file %s indexed %d", inFile, idx)
		}
		sourceFile := builder.NewFileNode(common.BuildCleanPath(as.WorkDir, inFile, false), src)
		targetFile := builder.NewFileNode(common.BuildCleanPath(as.WorkDir, as.Output[idx], false), obj)
		targetFile.DerivedFrom = []*service.FileNode{sourceFile}
		fileNodes = append(fileNodes, targetFile)
	}
	return &service.BuildMessage{FileNodes: fileNodes}, nil
}

func (as *AsBuilder) parseCommandLine(args []string) {
	if as.Debug {
		as.Logger.Printf("Parsing arguments: %v", args)
	}

	// remove all flags we don't care about but that would break parsing
	as.Args = gnubuilder.CleanCmdLine(args, as.Logger, as.Debug, false, map[string]struct{}{}, mode)

	asFlags := pflag.NewFlagSet("as", pflag.ContinueOnError)
	asFlags.StringP("output", "o", undef, "output")

	if as.Debug {
		as.Logger.Printf("Parsing cleaned commandline: %v", as.Args)
	}
	err := asFlags.Parse(as.Args)
	if err != nil {
		as.Logger.Fatalf("Unrecoverable commandline parsing error: %s", err)
	}

	as.Input = asFlags.Args()

	if output, err := asFlags.GetString("output"); err == nil && output != undef {
		as.Output = []string{output}
	} else {
		for _, input := range as.Input {
			objectname := strings.TrimSuffix(input, filepath.Ext(input)) + ".o"
			as.Output = append(as.Output, objectname)
		}
	}
}

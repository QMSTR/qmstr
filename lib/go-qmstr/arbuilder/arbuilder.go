package arbuilder

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/QMSTR/qmstr/lib/go-qmstr/builder"
	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

type command int

const (
	Undef command = iota
	Delete
	Move
	Print
	QuickAppend
	Replace
	Index
	Display
	Extract
)

var arCmdModPattern = regexp.MustCompile("-??([dmpqrstx]{1})([abcDfilMNoPsSTuvV]*)")
var relposCountModPattern = regexp.MustCompile("[abN]{1}")

type ArBuilder struct {
	Command         command
	Modifiers       string
	CommandLineArgs []string
	Input           []string
	Output          string
	WorkDir         string
	builder.GeneralBuilder
}

func NewArBuilder(workDir string, logger *log.Logger, debug bool) *ArBuilder {
	return &ArBuilder{Undef, "", nil, []string{}, "", workDir, builder.NewGeneralBuilder(logger, debug)}
}

func (a *ArBuilder) GetPrefix() (string, error) {
	return "", errors.New("ar not prefixed")
}

func (a *ArBuilder) Analyze(commandline []string) ([]*service.FileNode, error) {
	if len(commandline) < 3 {
		return nil, fmt.Errorf("failed to analyze \"%s\" too few arguments", commandline)
	}

	commandline, err := a.processFlags(commandline)
	if err != nil {
		if err.Error() == "noop" {
			os.Exit(0)
		}
		return nil, err
	}

	cmdMod := arCmdModPattern.FindStringSubmatch(commandline[1])

	switch cmdMod[1] {
	case "r":
		a.Command = Replace
	case "d":
		a.Command = Delete
	case "m":
		a.Command = Move
	case "p":
		a.Command = Print
	case "q":
		a.Command = QuickAppend
	case "s":
		a.Command = Index
	case "t":
		a.Command = Display
	case "x":
		a.Command = Extract
	}

	a.CommandLineArgs = commandline[2:]

	if len(cmdMod) == 3 && cmdMod[2] != "" {
		a.Modifiers = cmdMod[2]
	}

	a.processModifiers()

	a.Output = a.CommandLineArgs[0]
	a.Input = a.CommandLineArgs[1:]

	msg, err := a.getResults()
	if err != nil {
		return nil, fmt.Errorf("Failed to generate result message: %v", err)
	}

	return msg, nil
}

func (a *ArBuilder) getResults() ([]*service.FileNode, error) {
	if a.Command == Replace || a.Command == QuickAppend {
		a.Logger.Printf("archiving")
		fileNodes := []*service.FileNode{}
		linkedTarget := builder.NewFileNode(common.BuildCleanPath(a.WorkDir, a.Output, false), service.FileNode_TARGET)
		dependencies := []*service.FileNode{}
		for _, inFile := range a.Input {
			inputFileNode := &service.FileNode{}
			ext := filepath.Ext(inFile)
			switch ext {
			case ".o":
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(a.WorkDir, inFile, false), service.FileNode_INTERMEDIATE)
			case ".a":
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(a.WorkDir, inFile, false), service.FileNode_TARGET)
			default:
				inputFileNode = builder.NewFileNode(common.BuildCleanPath(a.WorkDir, inFile, false), service.FileNode_INTERMEDIATE)
			}
			dependencies = append(dependencies, inputFileNode)
		}
		linkedTarget.DerivedFrom = dependencies
		fileNodes = append(fileNodes, linkedTarget)
		return fileNodes, nil
	}
	return nil, errors.New("Command not supported")
}

func (a *ArBuilder) processFlags(commandline []string) ([]string, error) {
	cleanIdx := []int{}
	for idx, arg := range commandline {
		switch arg {
		case "--help", "--version":
			return nil, errors.New("noop")
		case "--target", "--plugin":
			cleanIdx = append(cleanIdx, idx, idx+1)
			continue
		case "-X32_64":
			cleanIdx = append(cleanIdx, idx)
			continue
		}
		if strings.HasPrefix(arg, "@") {
			return nil, errors.New("Reading commandline options from file is not supported")
		}
		if strings.HasPrefix(arg, "--target=") || strings.HasPrefix(arg, "--plugin=") {
			cleanIdx = append(cleanIdx, idx)
			continue
		}
	}

	commandline = builder.CleanCmd(commandline, cleanIdx, a.Debug, a.Logger)

	return commandline, nil
}

func (a *ArBuilder) processModifiers() error {
	if relposCountModPattern.MatchString(a.Modifiers) {
		// remove relpos or count argument
		a.CommandLineArgs = a.CommandLineArgs[1:]
	}
	return nil
}

func (a *ArBuilder) GetName() string {
	return "GNU ar builder"
}

package stripbuilder

import (
	"errors"
	"log"
	"strings"

	"github.com/QMSTR/qmstr/lib/go-qmstr/builder"
	"github.com/QMSTR/qmstr/lib/go-qmstr/common"
	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

type StripBuilder struct {
	Workdir string
	Input   []string
	Output  []string
	builder.GeneralBuilder
}

func NewStripBuilder(workDir string, logger *log.Logger, debug bool) *StripBuilder {
	return &StripBuilder{workDir, []string{}, []string{}, builder.NewGeneralBuilder(logger, debug)}
}

func (s *StripBuilder) GetPrefix() (string, error) {
	return "", errors.New("strip not prefixed")
}

func (s *StripBuilder) GetName() string {
	return "strip builder"
}

func (s *StripBuilder) Analyze(commandline []string) ([]*service.FileNode, error) {
	if s.Debug {
		s.Logger.Printf("%s parsing commandline %v", s.GetName(), commandline)
	}

	err := s.processFlags(commandline[1:])
	if err != nil {
		return nil, err
	}

	fileNodes := []*service.FileNode{}
	for _, input := range s.Input {
		inputTarget := builder.NewFileNode(common.BuildCleanPath(s.Workdir, input, false), true)
		outputTarget := builder.NewFileNode(common.BuildCleanPath(s.Workdir, input, false), false)
		outputTarget.DerivedFrom = []*service.FileNode{inputTarget}
		fileNodes = append(fileNodes, outputTarget)
		s.Logger.Printf("Striping %s:%s", inputTarget.GetPath(), inputTarget.FileData.GetHash())
	}

	return fileNodes, nil
}

func (s *StripBuilder) processFlags(args []string) error {
	if s.Debug {
		s.Logger.Printf("Parsing arguments: %v", args)
	}

	cleanIdx := []int{}
	for idx, arg := range args {
		switch arg {
		case "--remove-section":
			cleanIdx = append(cleanIdx, idx, idx+1)
			continue
		case "-g", "-S", "-d", "--strip-debug", "-s", "--strip-unneeded", "--enable-deterministic-archives", "-D":
			cleanIdx = append(cleanIdx, idx)
			continue
		}
		if strings.HasPrefix(arg, "--remove-section=") {
			cleanIdx = append(cleanIdx, idx)
		}
	}
	s.Input = builder.CleanCmd(args, cleanIdx, s.Debug, s.Logger)

	if len(s.Input) <= 0 {
		return builder.ErrNoTargetsProvided
	}
	s.Output = s.Input

	return nil
}

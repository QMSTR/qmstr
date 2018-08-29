package arbuilder

import (
	"fmt"
	"log"

	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/service"
)

type ArBuilder struct {
	Input   []string
	Output  string
	WorkDir string
	builder.GeneralBuilder
}

func NewArBuilder(workDir string, logger *log.Logger, debug bool) *ArBuilder {
	return &ArBuilder{[]string{}, "", workDir, builder.GeneralBuilder{logger, debug}}
}

func (g *ArBuilder) Analyze(commandline []string) (*service.BuildMessage, error) {
	return nil, fmt.Errorf("failed to analyze \"%s\"", commandline)
}

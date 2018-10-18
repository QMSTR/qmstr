package builder

import (
	"log"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

var (
	ErrBuilderModeNotImplemented = errors.New("Mode not implemented")
	ErrNoTargetsProvided         = errors.New("No targets provided")
)

type Builder interface {
	Analyze(commandline []string) (*service.BuildMessage, error)
	GetPushFile() (*service.PushFileMessage, error)
	GetName() string
	GetPrefix() (string, error)
	SetStdinChannel(chan []byte)
}

type GeneralBuilder struct {
	Logger       *log.Logger
	Debug        bool
	Afs          afero.Fs
	StdinChannel chan []byte
}

func NewGeneralBuilder(logger *log.Logger, debug bool) GeneralBuilder {
	return GeneralBuilder{Logger: logger, Debug: debug, Afs: afero.NewOsFs()}
}

func (gb *GeneralBuilder) SetStdinChannel(stdin chan []byte) {
	gb.StdinChannel = stdin
}

func (gb *GeneralBuilder) GetPushFile() (*service.PushFileMessage, error) {
	return nil, errors.New("No file to push")
}

func NewFileNode(path string, fileType string) *service.FileNode {
	filename := filepath.Base(path)
	hash, err := common.HashFile(path)
	broken := false
	if err != nil {
		hash = "nohash" + path
		broken = true
	}
	return &service.FileNode{Name: filename, Type: fileType, Path: path, Hash: hash, Broken: broken}
}

func CleanCmd(commandline []string, cleanIdx []int, debug bool, logger *log.Logger) []string {
	for i, ci := range cleanIdx {
		realIdx := ci - i
		if debug {
			logger.Printf("Clearing argument: %v", commandline[realIdx])
		}
		if realIdx == len(commandline)-1 {
			commandline = commandline[:realIdx-1]
			break
		}
		rest := commandline[realIdx+1:]
		commandline = append(commandline[:realIdx], rest...)
		if debug {
			logger.Printf("new slice is %v", commandline)
		}
	}
	return commandline
}

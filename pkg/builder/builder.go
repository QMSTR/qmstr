package builder

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

var (
	ErrBuilderModeNotImplemented = errors.New("Mode not implemented")
	ErrNoTargetsProvided         = errors.New("No targets provided")
)

type Builder interface {
	Analyze(commandline []string) (*service.BuildMessage, error)
	GetName() string
	GetPrefix() (string, error)
}

type GeneralBuilder struct {
	Logger *log.Logger
	Debug  bool
	Afs    afero.Fs
}

func NewGeneralBuilder(logger *log.Logger, debug bool) GeneralBuilder {
	return GeneralBuilder{Logger: logger, Debug: debug, Afs: afero.NewOsFs()}
}

func NewFileNode(path string, fileType string) *service.FileNode {
	filename := filepath.Base(path)
	hash, err := hash(path)
	broken := false
	if err != nil {
		hash = "nohash" + path
		broken = true
	}
	return &service.FileNode{Name: filename, Type: fileType, Path: path, Hash: hash, Broken: broken}
}

func hash(fileName string) (string, error) {
	h := sha1.New()
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := f.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return "", err
		}
		h.Write(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
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

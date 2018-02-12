package compiler

import (
	"log"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/wrapper"
)

type Compiler interface {
	Analyze(commandline []string) (*pb.BuildMessage, error)
}

type GeneralCompiler struct {
	logger *log.Logger
	debug  bool
}

func GetCompiler(prog string, workDir string, logger *log.Logger, debug bool) Compiler {
	switch prog {
	case "gcc", "g++":
		return NewGccCompiler(workDir, logger, debug)
	default:
		log.Printf("Compiler %s not available", prog)
	}

	return nil
}

func NewFile(path string) (*pb.File, error) {
	hash, err := wrapper.Hash(path)
	if err != nil {
		return &pb.File{Hash: "nohash" + path, Path: path}, nil
	}
	return &pb.File{Hash: hash, Path: path}, nil
}

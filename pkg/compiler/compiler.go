package compiler

import (
	"log"
	"path/filepath"

	pb "github.com/QMSTR/qmstr/pkg/service"
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

func NewFileNode(path string, fileType string) *pb.FileNode {
	filename := filepath.Base(path)
	hash, err := wrapper.Hash(path)
	broken := false
	if err != nil {
		hash = "nohash" + path
		broken = true
	}
	return &pb.FileNode{NodeType: pb.NodeTypeFileNode, Name: filename, Type: fileType, Path: path, Hash: hash, Broken: broken}
}

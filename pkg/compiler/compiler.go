package compiler

import (
	"log"

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
	hash, err := wrapper.Hash(path)
	if err != nil {
		return &pb.FileNode{Path: path, Hash: "nohash" + path, Broken: true}
	}
	return &pb.FileNode{Type: fileType, Path: path, Hash: hash, Broken: false}
}

func NewFileNodeDerivedFrom(path string, fileType string, derivedNode []*pb.FileNode) *pb.FileNode {
	hash, err := wrapper.Hash(path)
	if err != nil {
		return &pb.FileNode{Path: path, Hash: "nohash" + path, Broken: true, DerivedFrom: derivedNode}
	}
	return &pb.FileNode{Type: fileType, Path: path, Hash: hash, Broken: false, DerivedFrom: derivedNode}
}

package compiler

import (
	"log"

	pb "github.com/QMSTR/qmstr/pkg/buildservice"
)

type Compiler interface {
	Analyze(commandline []string) (*pb.BuildMessage, error)
}

func GetCompiler(prog string, workDir string, logger *log.Logger) Compiler {
	switch prog {
	case "gcc", "g++":
		return NewGccCompiler(workDir, logger)
	default:
		log.Printf("Compiler %s not available", prog)
	}

	return nil
}

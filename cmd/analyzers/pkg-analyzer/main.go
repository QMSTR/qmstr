package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/service"
)

const queryType = "linkedtarget"

type PkgAnalyzer struct {
	targetsSlice []string
	targetsDir   string
}

func main() {
	analyzer := analysis.NewAnalyzer(&PkgAnalyzer{})
	if err := analyzer.RunAnalyzerModule(); err != nil {
		log.Printf("%v failed: %v\n", analyzer.GetModuleName(), err)
		os.Exit(master.ReturnAnalyzerFailed)
	}
}

func (pkganalyzer *PkgAnalyzer) Configure(configMap map[string]string) error {
	if _, ok := configMap["targets"]; !ok {
		log.Println("No linked targets found to be linked to package node.\n Please provide linked targets in the configuration.")
		return errors.New("Misconfigured package analyzer")
	}
	pkganalyzer.targetsSlice = strings.Split(configMap["targets"], ";")

	if _, ok := configMap["targetdir"]; !ok {
		log.Println("No target directories found. Please provide target directories in the configuration.")
		return errors.New("Misconfigured package analyzer")
	}
	pkganalyzer.targetsDir = configMap["targetdir"]

	return nil
}

// Analyze finds the targets in db which we are going to connect to the package node
func (pkganalyzer *PkgAnalyzer) Analyze(controlService service.ControlServiceClient, session string) error {
	queryNode := &service.FileNode{Type: queryType}

	pkgNodeResp, err := controlService.GetPackageNode(context.Background(), &service.PackageRequest{Session: session})
	stream, err := controlService.GetFileNode(context.Background(), queryNode)
	if err != nil {
		log.Printf("Could not get file node %v", err)
		return err
	}

	for {
		fileNode, err := stream.Recv()
		if err == io.EOF {
			break
		}

		for _, target := range pkganalyzer.targetsSlice {
			if fileNode.Path == filepath.Join(pkganalyzer.targetsDir, target) {
				log.Printf("Adding node %v to package targets.", fileNode.Path)
				pkgNodeResp.PackageNode.Targets = append(pkgNodeResp.PackageNode.Targets, fileNode)
			}
		}
	}
	return nil
}

func (pkganalyzer *PkgAnalyzer) PostAnalyze() error {
	return nil
}

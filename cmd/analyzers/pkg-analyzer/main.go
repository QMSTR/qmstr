package main

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/service"
)

type PkgAnalyzer struct {
	targetsSlice []string
	pkgNode      *service.PackageNode
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

func (pkganalyzer *PkgAnalyzer) Analyze(node *service.FileNode) (*service.InfoNodeSlice, error) {
	for _, target := range pkganalyzer.targetsSlice {
		if node.Path == filepath.Join(pkganalyzer.targetsDir, target) {
			log.Printf("Adding node %v to package targets.", node.Path)
			pkganalyzer.pkgNode.Targets = append(pkganalyzer.pkgNode.Targets, node)
		}
	}
	return &service.InfoNodeSlice{}, nil
}

func (pkganalyzer *PkgAnalyzer) PostAnalyze() error {
	return nil
}

func (pkganalyzer *PkgAnalyzer) GetPackageNode() *service.PackageNode {
	return pkganalyzer.pkgNode
}

func (pkganalyzer *PkgAnalyzer) SetPackageNode(pkgNode *service.PackageNode) {
	pkganalyzer.pkgNode = pkgNode
}

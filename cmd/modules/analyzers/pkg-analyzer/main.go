package main

import (
	"context"
	"errors"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/common"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

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
func (pkganalyzer *PkgAnalyzer) Analyze(controlService service.ControlServiceClient, analysisService service.AnalysisServiceClient, token int64) error {
	pkgNode, err := controlService.GetPackageNode(context.Background(), &service.PackageRequest{})
	if err != nil {
		return err
	}

	PackageNodeMsgs := []*service.PackageNodeMessage{}
	queryNode := &service.FileNode{FileType: service.FileNode_TARGET}

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
			re := regexp.MustCompile(filepath.Join(pkganalyzer.targetsDir, target))
			if re.MatchString(fileNode.Path) {
				hash, err := common.HashFile(fileNode.Path)
				if err != nil {
					return err
				}
				if hash == fileNode.Hash {
					pkgNode.Targets = append(pkgNode.Targets, fileNode)
					break
				}
			}
		}
	}

	PackageNodeMsgs = append(PackageNodeMsgs, &service.PackageNodeMessage{Token: token, Packagenode: pkgNode})
	sendStream, err := analysisService.SendPackageNode(context.Background())
	if err != nil {
		return err
	}
	for _, fnodeMsg := range PackageNodeMsgs {
		sendStream.Send(fnodeMsg)
	}

	reply, err := sendStream.CloseAndRecv()
	if err != nil {
		return err
	}
	if reply.Success {
		log.Println("Package Analyzer sent FileNodes")
	}
	return nil
}

func (pkganalyzer *PkgAnalyzer) PostAnalyze() error {
	return nil
}

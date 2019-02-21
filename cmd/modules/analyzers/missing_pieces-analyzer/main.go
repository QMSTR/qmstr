package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/QMSTR/qmstr/pkg/analysis"
	"github.com/QMSTR/qmstr/pkg/master"
	"github.com/QMSTR/qmstr/pkg/service"
	yaml "gopkg.in/yaml.v2"
)

type MissingPiecesAnalyzer struct {
	File []*service.FileNode
}

func main() {
	analyzer := analysis.NewAnalyzer(&MissingPiecesAnalyzer{})
	if err := analyzer.RunAnalyzerModule(); err != nil {
		log.Printf("%v failed: %v\n", analyzer.GetModuleName(), err)
		os.Exit(master.ReturnAnalyzerFailed)
	}
}

func (mpanalyzer *MissingPiecesAnalyzer) Configure(configMap map[string]string) error {
	if inputfile, ok := configMap["inputfile"]; ok {
		if _, err := os.Stat(inputfile); os.IsNotExist(err) {
			return fmt.Errorf("File %s not found", inputfile)
		}

		f, err := os.Open(inputfile)
		if err != nil {
			return err
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		if err := mpanalyzer.readData(data); err != nil {
			return fmt.Errorf("Failed to read config from %s: %v", inputfile, err)
		}
		return nil
	}
	return fmt.Errorf("Misconfigured Missing Pieces Analyzer. No input file declared")
}

func (mpanalyzer *MissingPiecesAnalyzer) Analyze(controlService service.ControlServiceClient, analysisService service.AnalysisServiceClient, token int64) error {
	sendStream, err := analysisService.SendFileNode(context.Background())
	if err != nil {
		return err
	}
	for _, fnode := range mpanalyzer.File {
		log.Printf("send filenode for %s", fnode.Path)
		sendStream.Send(&service.FileNodeMessage{Filenode: fnode, Token: token})
	}

	reply, err := sendStream.CloseAndRecv()
	if err != nil {
		return err
	}
	if reply.Success {
		log.Println("Missing Pieces Analyzer sent FileNodes")
	}
	return nil
}

func (mpanalyzer *MissingPiecesAnalyzer) PostAnalyze() error {
	return nil
}

func (mpanalyzer *MissingPiecesAnalyzer) readData(data []byte) error {
	err := yaml.Unmarshal(data, mpanalyzer)
	if err != nil {
		return err
	}
	if mpanalyzer.File == nil {
		return fmt.Errorf("no data found -- check indentation")
	}
	return nil
}

package analysis

import (
	"log"
	"strings"

	"github.com/QMSTR/qmstr/pkg/config"
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type Analyzer interface {
	Analyze(node *AnalysisNode) error
}

type Analysis struct {
	Name     string
	Nodes    []AnalysisNode
	Analyzer Analyzer
}

type AnalysisNode struct {
	actualNode service.FileNode
	pathSub    []*config.PathSubstitution
	db         *database.DataBase
	dirty      bool
}

func NewAnalysisNode(actualNode service.FileNode, pathSub []*config.PathSubstitution, db *database.DataBase) AnalysisNode {
	return AnalysisNode{actualNode: actualNode, pathSub: pathSub, db: db, dirty: false}
}

func (an *AnalysisNode) GetPath() string {
	actualPath := an.actualNode.Path
	for _, pathsubmsg := range an.pathSub {
		actualPath = strings.Replace(actualPath, pathsubmsg.Old, pathsubmsg.New, 1)
	}
	return actualPath
}

func (an *AnalysisNode) GetName() string {
	return an.actualNode.Name
}

func (an *AnalysisNode) Store() error {
	if an.dirty {
		_, err := an.db.AlterNode(&an.actualNode)
		if err != nil {
			return err
		}
	}
	return nil
}

func RunAnalysis(analysis Analysis) {
	log.Printf("Starting analysis: %s", analysis.Name)
	for _, node := range analysis.Nodes {
		err := analysis.Analyzer.Analyze(&node)
		if err != nil {
			log.Printf("Analysis of %s failed: %v\n", node.GetPath(), err)
			panic("Analysis corrupt")
		}
		err = node.Store()
		if err != nil {
			log.Printf("Storing failed: %v\n", err)
			panic("Analysis corrupt")
		}
	}
	log.Printf("Analysis finished")
}

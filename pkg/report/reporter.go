package report

import (
	"github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/database"
)

type ReportNode struct {
	actualNode database.Node
	db         *database.DataBase
}

type Reporter interface {
	Generate(node []ReportNode) (*buildservice.ReportResponse, error)
}

func NewReportNode(actualNode database.Node, db *database.DataBase) ReportNode {
	return ReportNode{actualNode: actualNode, db: db}
}

func (an *ReportNode) GetPath() string {
	return an.actualNode.Path
}

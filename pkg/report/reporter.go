package report

import (
	"github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/database"
)

type Reporter interface {
	Generate(node []*database.Node) (*buildservice.ReportResponse, error)
}

package report

import (
	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

type Reporter interface {
	Generate(node []*database.Node) (*service.ReportResponse, error)
}

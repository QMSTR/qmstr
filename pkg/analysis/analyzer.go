package analysis

import "github.com/QMSTR/qmstr/pkg/database"

type Analyzer interface {
	Analyze(node database.Node) (error, interface{})
}

package database

import (
	"sync/atomic"

	"github.com/QMSTR/qmstr/pkg/service"
)

// AddProjectNode adds a node to the insert queue
func (db *DataBase) AddProjectNode(node *service.ProjectNode) {
	atomic.AddUint64(&db.pending, 1)
	db.insertQueue <- node
}

package database

import (
	"errors"
	"sync/atomic"

	"github.com/QMSTR/qmstr/pkg/service"
)

var ErrNoProjectNode = errors.New("No project node found")

// AddProjectNode adds a node to the insert queue
func (db *DataBase) AddProjectNode(node *service.ProjectNode) {
	atomic.AddUint64(&db.pending, 1)
	db.insertQueue <- node
}

func (db *DataBase) GetProjectNode() (*service.ProjectNode, error) {
	var ret map[string][]*service.ProjectNode

	q := `{
		getProjectNode(func: has(projectNodeType)) {
			uid
			name
		}
	}`

	err := db.queryNodesSimple(q, &ret)
	if err != nil {
		return nil, err
	}

	projNodes := ret["getProjectNode"]
	if len(projNodes) < 1 {
		return nil, ErrNoProjectNode
	}

	return projNodes[0], nil
}

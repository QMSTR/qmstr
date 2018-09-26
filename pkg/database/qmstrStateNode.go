package database

import (
	"errors"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

func (db *DataBase) AddQmstrStateNode(qNode *service.QmstrStateNode) (string, error) {
	return dbInsert(db.client, qNode)
}

func (db *DataBase) GetQmstrStateNode() (*service.QmstrStateNode, error) {
	var ret map[string][]*service.QmstrStateNode

	q := `query QmstrStateNode() {
		getQmstrStateNode(func: has(qmstrStateNodeType)) @recurse(loop: false) {
			uid
			session
			phase
		  }}`

	err := db.queryNodesSimple(q, &ret)
	if err != nil {
		return nil, err
	}

	stateNodes := ret["getQmstrStateNode"]
	if len(stateNodes) < 1 {
		return nil, errors.New("No package node found")
	}

	return stateNodes[0], nil

}

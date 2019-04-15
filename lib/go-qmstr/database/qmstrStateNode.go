package database

import (
	"errors"

	"github.com/QMSTR/qmstr/lib/go-qmstr/service"
)

func (db *DataBase) AddQmstrStateNode(qNode *service.QmstrStateNode) (string, error) {
	qmstrState, err := db.GetQmstrStateNode()
	if err == nil {
		qNode.Uid = qmstrState.Uid
	}
	return dbInsert(db.client, qNode)
}

func (db *DataBase) GetQmstrStateNode() (*service.QmstrStateNode, error) {
	var ret map[string][]*service.QmstrStateNode

	q := `{
		getQmstrStateNode(func: has(qmstrStateNodeType)) @recurse(loop: false) {
			uid
			phase
		}
	}`

	err := db.queryNodesSimple(q, &ret)
	if err != nil {
		return nil, err
	}

	stateNodes := ret["getQmstrStateNode"]
	if len(stateNodes) < 1 {
		return nil, errors.New("No qmstr state nod found")
	}

	return stateNodes[0], nil

}

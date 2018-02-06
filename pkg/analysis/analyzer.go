package analysis

import (
	"strings"

	"github.com/QMSTR/qmstr/pkg/buildservice"
	"github.com/QMSTR/qmstr/pkg/database"
)

type Analyzer interface {
	Analyze(node *AnalysisNode) error
}

type AnalysisNode struct {
	actualNode *database.Node
	pathSub    []*buildservice.PathSubstitutionMessage
	db         *database.DataBase
	dirty      bool
}

func NewAnalysisNode(actualNode *database.Node, pathSub []*buildservice.PathSubstitutionMessage, db *database.DataBase) *AnalysisNode {
	return &AnalysisNode{actualNode: actualNode, pathSub: pathSub, db: db, dirty: false}
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

func (an *AnalysisNode) SetLicense(spdxLicenseIdentifier string) error {
	uid, err := an.db.GetLicenseUid(spdxLicenseIdentifier)
	if err != nil {
		return err
	}
	an.actualNode.License = database.License{Uid: uid}
	an.dirty = true
	return nil
}

func (an *AnalysisNode) Store() error {
	if an.dirty {
		_, err := an.db.AlterNode(an.actualNode)
		if err != nil {
			return err
		}
	}
	return nil
}

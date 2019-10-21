package database

import "github.com/QMSTR/qmstr/lib/go-qmstr/service"

func (db *DataBase) GetAnalyzerByName(name string) (*service.Analyzer, error) {
	var ret map[string][]*service.Analyzer

	q := `query AnalyzerByName($AnaName: string){
		  getAnalyzerByType(func: has(analyzerNodeType)) @filter(eq(name, $AnaName)) {
			uid
			name
		  }}`

	vars := map[string]string{"$AnaName": name}

	err := db.queryNodes(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	if len(ret["getAnalyzerByName"]) < 1 {
		// No such analyzer
		analyzer := &service.Analyzer{Uid: "_:analyzer", Name: name}
		uid, err := dbInsert(db.client, analyzer)
		if err != nil {
			return nil, err
		}
		analyzer.Uid = uid["analyzer"]
		return analyzer, nil
	}

	return ret["getAnalyzerByName"][0], nil
}

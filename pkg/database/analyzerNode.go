package database

import "github.com/QMSTR/qmstr/pkg/service"

func (db *DataBase) GetAnalyzerByName(name string) (*service.Analyzer, error) {
	ret := map[string][]*service.Analyzer{}

	q := `query AnalyzerByName($AnaName: string){
		  getAnalyzerByType(func: has(analyzerNodeType)) @filter(eq(name, $AnaName)) {
			uid
			hash
			path
			derivedFrom
		  }}`

	vars := map[string]string{"$AnaName": name}

	err := db.queryAnalyzer(q, vars, &ret)
	if err != nil {
		return nil, err
	}

	if len(ret["getAnalyzerByName"]) < 1 {
		// No such analyzer
		analyzer := &service.Analyzer{Name: name}
		uid, err := dbInsert(db.client, analyzer)
		if err != nil {
			return nil, err
		}
		analyzer.Uid = uid
		return analyzer, nil
	}

	return ret["getAnalyzerByName"][0], nil
}

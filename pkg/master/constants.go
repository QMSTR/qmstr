package master

const (
	ReturnSuccess int = iota
	ReturnAnalyzerFailed
	ReturnAnalysisServiceFailed
	ReturnAnalysisServiceCommFailed
	ReturnReporterFailed
	ReturnReportServiceCommFailed
	ReturnReportServiceFailed
)

const (
	PhaseIDInit int32 = iota
	PhaseIDBuild
	PhaseIDAnalysis
	PhaseIDReport
	PhaseIDFailed
)

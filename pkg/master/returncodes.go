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

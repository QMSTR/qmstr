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

const ServerCacheDir string = "/var/cache/qmstr"
const ServerOutputDir string = "/var/qmstr"

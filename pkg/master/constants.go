package master

const (
	phaseIDInit int32 = iota
	phaseIDBuild
	phaseIDAnalysis
	phaseIDReport
	phaseIDFailed
)

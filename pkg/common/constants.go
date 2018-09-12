package common

//Environment variable names
const (
	// QMSTRADDRENV is the name of the environment variable that holds connection string to access qmstr-master server
	QMSTRADDRENV = "QMSTR_MASTER"
	// QMSTRDEBUGENV is the name of the environment variable that defines if qmstr runs in debug mode
	QMSTRDEBUGENV = "QMSTR_DEBUG"
)

// ContainerBuildDir is where the source code gets mounted in the qmstr-master container
const ContainerBuildDir = "/buildroot"

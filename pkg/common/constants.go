package common

//Environment variable names
const (
	// QMSTRADDRENV is the name of the environment variable that holds connection string to access qmstr-master server
	QMSTRADDRENV = "QMSTR_MASTER"
	// QMSTRDEBUGENV is the name of the environment variable that defines if qmstr runs in debug mode
	QMSTRDEBUGENV = "QMSTR_DEBUG"
	// CCACHEDIRENV is the name of the environment variable that stores the path to the ccache cache directory
	CCACHEDIRENV = "CCACHE_DIR"
)

// ContainerBuildDir is where the source code gets mounted in the qmstr-master container
const ContainerBuildDir = "/buildroot"

// ContainerCcacheDir is where the cccache dir gets mountet to
const ContainerCcacheDir = "/ccache"

// ContainerQmstrHomeDir is the HOME dir of the user running a client container
const ContainerQmstrHomeDir = "/home/qmstruser"

const ContainerGraphExportDir = "/var/qmstr/export"
const ContainerGraphImportPath = "/var/qmstr/qmstr.import.tar"

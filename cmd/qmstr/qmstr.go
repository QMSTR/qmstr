package main

import (
	"log"
	"os"
	"path"
	"path/filepath"
)

func main() {
	SetupCompilerInstrumentation()
}

// SetupCompilerInstrumentation creates the QMSTR instrumentation symlinks in the given path
func SetupCompilerInstrumentation() {
	executable, err := os.Executable()
	if err != nil {
		log.Fatalf("")
	}
	ownPath, err := filepath.Abs(filepath.Dir(executable))
	if err != nil {
		log.Fatalf("unable to determine path to executable: %v", err)
	}
	const wrapper = "qmstr-wrapper"
	wrapperPath := path.Join(ownPath, wrapper)
	if _, err := os.Stat(wrapperPath); err != nil {
		log.Fatalf("cannot find %s at %s: %v", wrapper, wrapperPath, err)
	}

	//create a "bin" directory in the temporary directory
	//create the symlinks to qmstr-wrapper in there
	//extend the PATH variable to include the created bin/ directory
}

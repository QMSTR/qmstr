package logging

import (
	"io/ioutil"
	golog "log"
	"os"
)

// Logging holds a configured Quartermaster logger
type Logging struct {
	// Debug receives log messages in verbose mode
	Debug *golog.Logger
	// Log is the standard logger
	Log *golog.Logger
}

// Setup sets up logging
func Setup(verbose bool) Logging {
	var log Logging

	log.Log = golog.New(os.Stderr, "", golog.Ldate|golog.Ltime)
	if verbose {
		log.Debug = golog.New(os.Stderr, "DEBUG: ", golog.Ldate|golog.Ltime)
	} else {
		log.Debug = golog.New(ioutil.Discard, "", golog.Ldate|golog.Ltime)
	}

	return log
}

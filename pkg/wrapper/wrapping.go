package wrapper

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"

	"github.com/QMSTR/qmstr/pkg/arbuilder"
	"github.com/QMSTR/qmstr/pkg/builder"
	"github.com/QMSTR/qmstr/pkg/gnubuilder/gccbuilder"
	"github.com/QMSTR/qmstr/pkg/gnubuilder/ldbuilder"
)

// Wrapper represents a wrapper to call a program
type Wrapper struct {
	logger          *log.Logger
	Program         string
	commandlineArgs []string
	debug           bool
	Builder         builder.Builder
}

// NewWrapper returns an instance of a Wrapper for the given command line
func NewWrapper(commandline []string, workdir string, logger *log.Logger, debug bool) (*Wrapper, error) {
	// extract the compiler that was supposed to run
	w := Wrapper{}
	w.logger = logger
	w.Program = filepath.Base(commandline[0])
	w.debug = debug
	//extract the arguments
	w.commandlineArgs = commandline[1:]
	b, err := getBuilder(w.Program, workdir, logger, debug)
	if err != nil {
		return nil, err
	}
	w.Builder = b
	return &w, nil
}

func getBuilder(prog string, workDir string, logger *log.Logger, debug bool) (builder.Builder, error) {
	switch prog {
	case "gcc", "g++":
		return gccbuilder.NewGccBuilder(workDir, logger, debug), nil
	case "ar":
		return arbuilder.NewArBuilder(workDir, logger, debug), nil
	case "ld":
		return ldbuilder.NewLdBuilder(workDir, logger, debug), nil
	}
	return nil, fmt.Errorf("Builder %s not available", prog)
}

// Wrap calls the actual program to be wrapped and preserves output and return value
func (w *Wrapper) Wrap() {

	if w.Program == "qmstr-wrapper" {
		log.Fatal("This is not how you should invoke the qmstr-wrapper.\n\tSee https://github.com/QMSTR/qmstr for more information on how to use the QMSTR.")
	}

	// find and run actual program
	actualProg, err := FindActualProgram(w.Program)
	if err != nil {
		log.Fatalf("actual executable was not found: %v", err)
	}

	// setup next compiler wrapper
	if prefix, err := w.Builder.GetPrefix(); err == nil {
		w.logger.Printf("Using chained compiler wrapper %s", prefix)
		w.commandlineArgs = append([]string{actualProg}, w.commandlineArgs...)
		actualProg = prefix
	}

	cmd := exec.Command(actualProg, w.commandlineArgs...)
	var stdoutbuf, stderrbuf bytes.Buffer
	cmd.Stdout = &stdoutbuf
	cmd.Stderr = &stderrbuf

	// connect stdin to pass piped data through and save for analysis
	stdin, err := cmd.StdinPipe()
	if err != nil {
		w.logger.Panic(err)
	}

	stdinChannel := make(chan []byte, 1024)
	stdinHandler := func(stdin io.WriteCloser, c chan []byte) {
		defer stdin.Close()
		tee := io.TeeReader(os.Stdin, stdin)
		r := bufio.NewReader(tee)
		nBytes, nChunks := int64(0), int64(0)
		buf := make([]byte, 0, 1024)
		for {
			if w.debug {
				w.logger.Println("Reading data from stdin")
			}
			n, err := r.Read(buf[:cap(buf)])
			buf = buf[:n]
			if n == 0 {
				if err == nil {
					continue
				}
				if err == io.EOF {
					break
				}
				w.logger.Fatal(err)
			}
			nChunks++
			nBytes += int64(len(buf))
			if w.debug {
				w.logger.Println("Writing data to channel")
			}
			c <- buf
			if w.debug {
				w.logger.Printf("data: %s", buf)
			}
			if err != nil && err != io.EOF {
				w.logger.Fatal(err)
			}
		}
	}

	go stdinHandler(stdin, stdinChannel)

	if w.debug {
		w.logger.Printf("Starting wrapped program [%s %s]\n", cmd.Path, cmd.Args[:])
	}

	err = cmd.Run()
	if err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				// preserve stderr
				if stderr := stderrbuf.String(); len(stderr) > 0 {
					w.logger.Printf("Compiler %s failed: %v", actualProg, err)
					fmt.Fprintf(os.Stderr, "%s", stderr)
				}
				// preserve non-zero return code
				os.Exit(status.ExitStatus())
			}
		} else {
			log.Fatalf("Calling compiler %v failed: %v", actualProg, err)
		}
	}

	// preserve stdout
	if stdout := stdoutbuf.String(); len(stdout) > 0 {
		fmt.Fprintf(os.Stdout, "%s", stdout)
	}

	if w.debug {
		w.logger.Print("Actual compiler finished successfully")
	}
}

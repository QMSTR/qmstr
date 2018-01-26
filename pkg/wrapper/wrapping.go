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
)

// Wrapper represents a wrapper to call a program
type Wrapper struct {
	logger          *log.Logger
	Program         string
	commandlineArgs []string
	debug           bool
}

// NewWrapper returns an instance of a Wrapper for the given command line
func NewWrapper(commandline []string, logger *log.Logger, debug bool) *Wrapper {
	// extract the compiler that was supposed to run
	w := Wrapper{}
	w.logger = logger
	w.Program = filepath.Base(commandline[0])
	//extract the arguments
	w.commandlineArgs = commandline[1:]
	return &w
}

// Wrap calls the actual program to be wrapped and preserves output and return value
func (w *Wrapper) Wrap() {

	if w.Program == "qmstr-wrapper" {
		log.Fatal("This is not how you should invoke the qmstr-wrapper.\n\tSee https://github.com/QMSTR/qmstr for more information on how to use the QMSTR.")
	}

	// find and run actual compiler
	actualProg, err := FindActualProgram(w.Program)
	if err != nil {
		log.Fatalf("Actual compiler was not found. %v", err)
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
		w.logger.Println("Starting wrapped program")
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

// CheckExecutable checks the given file to be no directory and executable flagged
func CheckExecutable(file string) error {
	d, err := os.Stat(file)
	if err != nil {
		return err
	}
	if m := d.Mode(); !m.IsDir() && m&0111 != 0 {
		return nil
	}
	return os.ErrPermission
}

// FindActualProgram discovers the actual program that is wrapper on the PATH
func FindActualProgram(prog string) (string, error) {
	path := os.Getenv("PATH")
	foundWrapper := false
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := filepath.Join(dir, prog)
		if err := CheckExecutable(path); err == nil {
			if foundWrapper {
				return path, nil
			}
			// First hit is the wrapper
			foundWrapper = true
		}
	}
	return "", fmt.Errorf("executable file %s not found in [%s]", prog, path)
}

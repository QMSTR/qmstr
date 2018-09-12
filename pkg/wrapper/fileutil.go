package wrapper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/pkg/common"
)

// FindActualProgram discovers the actual program that is wrapped on the PATH
func FindActualProgram(prog string) (string, error) {
	me, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("could not find initial executable that started this process: %v", err)
	}
	meAbs, err := filepath.Abs(me)
	if err != nil {
		return "", fmt.Errorf("failed finding absolute path: %v", err)
	}

	paths := FindExecutablesOnPath(prog)
	countProgs := len(paths)
	for i, path := range paths {
		pathAbs, err := filepath.Abs(path)
		if err != nil {
			continue
		}

		pathtocheck, err := filepath.EvalSymlinks(pathAbs)
		if err != nil {
			return "", fmt.Errorf("failed evaluation: %v", err)
		}

		if filepath.Clean(meAbs) == pathtocheck {
			// +1 next program after me is the target
			if countProgs <= i+1 {
				break
			}
			return paths[i+1], nil
		}
	}
	return "", fmt.Errorf("executable file %s not found in %v", prog, paths)
}

// FindExecutablesOnPath finds and returns all reachable executables for the given progname
func FindExecutablesOnPath(progname string) []string {
	var paths []string
	path := os.Getenv("PATH")
	for _, dir := range filepath.SplitList(path) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		path := filepath.Join(dir, progname)
		if err := common.CheckExecutable(path); err == nil {
			paths = append(paths, path)
		}
	}
	return paths
}

package wrapper

import (
	"fmt"
	"os"
	"path/filepath"
)

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

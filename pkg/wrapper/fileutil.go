package wrapper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/pkg/common"
)

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
		if err := common.CheckExecutable(path); err == nil {
			if foundWrapper {
				return path, nil
			}
			// First hit is the wrapper
			foundWrapper = true
		}
	}
	return "", fmt.Errorf("executable file %s not found in [%s]", prog, path)
}

package wrapper

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func BuildCleanPath(base string, subpath string) string {
	if filepath.IsAbs(subpath) {
		return filepath.Clean(subpath)
	}

	if !filepath.IsAbs(base) {
		// ignore error and use non absolute path
		base, _ = filepath.Abs(base)
	}
	tmpPath := filepath.Join(base, subpath)
	return filepath.Clean(tmpPath)
}

func Hash(fileName string) (string, error) {
	h := sha256.New()
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	nBytes, nChunks := int64(0), int64(0)
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := f.Read(buf[:cap(buf)])
		buf = buf[:n]
		if n == 0 {
			if err == nil {
				continue
			}
			if err == io.EOF {
				break
			}
			return "", err
		}
		nChunks++
		nBytes += int64(len(buf))
		h.Write(buf)
		if err != nil && err != io.EOF {
			return "", err
		}
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
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

// exists checks if file exists and is not a directory
func exists(file string) bool {
	if d, err := os.Stat(file); err == nil {
		if d.IsDir() {
			return false
		}
		return true
	}
	return false
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

// FindActualProgram discovers the actual program that is wrapper on the PATH
func FindActualLibraries(libs []string, libpath []string) ([]string, error) {
	actualLibPaths := []string{}
	syslibpath := []string{"/lib", "/usr/lib"}
	for _, lib := range libs {
		for _, dir := range append(libpath, syslibpath...) {
			if dir == "" {
				// Unix shell semantics: path element "" means "."
				dir = "."
			}
			err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
				if f.Name() == fmt.Sprintf("lib%s.so", lib) {
					actualLibPaths = append(actualLibPaths, path)
					return fmt.Errorf("Found %s", path)
				}
				return nil

			})
			if err != nil {
				break
			}
		}
	}

	if len(actualLibPaths) == len(libs) {
		return actualLibPaths, nil
	}

	return actualLibPaths, fmt.Errorf("Missing libraries from %v in %v", libs, actualLibPaths)
}

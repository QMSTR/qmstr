package wrapper

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

const (
	libPathVar = "LIBRARY_PATH"
)

func BuildCleanPath(base string, subpath string, abs bool) string {
	if filepath.IsAbs(subpath) {
		return filepath.Clean(subpath)
	}

	if abs && !filepath.IsAbs(base) {
		// ignore error and use non absolute path
		base, _ = filepath.Abs(base)
	}

	tmpPath := filepath.Join(base, subpath)
	return filepath.Clean(tmpPath)
}

func Hash(fileName string) (string, error) {
	h := sha1.New()
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
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

// FindActualLibraries discovers the actual libraries on the path
func FindActualLibraries(libs []string, libpath []string) ([]string, error) {
	actualLibPaths := []string{}
	libpathvar, present := os.LookupEnv(libPathVar)
	if present && libpathvar != "" {
		libpath = append([]string{libpathvar}, libpath...)
	}
	var libprefix string
	var libsuffix []string
	var syslibpath []string
	switch runtime.GOOS {
	case "linux":
		libprefix = "lib"
		libsuffix = []string{".so"}
		syslibpath = []string{"/lib", "/usr/lib", "/usr/local/lib", "/lib64"}
	case "darwin":
		libprefix = "lib"
		libsuffix = []string{".dylib", ".so"}
		syslibpath = []string{"/usr/lib", "/usr/local/lib"}
	case "windows":
		libprefix = ""
		libsuffix = []string{".dll"}
		syslibpath = []string{""}
	}
	for _, lib := range libs {
		for _, dir := range append(libpath, syslibpath...) {
			if dir == "" {
				// Unix shell semantics: path element "" means "."
				dir = "."
			}
			err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
				for _, suffix := range libsuffix {
					if f.Name() == fmt.Sprintf("%s%s%s", libprefix, lib, suffix) {
						actualLibPaths = append(actualLibPaths, path)
						return fmt.Errorf("Found %s", path)
					}
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

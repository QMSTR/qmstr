package builder

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/QMSTR/qmstr/pkg/qmstr/service"
)

const (
	libPathVar = "LIBRARY_PATH"
)

type Builder interface {
	Analyze(commandline []string) (*service.BuildMessage, error)
	GetName() string
	GetPrefix() (string, error)
}

type GeneralBuilder struct {
	Logger *log.Logger
	Debug  bool
}

func NewFileNode(path string, fileType string) *service.FileNode {
	filename := filepath.Base(path)
	hash, err := hash(path)
	broken := false
	if err != nil {
		hash = "nohash" + path
		broken = true
	}
	return &service.FileNode{Name: filename, Type: fileType, Path: path, Hash: hash, Broken: broken}
}

func hash(fileName string) (string, error) {
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

// FindActualLibraries discovers the actual libraries on the path
func FindActualLibraries(actualLibs map[string]string, linkLibs []string, libPath []string, staticLink bool, staticLibs map[string]struct{}) error {
	libpathvar, present := os.LookupEnv(libPathVar)
	if present && libpathvar != "" {
		libPath = append([]string{libpathvar}, libPath...)
	}
	var libprefix string
	var libsuffix []string
	var syslibpath []string
	switch runtime.GOOS {
	case "linux":
		libprefix = "lib"
		libsuffix = []string{".so", ".a"}
		syslibpath = []string{"/lib", "/usr/lib", "/usr/local/lib", "/lib64"}
	case "darwin":
		libprefix = "lib"
		libsuffix = []string{".dylib", ".so", ".a"}
		syslibpath = []string{"/usr/lib", "/usr/local/lib"}
	case "windows":
		libprefix = ""
		libsuffix = []string{".dll"}
		syslibpath = []string{""}
	}

	for _, dir := range append(libPath, syslibpath...) {
		if dir == "" {
			// Unix shell semantics: path element "" means "."
			dir = "."
		}
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			for _, lib := range linkLibs {
				// is lib located
				if _, ok := actualLibs[lib]; ok {
					continue
				}
				var suffixes []string
				// forced static lib or linking statically
				if _, ok := staticLibs[lib]; ok || staticLink {
					suffixes = []string{".a"}
				} else {
					suffixes = libsuffix
				}

				for _, suffix := range suffixes {
					if f.Name() == fmt.Sprintf("%s%s%s", libprefix, lib, suffix) {
						actualLibs[lib] = path
					}
				}
			}
			return nil
		})
	}

	if len(actualLibs) == len(linkLibs) {
		return nil
	}

	return fmt.Errorf("Missing libraries from %v in %v", linkLibs, actualLibs)
}

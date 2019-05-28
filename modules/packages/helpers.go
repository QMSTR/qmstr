package packages

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/lib/go-qmstr/validation"
)

func PackageFromFile(path string) (validation.Package, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".deb":
		return NewDebPackage(f)
	default:
		return nil, fmt.Errorf("unknown package %s", ext)
	}
}

func validateFileInfos(a, b validation.FileInfo) error {
	if a.SHA1 != b.SHA1 {
		return validation.HashMissmatchError{
			Algo: "SHA1",
			A:    a.SHA1,
			B:    b.SHA1,
		}
	}
	return nil
}

func indexByName(fis []validation.FileInfo) map[string]validation.FileInfo {
	m := map[string]validation.FileInfo{}
	for _, fi := range fis {
		m[fi.Name] = fi
	}
	return m
}

func indexByHash(fis []validation.FileInfo) map[string]validation.FileInfo {
	m := map[string]validation.FileInfo{}
	for _, fi := range fis {
		m[fi.SHA1] = fi
	}
	return m
}

package common

import (
	"crypto/sha1"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/QMSTR/qmstr/pkg/database"
	"github.com/QMSTR/qmstr/pkg/service"
)

var nonPosixChars = regexp.MustCompile(`[^A-Za-z0-9\._-]`)

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

// IsFileExist checks if file IsFileExist and is not a directory
func IsFileExist(file string) bool {
	if d, err := os.Stat(file); err == nil {
		if d.IsDir() {
			return false
		}
		return true
	}
	return false
}

func SetRelativePath(node *service.FileNode, buildPath string, pathSub []*service.PathSubstitution) error {
	for _, substitution := range pathSub {
		node.Path = strings.Replace(node.Path, substitution.Old, substitution.New, 1)
	}
	if !filepath.IsAbs(node.Path) {
		return nil
	}
	relPath, err := filepath.Rel(buildPath, node.Path)
	if err != nil {
		return err
	}
	node.Path = relPath
	return nil
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
		if err := CheckExecutable(path); err == nil {
			paths = append(paths, path)
		}
	}
	return paths
}

func HashFile(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	return Hash(f)
}

func Hash(r io.Reader) (string, error) {
	h := sha1.New()
	buf := make([]byte, 0, 4*1024)
	for {
		n, err := r.Read(buf[:cap(buf)])
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

func SanitizeFileNode(f *service.FileNode, base string, pathSub []*service.PathSubstitution, db *database.DataBase, parentPath string) error {
	if err := SetRelativePath(f, base, pathSub); err != nil {
		return err
	}
	if f.Hash == "" {
		log.Printf("No hash for file %s", f.Path)
		var hash string
		var err error
		if f.Path == parentPath {
			log.Println("Override detected")
			hash, err = db.GetFileNodeHashByPath(f.Path)
			if err != nil {
				return fmt.Errorf("Corrupted data provided. File does not exist: %v", err)
			}
			log.Printf("Found original hash %s in database\n", hash)
		} else {
			hash, err = HashFile(filepath.Join(base, f.Path))
			if err != nil {
				return err
			}
			log.Printf("Calculated hash %s\n", hash)
		}
		f.Hash = hash
	}
	fileParts := strings.Split(f.Path, "/")
	// catch tmp files
	if fileParts[0] == ".." && f.FileType == service.FileNode_SOURCE {
		f.FileType = service.FileNode_INTERMEDIATE
	}

	for _, d := range f.DerivedFrom {
		if err := SanitizeFileNode(d, base, pathSub, db, f.Path); err != nil {
			return err
		}
	}
	for _, dep := range f.Dependencies {
		if err := SanitizeFileNode(dep, base, pathSub, db, f.Path); err != nil {
			return err
		}
	}
	return nil
}

func GetPosixFullyPortableFilename(filename string) string {
	posixFilename := nonPosixChars.ReplaceAllString(filename, "_")
	return posixFilename
}

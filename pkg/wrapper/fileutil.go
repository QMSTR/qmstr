package wrapper

import "path/filepath"

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

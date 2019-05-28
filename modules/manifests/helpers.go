package manifests

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/QMSTR/qmstr/lib/go-qmstr/validation"
)

func ManifestFromFile(path string) (validation.Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	ext := filepath.Ext(path)
	switch ext {
	case ".spdx":
		return NewSPDXManifest(f)
	default:
		return nil, fmt.Errorf("unknown manifest %s", ext)
	}
}

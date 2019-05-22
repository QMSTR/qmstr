package pkg

import "fmt"

// PackageInfo holds generic package metadata for manifest and package
type PackageInfo struct {
	Name string
}

// FileInfo holds generic file metadata for manifest and package
type FileInfo struct {
	Name string
	SHA1 string
}

// Manifest the interface to be implemented by manifests
type Manifest interface {
	PackageInfo() PackageInfo
	FileInfo() []FileInfo
}

// Package the interface to be implemented by packages
type Package interface {
	Validate(manifest Manifest) error
}

// HashMissmatchError is returned by Validate methods when hashes missmatch
type HashMissmatchError struct {
	Algo string
	A, B string
}

func (e HashMissmatchError) Error() string {
	return fmt.Sprintf("hash %s missmatch: %s <> %s", e.Algo, e.A, e.B)
}

// FileNotInManifestError is returned by Validate methods when a file is
// documented in the manifest but is not in the package
type FileNotInManifestError struct {
	Name string
}

func (e FileNotInManifestError) Error() string {
	return fmt.Sprintf("file %s not documented in manifest", e.Name)
}

// FileNotInPackageError is returned by Validate methods when a file is
// in the package in the manifest but is not documented in the manifest
type FileNotInPackageError struct {
	Name string
}

func (e FileNotInPackageError) Error() string {
	return fmt.Sprintf("file %s not present in package", e.Name)
}

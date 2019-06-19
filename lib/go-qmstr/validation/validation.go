package validation

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

// Error extents the error interface with ExitCode method
type Error interface {
	error
	ExitCode() int
}

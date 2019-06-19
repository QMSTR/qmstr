package validation

import "fmt"

// UnknownError is for errors unknown to the validator
type UnknownError struct {
}

func (e UnknownError) Error() string {
	return "unknown error"
}

func (e UnknownError) ExitCode() int {
	return 1
}

// HashMissmatchError is returned by Validate methods when hashes missmatch
type HashMissmatchError struct {
	Algo string
	A, B string
}

func (e HashMissmatchError) Error() string {
	return fmt.Sprintf("hash %s missmatch: %s <> %s", e.Algo, e.A, e.B)
}

func (e HashMissmatchError) ExitCode() int {
	return 2
}

// FileNotInManifestError is returned by Validate methods when a file is
// documented in the manifest but is not in the package
type FileNotInManifestError struct {
	Name string
}

func (e FileNotInManifestError) Error() string {
	return fmt.Sprintf("file %s not documented in manifest", e.Name)
}

func (e FileNotInManifestError) ExitCode() int {
	return 3
}

// FileNotInPackageError is returned by Validate methods when a file is
// in the package in the manifest but is not documented in the manifest
type FileNotInPackageError struct {
	Name string
}

func (e FileNotInPackageError) Error() string {
	return fmt.Sprintf("file %s not present in package", e.Name)
}

func (e FileNotInPackageError) ExitCode() int {
	return 4
}

// MissingLicenseInfoError is returned when a manifest has no license info for
// a file described
type MissingLicenseInfoError struct {
	Name string
}

func (e MissingLicenseInfoError) Error() string {
	return fmt.Sprintf("file %s has no licese info", e.Name)
}

func (e MissingLicenseInfoError) ExitCode() int {
	return 5
}

// MissingCopyrightInfoError is returned when a manifest has no copyright info
// for a file described
type MissingCopyrightInfoError struct {
	Name string
}

func (e MissingCopyrightInfoError) Error() string {
	return fmt.Sprintf("file %s has no copyright info", e.Name)
}

func (e MissingCopyrightInfoError) ExitCode() int {
	return 6
}

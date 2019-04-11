---
title: "The QMSTR Validator"
date: 2019-04-01T15:26:15Z
draft: false
weight: 10
---

The QMSTR Validator is a command line program that validates the
documentation for a package or a distribution. It verifies the content
of packages against the compliance documentation and answers three
questions:

* Does the documentation match the package? The documentation matches
  the package if the checksums for the files contained in the package
  match the information in the documentation. This check will fail,
  for example, if a binary has been changed or rebuilt after the
  documentation was created so that the checksums do not match.
* Does the package only contain files that are documented? This test
  will fail if there are files in the package that are not described
  in the documentation.
* Is the documentation complete? The documentation is complete if it
  describes all files in the package and contains all required
  information about these files. This check will fail, for example,
  if license or authorship information is missing even though the file
  is listed in the documentation.

To validate a package, the QMSTR Validator must have been implemented
to support that package format. For each format, a standard or a best
practice needs to exist that describes how the package compliance
documentation and the package payload content are shipped
together. Initially, the validator supports the Debian file format as
a starting point.

## Invoking the QMSTR Validator

The QMSTR Validator is invoked on a package and the corresponding
manifest file:

	> qmstr validate curl-a.b.c.deb curl-a.b.c.deb.spdx
	Verifying package content against manifest curl-a.b.c.deb.spdx...
	* documentation matches the package content
	* all files in the package are documented
	* documentation is complete
	Package validation passed.

In case of a validation error, the return code of the validator
depends on which of the checks failed (1, 2 or 3). If the validation
succeeds the return code is zero.

Depending on the file format, the manifest file may be included in the
package. Support for that may be added to the validator.

## Verifying distributions

In the QMSTR context, a distribution is a set of packages that is
shipped together as a unit. The validator can be used to evaluate
distributions. The validation succeeds if the validation for each
package contained in the distribution succeeds.

TODO How this can be done has to be developed and depends on how
distributions are represented.

## Development plan

The first working version of the QMSTR Validator will be implemented
against the Debian package file format and plain tar files.

## Dependencies

The QMSTR Validator runs stand-alone and does not require a QMSTR
master to be available.

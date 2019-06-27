---
title: "The QMSTR Validator"
date: 2019-04-01T15:26:15Z
draft: false
weight: 10
---

The QMSTR Validator is a command line program that validates the
manifests for a package or a distribution. It verifies the content
of packages against the compliance documentation and answers three
questions:

* Does the manifest match the package? The documentation matches
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
  is listed in the manifest.

To validate a package, the QMSTR Validator must have been implemented
to support that package format. For each format, a standard or a best
practice needs to exist that describes how the package compliance
documentation and the package payload content are shipped
together. Initially, the validator supports the Debian file format as
a starting point.

The validator validates one package, however "traversing into
elements". That means it first validates that the distribution is
"correct", and then validates that the same goes for all element
packages. There should be no separate concepts for packages and
distributions. Distributions are simply packages that contain other
packages instead of files to be installed. Or - if you "install" a
distribution file into a directory, you get the element packages.

## Terminology

* *Package*: A package is a concrete "form of archive" intended for
  distribution that is created during a build. It could be a tarball,
  a Debian or RPM package, or even a file system overlay that contains
  only what was packaged. For the first implementation, we want to
  focus specifically on Debian packages and then abstract the concept.
* *Distribution*: A distribution is a collection of packages. The
  concept is not related to "Linux distribution", but to a software
  delivery. It is a file-like entity, not a process (as in
  handover). If *distribution* is imagined as a file, it is also a
  *package*. The concepts are purposefully recursive, so that
  validation can be applied to both exactly the same way.
* *Software delivery*: A software delivery is the handover and
  acceptance of a package (or set of packages, a distribution) from a
  supplier to a customer. Software delivery always involves two
  parties, one that provides a package with a maniest and one that
  receives those and validates them against each other.
* *Manifest*: A manifest in the context of validation is a SPDX file
  that that describes exactly one package, including its sub-packages
  (which are files), but not the content of the sub-packages. The
  manifest for a distribution describes the content of the
  distribution (the packages in the collection). The elements of the
  distribution contain or bring their own manifests.

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

The basic algorithm is

* validate the package itself
* for each element package
  * validate the element package
* succeed if both steps validate, fail if not

## Development plan

### QMSTR 0.6: TRL 6

For QMSTR 0.6, validation is suppossed to
reach
[H2020 TRL 6](https://ec.europa.eu/info/funding-tenders/opportunities/portal/screen/support/faq;keywords=/2890) ("technology
demonstrated in relevant environment "). In detail, this means:

* A first working version of the QMSTR Validator will be implemented
  against the Debian package file format.
* The package manifest reporter will be implemented to produce
  manifests that match the format expectations of the validator
  (one-to-one match between package amd manifest, mandatory manifest
  content).
* The validator will be implemented to validate the package against
  the manifests generated by the package manifest reporter.
* The QMSTR CURL demo will be extended to showcase the validation
  feature. The demonstration will consist of two separate steps, build
  and validation (the steps are separate since they could be performed
  by different entities in real life, simulating a software delivery):
  * On the build side, the demo should build curl and create the .deb
    packages for it, as well as the manifests for the packages.
  * On the validation side, the demo should validate the packages
    against the generated manifests.

The validation feature will be merged into the QMSTR master branch
after this demonstration.

### QMSTR 0.7: TRL 7

The validation feature will be demonstrated at conferences in
Q4/2019. Based on the feedback, the feature will be completed at TRL 7
("system prototype demonstration in operational environment"). Other
packages formats may be added, in particular plain tar files and
Android .apk packages.

## Dependencies

The QMSTR Validator runs stand-alone and does not require a QMSTR
master to be available. The package manifest reporter runs as part of
the QMSTR reporting phase.

## Potential additional features (iceboxed)

### Validation Levels

Once the validation basics are complete, it is conceivable that the
results may correspond to different validation levels:

* level 1: manifest has at least checksums
* level 2: manifest has rights holder and license information for all artifacts
* level 3: ???

### Delivery receipts

On delivery, the validator could generate and cryptographically sign a
delivery receipt. Delivery receipts could be collected in a central
place, like an online repository or a blockchain.

### Self-contained packages

Package formats could be extended so that they contain the
manifests. This means that the software delivery would only contain
package files that can be validated against themselfes.

### Additional package formats

Feature completeness for the implemented formats currently has higher
priority than adding additional formats. However, if there is enough
interest, support for additional package formats should be added.

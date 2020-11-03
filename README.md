# Quartermaster - the FOSS Compliance Toolchain that is itself FOSS

[Quartermaster](http://qmstr.org) is a suite of command line tools and build system extensions that instruments software builds to create
FOSS compliance documentation and support compliance decisions. It executes as part of a software build process to generate reports about the analyzed product.

## Usage

See the [deployment instructions](deploy/README.md#deployment-instructions).

## How it works

Quartermaster runs adjacent to a software build process. A master
process collects information about the software that is built. Once
the build is complete, the master executes a number of analysis tools,
and finally a number of reporters. Which exactly is configured in  a
configuration file called `qmstr.yaml` that usually resides in the
root directory of the repository.

## Current status

As of now, Quartermaster only analyzes Maven projects and [stores license and compliance information in its database](deploy/README.md#results).

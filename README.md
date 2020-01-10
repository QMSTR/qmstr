# Quartermaster - the FOSS Compliance Toolchain that is itself FOSS

[Quartermaster](http://qmstr.org) is a suite of command line tools and build system extensions that instruments software builds to create
FOSS compliance documentation and support compliance decisions. It executes as part of a software build process to generate reports about the analysed product.

## Continuous Integration Status

| master (build and unit tests) | documentation ([shelldoc](https://github.com/endocode/shelldoc)) |
| ------------- | ------------- |
| [![Build Status](https://ci.endocode.com/buildStatus/icon?job=QMSTR%2Fqmstr%2Fdevelopment)](https://ci.endocode.com/job/QMSTR/job/qmstr/job/development/) | [![Build Status](https://ci.endocode.com/buildStatus/icon?job=QMSTR/DocTest)](https://ci.endocode.com/view/QMSTR/job/QMSTR/job/DocTest/) |

## Compilation and Installation

The client and master side installation of Quartermaster is Makefile
based. See the [installation instructions](https://qmstr.org/documentation/introduction/installation/) for
details. The client side tools need to be installed on all build
machines. The master is only needed once per build process.

To develop on Quartermaster, a [few extra steps](https://qmstr.org/documentation/developer)
may be required.

## Basics

Quartermaster runs adjacent to a software build process. A master
process collects information about the software that is build. Once
the build is complete, the master executes a number of analysis tools,
and finally a number of reporters. Which exactly is configured in  a
configuration file called `qmstr.yaml` that usually resides in the
root directory of the repository.

All modules are executed in the context of the master, not the build
machine. The master ships all dependencies of the modules without
affecting the build clients file system (it runs in a container).

Read more in the [Getting started](https://qmstr.org/documentation/introduction/getting-started/) page.

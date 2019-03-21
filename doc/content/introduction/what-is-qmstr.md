---
title: "What is QMSTR?"
date: 2019-01-17T15:26:15Z
lastmod: 2019-01-17T15:26:15Z
draft: false
weight: 1
---

Quartermaster is an integrated free and open source software (FOSS)
toolchain that implements industry best practises of license
compliance management. It is FOSS because the QMSTR community believes
that software that manages FOSS compliance becomes critical
infrastructure to the wider open source community and should itself be
FOSS. QMSTR is developed in collaboration with users, legal experts
and all other interested stakeholders.

## QMSTR: A command-line workflow toolchain for FOSS compliance

QMSTR is implemented as a command line toolchain. It assumes that
almost everywhere, software is built and packaged using command line
based build systems that automate the build steps. By integrating into
the common software development tooling, QMSTR programs may be used in
terminals, scripts, Makefiles, Jenkinsfiles and other places. QMSTR
aggregates the knowledge it acquires about the software being built in
a graph database held in a central master process.

The workflow QMSTR implements proceeds through three separate phases,
a build phase, an analysis phase and a reporting phase. The purpose of
the build phase is to produce a build graph. For that, QMSTR
integrates natively into build systems like Maven, or uses a fallback
of wrapping the compiler and linker if the build system cannot
otherwise produce a build graph (as in traditional Makefile based
builds). This means QMSTR can support any build system that produces
or can be shoehorned into revealing a build graph. The build graph
describes which sources and dependencies are used to create which
targets and in what way, and how these artifacts are being assembled
into a package for distribution.

The analysis phase utilizes a range of tools to augment the knowledge
graph with metadata about the software for example by executing
static analysis tools or retrieving information from other data
sources. By combining the build graph with the results of the analysis
phase, QMSTR is able to determine accurately what license and
authorship metadata is relevant for the package that is finally
distributed. It is also able to identify which elements of the source
code are tests, examples or other fragments that do not get
distributed and do not affect license and authorship of the
distributed package.

The reporting phase formats the results in the knowledge graph for a
specific need. Creating a package SPDX manifest is done by using a
reporter, as is triggering notification or producing CI test
results. During the reporting phase, no additional data may be added
to the knowledge graph anymore. It is considered frozen.

## Modular analysis and reporting architecture

The analysis and reporting functionality in QMSTR is implemented in
modules. Modules are stand-alone programs that are executed by the
master and have access to the knowledge graph and the source and build
directories. Modules communicate with the master using the gRPC based
master API. This means that modules may be written in any language
that comes with support for gRPC communication. The modules shipped
with QMSTR are commonly implemented in Go or Python. QMSTR for example
ships modules to perform license or authorship analysis using Git or
Scancode or to read and write SPDX manifests.

Implementing custom modules and adding them to the master installation
is relatively straightforward. Next to the modules shipped with QMSTR,
users may add custom analysis or reporting modules that have the same
access to the knowledge graph and the source and build artifacts as
the built-in ones. This facilitates the implementation for example of
checks on compliance with specific business policies or the integrate
reporting with internal software delivery platforms. Since modules run
as stand-alone processes and communicate with the master using gRPC,
the implementer is free in the choice of license and may even create
custom proprietary modules. All modules shipped with QMSTR are however
FOSS.

## APIs instead of file formats

The QMSTR master is started right before the build process begins and
terminates after the reporting phase ends. There are no long-running
services as part of a QMSTR installation. To persist any results,
reports need to be created and stored. All data collected during the
build and analysis phase is held in a graph database managed by the
master and made available to the modules by the master API. Without
reporting (or caching), the data is ephemeral. QMSTR specifically
avoids creating file formats for data exchange. SPDX exists for that
purporse and is used in different places in QMSTR. Other formats can
be consumed by analysis or created by reporting modules.

## Integration into CI/CD pipelines

Software integration today should be automated and performed in
continuous integration (CI) pipelines. Since the steps performed by CI
pipelines are commonly regular command line instructions, QMSTR can be
integrated into a CI pipeline in the same style as the other
instructions for the individual build steps. Endocode runs
a
[public cURL with QMSTR demo](https://ci.endocode.com/view/QMSTR/job/QMSTR/job/qmstr-cURL-demo/) that
showcases how all the steps necessary to perform a build with QMSTR
can be integrated into CI without the need for specific integration
(open the "Blue Ocean" view for more details).

## Summary: Building blocks of an open compliance program

QMSTR creates the license compliance information that comes with a
distributed software package. If everything is set up properly and the
metadata in the project source code is well-maintained, the resulting
reports reduce legal and contractual uncertainty and streamline
delivery along supply chains.

QMSTR is however only one building block of an open compliance
program. It complements (and is partly developed in cooperation with
the communities) [SPDX](https://spdx.org/) as the main license data
exchange format and [OpenChain](https://www.openchainproject.org/) as
the specification of supply chain requirement. The combination of
formats, tooling and processes -- all available under FOSS licenses --
provides all the necessary tools to maintain FOSS license compliance.

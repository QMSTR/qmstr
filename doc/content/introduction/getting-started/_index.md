---
title: "Getting started"
date: 2019-01-17T15:26:15Z
draft: false
weight: 4
---

There are usually three groups of users interested in working with
Quartermaster: open source officers, continuous integration operators,
and software developers. Open source officers use documentation
produced by Quartermaster to manage compliance of software
products. Continuous integration (CI) operators set up software builds
that are instrumented with Quartermaster to produce compliance
documentation when software packages are created. Software developers
integrate Quartermaster into their development setup to monitor
compliance and reproduce test results. This tutorial will begin to
explain the basics and building blocks of Quartermaster from the
perspective of a software developer. The CI operator and open source
office perspective will be covered in later articles.

## Prerequisites

To execute Quartermaster as part of a software build, the client side
tools and the master container image need to be available on the
system. Follow the [installation instructions](Installation.md) to get
those installed.

Verify that the client side tools are available by querying the
version of the `qmstrctl` tool:

    > qmstrctl version
    This is qmstrctl version 0.5

To check that the master image exists, run

    > docker images | grep qmstr/master
    ...

For the remainder of this tutorial, Git, internet access and the basic
development tools to build a C library on Linux will be required.

## Next

In the next pages you will be guided to build different projects with Quartermaster.

Before even trying to get started with Quartermaster, make sure that
the project that is going to be analyzed (we call it the _project
under analysis_) builds properly in your
environment. Otherwise it may be difficult later to separate issues
with running Quartermaster from regular build errors. We generally
recommend to first make sure the project under analysis builds and
then add Quartermaster instrumentation in a second step.

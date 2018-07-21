# Getting started with Quartermaster

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
office perspetive will be covered in later articles.

## Prerequisites

To execute Quartermaster as part of a software build, the client side
tools and the master container image need to be available on the
system. Follow the [installation instructions](Installation.md) to get
those installed.

Verify that the client side tools are available by querying the version of the `qmstrctl` tool:

    > qmstrctl version
    This is qmstrctl version 0.1.

To check that the master image exists, run

    >docker images | grep qmstr/master
    ...

For the remainder of this tutorial, Git, internet access and the basic
development tools to build a C library on Linux will be required.

## Step 0: Making sure your software builds

Before even trying to get started with Quartermaster, make sure that
the project that is going to be analyzed builds properly in your
environment. Otherwise it may be difficult later to separate issues
with running Quartermaster from regular build errors.

This tutorial will uses the JSON-C library as the
library-under-analysis.

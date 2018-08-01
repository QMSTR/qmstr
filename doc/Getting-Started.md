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

Verify that the client side tools are available by querying the
version of the `qmstrctl` tool:

    > qmstrctl version
    This is qmstrctl version 0.1.

To check that the master image exists, run

    > docker images | grep qmstr/master
    ...

For the remainder of this tutorial, Git, internet access and the basic
development tools to build a C library on Linux will be required.

## Step 0: Making sure your software builds

Before even trying to get started with Quartermaster, make sure that
the project that is going to be analyzed (we call it the _project
under analysis_) builds properly in your
environment. Otherwise it may be difficult later to separate issues
with running Quartermaster from regular build errors. We generally
recommend to first make sure the project under analysis builds and
then add Quartermaster instrumentation in a second step.

This tutorial will use the JSON-C library as the
project under analysis. The files referenced in the following steps
are locatedd in the [tutorial](tutorial/) subdirectory. Let's retrieve
the JSON-C source code first, in a specific revision that we know
works with this tutorial:

	> cd doc/tutorial
	> git clone git@github.com:json-c/json-c.git
	...
	> cd json-c
	> git reset --hard bf29aa0f
	...

All JSON-C specific parts of the build have been automated in a [build
script](tutorial/build-json-c.sh). Let's run it to make sure our
environment is configured to build JSON-C from scratch:

	> ../build-json-c.sh
	...

If this works, wonderful. If not, please dig into the output and check
for errors. Quartermaster tracks what sources are compiled and what
targets are linked in your project. To be sure that everything gets
compiled again after instrumentation, let's clean the repository:

	> git clean -fxd
	> cd ..

## Step 1: Start a Quartermaster master process

Every build that is instrumented with Quartermaster needs a master
process. Every Quartermaster master needs a configuration file,
usually called [qmstr.yaml](tutorial/qmstr.yaml). The configuration
file is located in the tutorial/ directory because in our case, we are
avoiding to make changes to the project under analysis. Let's start
the master:

	> cd json-c
	> . `qmstrctl start --wait --config ../qmstr.yaml`

The `wait` flag makes sure that the command returns only after the
master has finished starting up and is fully operational. The `config`
flag points to the configuration file. If it is not specified,
`qmstrctl` looks for a file called `qmstr.yaml` in  the current
directory.






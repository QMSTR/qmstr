# Installing Quartermaster

A Quartermaster installation consists of a client and a master side
part. The client side builds the software under inspection, and
collects information about the objects that are being built. It
transmits that information to the master which performs the
analysis and creates the output artifacts (reports). The client side
uses programs native to the build environment. The master side always
runs in a containerized Linux system.

## Installing the clients

The main entry point into the installation tasks for Quartermaster is
the Makefile in the main repository. The default installation installs
the client programs into `/usr/local/bin`:

	> make install_qmstr_client
	...

Depending on the specifics of the local setup, a developer may want to
install the binaries into the GOPATH:

	> make install_qmstr_client_gopath
	...

If the installation completes successfully, the qmstrctl command is
now available:

	> qmstrctl version
	This is qmstrctl version 0.1.

Only the client side installation is required on a system that builds
software with Quartermaster instrumentation. All tools and programs
required to perform analysis and create reports are included in the
master container.

## Master installation

The master process runs as a container. The Dockerfile to build the
container resides in the main Quartermaster repository. The
Quartermaster project does not provide ready-made images in a
container repository, the images have to be built on the systems that
are supposed to run the master. This may change at a later time,
however at the moment it is the only mechanism to prepare to run the
master. To crate the master image, run make again:

	> make container
	...

This may take a while. It will build the master container, including
all Quartermaster modules and the tools that are used for
analysis. Once the images have been created, `qmstrctl` can be used to
start and manage the master process.

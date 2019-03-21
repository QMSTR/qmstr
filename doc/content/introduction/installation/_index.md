---
title: "Installing Quartermaster"
date: 2019-01-17T15:26:15Z
lastmod: 2019-01-17T15:26:15Z
draft: false
weight: 3
---

A Quartermaster installation consists of a client and a master side
part. The client side builds the software under inspection, and
collects information about the objects that are being built. It
transmits that information to the master which performs the analysis
and creates the output artifacts (reports). The client side uses
programs native to the build environment. The master side always runs
in a containerized Linux system. To run the master, the master
container image needs to be available on the system where `qmstrctl
start` is called. Other build slaves or distributed compile nodes only
need the client side programs, since the master only runs once per
build.

The core modules and tools of Quartermaster are written in the Go
programming language. However it is not primarily a Go
project. Modules may be written in any programming language that
provides GRPC bindings. They are usually written in the language of
the tools they integrate with. The SPDX parser module, for example, is
written in Python. The Gradle integration is written in Java. The
Quartermaster build system uses a Makefile to implement the different
steps required to build these modules in their specific ways.

## Prerequisites

Quartermaster leverages Protobuf & GRPC for the master - client
communication. The `protoc` Protobuf  compiler needs to be installed
and also the Go protobuf library and generator. The later can be
installed with

	> go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

### Host machine preparation:
Install The depencencies:

  - If you are running an Ubuntu machine:

    > sudo apt update

	> sudo apt install golang protobuf-compiler

    - Install Docker: https://docs.docker.com/install/linux/docker-ce/ubuntu/

  - If you are running a Fedora machine:

    > sudo dnf install golang protobuf-compiler

  - Install Docker: https://docs.docker.com/install/linux/docker-ce/fedora/


Add user to the docker group:
  - Create new group if it does not exist.

  > sudo groupadd docker

  -  Add current user to the group

  > sudo gpasswd -a $USER docker

  - Reload shell. For that log out and log back or execute the next command:

  > newgrp docker

  More information in:
  https://linoxide.com/linux-how-to/use-docker-without-sudo-ubuntu/

Install The `protoc` Protobuf  compiler:

	> go get -u github.com/golang/protobuf/{proto,protoc-gen-go}

Make $GOPATH/bin part of you $PATH, for that edit the ~/.bashrc file and
add the next at the end of it:

        # Make $GOPATH/bin part of you $PATH
        export GOPATH="${HOME}/go"
        export PATH="${PATH}:${GOPATH}/bin"

If that doesn't work try the next:

        # Make $GOPATH/bin part of you $PATH
        export PATH="${PATH}:${HOME}/go/bin"

## Installing the clients

Get the sources:

> go get -u github.com/QMSTR/qmstr

> cd ~/go/src/github.com/QMSTR/qmstr

* Note: If you are running `go get -u github.com/QMSTR/qmstr` on the Ci it
might crash and obtain the next error message:
`Package github.com/QMSTR/qmstr: no Go files in /home/go/src/github.com/QMSTR/qmstr`
If that's the case try the next:

> go get -u github.com/QMSTR/qmstr || true

The main entry point into the installation tasks for Quartermaster is
the Makefile in the main repository. The default installation installs
the client programs into `/usr/local/bin`:

	> make install_qmstr_client
	...

Depending on the specifics of the local setup, a developer may want to
install the binaries into the GOPATH:

	> make install_qmstr_client_gopath
	...

If the installation completes successfully, the `qmstrctl` command is
now available:

	> qmstrctl version
	This is qmstrctl version 0.2.

Only the client side installation and the information about how to
access the master are required on a system that builds software with
Quartermaster instrumentation. All tools and programs required to
perform analysis and create reports are included in the master
container.

For other make targets, inspect the Makefile.

### Further notes:

If the steps above didn't work try the next:

  > make

  > make install_qmstr_client_gopath


## Master installation

The master process runs as a container. The Dockerfile to build the
container resides in the main Quartermaster repository. The
Quartermaster project does not provide ready-made images in a
container repository, the images have to be built on or made available
to the systems that are supposed to run the master. This may change at
a later time, however at the moment it is the only mechanism to
prepare to run the master.

To crate the master image, run make again:

	> make container
	...

This may take a while. It will build the master container, including
all Quartermaster modules and install the dependencies and tools that
are used for analysis. Once the images have been created, `qmstrctl`
can be used to start and manage the master process.

## Summary

Every system that builds parts of your software (build slave,
executor, your CI may call the build clients in a different way) needs
to have the Quartermaster client tools installed. The easiest way to
ensure that is usually to include the instructions to install the
client programs in the automated setup of the build slaves. The master
only needs to run once per build process, even if parallel build
processes are used that deploy build jobs to a number of build
clients. In most scenarios, the master images need to be built on the
machine where the software build is started.

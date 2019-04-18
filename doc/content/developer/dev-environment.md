---
title: "Development environment"
date: 2019-04-01T15:26:15Z
weight: 20
---

To create a local development environment for Quartermaster, follow
these instructions. The master by default runs in a container. See
below for more details.

## Prerequisites

Quartermaster leverages Protobuf & GRPC for the master - client
communication. The `protoc` Protobuf compiler needs to be installed
together with the Go protobuf library and generator. The later will be installed via the `Makefile`.

### Host machine preparation

Install the depencencies:

  - If you are running an Ubuntu machine:

    > sudo apt update

	> sudo apt install golang protobuf-compiler

    - Install Docker: https://docs.docker.com/install/linux/docker-ce/ubuntu/

  - If you are running a Fedora machine:

    > sudo dnf install golang protobuf-compiler

  - Install Docker: https://docs.docker.com/install/linux/docker-ce/fedora/

## Dependencies

Quartermaster uses protobuf for the communication between clients and the master. Install `protobuf` from your package manager.

## Checkout sources

  > git clone https://github.com/QMSTR/qmstr.git

## Build and run the Quartermaster master process

Quartermaster uses a multi-stage [Dockerfile](masterserver/Dockerfile) to create various setups based on a common configuration. The DGraph database process and the Quartermaster master are executed in the container.

In order to build the different parts of the qmstr system `make` is used.

### Building qmstr

There are several targets defined in the `Makefile` to build the different parts of qmstr:
The client tools:
	- qmstrctl - client tool to communicate with the master server

	  > make qmstrctl

	- qmstr - wrapper tool to set up and run programs in a qmstrized environment

	  > make qmstr

Automatic targets are created to build every module inside the `modules/{builders,analyzers,reporters}` directories e.g to build the `spdx-analyzer` that can be found in `modules/analyzers/spdx-analyzer` run:

	> make spdx-analyzer

In order to assemble a full master image use the `master` target:

	> make master

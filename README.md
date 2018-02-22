# Quartermaster - the FOSS Compliance Toolchain that is itself FOSS

[Quartermaster](http://qmstr.org) is a suite of command line tools and build system extensions that instruments software builds to create
FOSS compliance documentation and support compliance decisions. It executes as part of a software build process to generate reports about the analysed product.

[![Build Status](https://ci.endocode.com/buildStatus/icon?job=QMSTR/build_and_unit_test_master)](https://ci.endocode.com/job/QMSTR/job/build_and_unit_test_master/)

## Compilation and Installation

To create a local development environment for Quartermaster, follow these instructions. The master by default runs in a container. See below for more details.

### Dependencies

Quartermaster uses protobuf for the communication between clients and the master. Install protobuf from your package manager.

Then install the protobuf go compiler plugin:

`go get -u github.com/golang/protobuf/protoc-gen-go`

### Install Quartermaster

Make sure $GOBIN is part of your $PATH.

Install the master server
`(go get github.com/QMSTR/qmstr/cmd/qmstr-master || go generate github.com/QMSTR/qmstr/cmd/qmstr-master; go get github.com/QMSTR/qmstr/cmd/qmstr-master)`

Install the wrapper
`go get github.com/QMSTR/qmstr/cmd/qmstr-wrapper`

Optional: install the cli
`go get github.com/QMSTR/qmstr/cmd/qmstr-cli`

Or if you dare `wget -O - http://github.com/QMSTR/qmstr/raw/master/install.sh | bash`

## Build and run the Quartermaster master process

Quartermaster uses a multi-stage [Dockerfile](ci/Dockerfile) to create various setups based on a common configuration. The DGraph database process and the Quartermaster master are executed in the container.

### Performing a clean build: builder

The `builder` stage compiles the Quartermaster toolchain in a clean Go environment. It does not contain any runtime dependencies.

`docker build -f ci/Dockerfile -t qmstr/build --target builder .`

### Running the unit tests

The `master_unit_tests` stage is a minimal extension of `builder` that adds an entry point to execute the master unit tests.

`docker build -f ci/Dockerfile -t qmstr/master_unit_tests --target master_unit_tests . && \
    docker run --rm -it qmstr/tests`

### Dependencies for analysis and reporting: `runtime`

The `runtime` stage contains a full operating system environment and the default tools that Quartermaster uses to perform analysis of the build graph and to create reports.

`docker build -f ci/Dockerfile -t qmstr/runtime --target runtime .`

### Build the "official" Quartermaster master container

The master container contains all analysis and reporting dependencies, and the compiled Quartermaster toolchain, but no development environment. It is the default way to run a Quartermaster master process.

`docker build -f ci/Dockerfile -t qmstr/master --target master .`

### Create a container for development

The development container is a combination of the runtime environment with a source volume to test local changes. It uses a volume to pass in source code under development, and builds the source code in it's entrypoint script.

`docker build -f ci/Dockerfile -t qmstr/dev --target dev .
docker run -it -p 50051:50051 -v $HOME/Go/src:/go/src <build_path>:/buildroot qmstr/dev`

...where `build_path` is the path to the source files you are about to compile.

# Setting up an environment to hack on Quartermaster

To create a local development environment for Quartermaster, follow
these instructions. The master by default runs in a container. See
below for more details.

## Dependencies

Quartermaster uses protobuf for the communication between clients and the master. Install `protobuf` from your package manager.

Then install the protobuf go compiler plugin:

	% go get -u github.com/golang/protobuf/protoc-gen-go

## Install Quartermaster

Make sure $GOBIN is part of your $PATH.

Install the master server

	% (go get github.com/QMSTR/qmstr/cmd/qmstr-master || go generate github.com/QMSTR/qmstr/cmd/qmstr-master; go get github.com/QMSTR/qmstr/cmd/qmstr-master)

Install the wrapper

	% go get github.com/QMSTR/qmstr/cmd/qmstr-wrapper

Optional: install the cli

	% go get github.com/QMSTR/qmstr/cmd/qmstrctl

Or if you dare:

	% wget -O - http://github.com/QMSTR/qmstr/raw/master/install.sh | bash

## Build and run the Quartermaster master process

Quartermaster uses a multi-stage [Dockerfile](docker/Dockerfile) to create various setups based on a common configuration. The DGraph database process and the Quartermaster master are executed in the container.

### Performing a clean build: builder

The `builder` stage compiles the Quartermaster toolchain in a clean Go environment. It does not contain any runtime dependencies.

	% docker build -f docker/Dockerfile -t qmstr/build --target builder .

### Running the unit tests

The `master_unit_tests` stage is a minimal extension of `builder` that adds an entry point to execute the master unit tests.

	% docker build -f docker/Dockerfile -t qmstr/master_unit_tests --target master_unit_tests .
	% docker run --rm -it qmstr/master_unit_tests

### Dependencies for analysis and reporting: `runtime`

The `runtime` stage contains a full operating system environment and the default tools that Quartermaster uses to perform analysis of the build graph and to create reports.

	% docker build -f docker/Dockerfile -t qmstr/runtime --target runtime .

### Build and run the "official" Quartermaster master container

The master container contains all analysis and reporting dependencies, and the compiled Quartermaster toolchain, but no development environment. It is the default way to run a Quartermaster master process.

	% docker build -f docker/Dockerfile -t qmstr/master --target master .
	% docker run -it -p 50051:50051 -v <build_path>:/buildroot qmstr/master

...where `build_path` is the path to the source files you are about to compile.

### Create and run a container for development

The development container is a combination of the runtime environment with a source volume to test local changes. It uses a volume to pass in source code under development, and builds the source code in it's entrypoint script.

	% docker build -f docker/Dockerfile -t qmstr/dev --target dev .
	% docker run -it -p 50051:50051 -v $HOME/Go/src:/go/src <build_path>:/buildroot qmstr/dev

...where `build_path` is the path to the source files you are about to compile.

### Debug in development container

You can use the development container to debug qmstr-master.

    % export QMSTR_DEV="<debugger_port>"
    % docker run -p 50051:50051 -p $QMSTR_DEV:2345 -eQMSTR_DEV=$QMSTR_DEV --security-opt seccomp=unconfined -v $HOME/Go/src:/go/src <build_path>:/buildroot qmstr/dev

Now you can connect your debugger.

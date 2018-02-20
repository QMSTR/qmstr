#!/bin/sh
#
# Temporary script to try and document the different container builds:

# Build the toolchain from the code in the master branches on Github:
docker rmi qmstr/build || true
docker build -f ci/Dockerfile -t qmstr/build --target builder .

# Run tests:
docker rmi qmstr/master_unit_tests || true
docker build -f ci/Dockerfile -t qmstr/master_unit_tests --target master_unit_tests . && \
    docker run --rm -it qmstr/tests

# Create a runtime environment that contains the built code and all
# the tools integrated in the master container:
docker rmi qmstr/runtime || true
docker build -f ci/Dockerfile -t qmstr/runtime --target runtime .

# Create the master container to be published to Docker hub: For this
# one, everything should be set up, so there should be no entrypoint
# script:
docker rmi qmstr/master || true
docker build -f ci/Dockerfile -t qmstr/master --target master .
# docker run -it qmstr/master

# Create the development container: A combination of the runtime
# environment with a source volume to test local changes:
docker rmi qmstr/dev || true
docker build -f ci/Dockerfile -t qmstr/dev --target dev .
docker run -it -v $HOME/Go/src:/go/src qmstr/dev

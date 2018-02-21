# qmstr

Quartermaster is a suite of command line tools and build system plugins that instruments software builds to create
FLOSS compliance documentation and support compliance decision making. It executes as part of a software build process
to generate reports about the analysed product.

[![Build Status](https://ci.endocode.com/buildStatus/icon?job=QMSTR/build_and_unit_test_master)](https://ci.endocode.com/job/QMSTR/job/build_and_unit_test_master/)

# Compile and Install

Install protobuf from your package manager

Install protobuf go compiler plugin
`go get -u github.com/golang/protobuf/protoc-gen-go`

Make sure $GOBIN is part of your $PATH.

Install the master server
`(go get github.com/QMSTR/qmstr/cmd/qmstr-master || go generate github.com/QMSTR/qmstr/cmd/qmstr-master; go get github.com/QMSTR/qmstr/cmd/qmstr-master)`

Install the wrapper
`go get github.com/QMSTR/qmstr/cmd/qmstr-wrapper`

Optional: install the cli
`go get github.com/QMSTR/qmstr/cmd/qmstr-cli`

Or if you dare `wget -O - http://github.com/QMSTR/qmstr/raw/master/install.sh | bash`

# Building and running qmstr in docker container

To build qmstr in docker container from repo root run
`docker build -f ci/Dockerfile -t qmstr --target builder .`

To build and run qmstr in docker container from repo root run
`docker build -f ci/Dockerfile -t qmstr/master . &&
docker run -p 50051:50051 -v <build_path>:/buildroot qmstr/master`

Where `build_path` is the path to the source files you are about to compile.

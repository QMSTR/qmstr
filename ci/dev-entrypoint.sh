#!/bin/bash
set -e

function start_dgraph() {
    dgraph zero &
    dgraph server --memory_mb=2048 --zero=localhost:5080 &
}

function start_dgraph_web {
    dgraph-ratel &
}

# Generate and build
go generate github.com/QMSTR/qmstr/cmd/qmstr-master
go install github.com/QMSTR/qmstr/cmd/qmstr-master

start_dgraph
start_dgraph_web

exec /go/bin/qmstr-master

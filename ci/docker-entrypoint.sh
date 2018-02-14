#!/bin/bash
set -e

function start_dgraph() {
    dgraph zero &
    dgraph server --memory_mb=2048 --zero=localhost:5080 &
}

start_dgraph

# Give dgraph some time to come up
sleep 15

exec /go/bin/qmstr-master

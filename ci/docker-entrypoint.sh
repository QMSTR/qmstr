#!/bin/bash
set -e

function start_dgraph() {
    dgraph zero &
    dgraph server --lru_mb=2048 --zero=localhost:5080 &
}

start_dgraph

exec /go/bin/qmstr-master --config /buildroot/qmstr.yaml

#!/bin/bash
set -e

function start_dgraph() {
    dgraph version
    dgraph zero &
    dgraph server --lru_mb=2048 --zero=localhost:5080 &
}


start_dgraph

exec /usr/local/bin/qmstr-master --config /buildroot/qmstr.yaml

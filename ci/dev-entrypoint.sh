#!/bin/bash
set -e

function start_dgraph() {
    dgraph version
    dgraph zero &
    dgraph server --lru_mb=2048 --zero=localhost:5080 &
}

function start_dgraph_web {
    dgraph-ratel &
}

# Generate and build
/usr/local/bin/install-qmstr.sh -d
/usr/local/bin/install-qmstr.sh -i
(cd $GOPATH/src/github.com/QMSTR/qmstr/cmd/reporters/qmstr-reporter-html && ./setup.sh /usr/share/qmstr $GOPATH/src/github.com/QMSTR/qmstr)

start_dgraph
start_dgraph_web

if [ -z "$QMSTR_DEV" ]; then
    exec /go/bin/qmstr-master --config /buildroot/qmstr.yaml
else
    echo "Running debug session"
    exec dlv debug github.com/QMSTR/qmstr/cmd/qmstr-master -l 0.0.0.0:2345 --headless=true --log=true -- --config /buildroot/qmstr.yaml
fi

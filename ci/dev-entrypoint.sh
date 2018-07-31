#!/bin/bash
set -e

source /common.inc

# Generate and build
pushd $GOPATH/src/github.com/QMSTR/qmstr/
make install_qmstr_server
pushd cmd/modules/reporters/qmstr-reporter-html
./setup.sh /usr/share/qmstr $GOPATH/src/github.com/QMSTR/qmstr
popd
popd

start_dgraph
start_dgraph_web

create_qmstr_user

if [ -z "$QMSTR_DEBUG" ]; then
    start_qmstr
else
    echo "Running debug session"
    exec dlv debug github.com/QMSTR/qmstr/cmd/qmstr-master -l 0.0.0.0:2345 --headless=true --log=true -- --config /qmstr/qmstr.yaml ${PATH_SUB:+--pathsub="$PATH_SUB"}
fi

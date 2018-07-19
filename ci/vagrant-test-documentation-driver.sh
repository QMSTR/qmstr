#!/bin/sh
#
# WARNING: Only call this to be executed in an automated CI
# setup. This script will modify the home directory of the current
# user!
#
# This script is to be called from the repository root.
set -e
export GOPATH=$HOME/test_go
export PATH=$GOPATH/bin:$PATH
BASEDIR=$GOPATH/github/QMSTR
mkdir -p $BASEDIR
cd $BASEDIR
ln -sf /vagrant qmstr
cd qmstr/
./ci/test-documentation.sh


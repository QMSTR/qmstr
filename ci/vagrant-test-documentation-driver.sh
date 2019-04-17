#!/bin/sh
#
# WARNING: Only call this to be executed in an automated CI
# setup. This script will modify the home directory of the current
# user and the host file system!
#
# This script is to be called from the repository root. It also
# assumes that /vagrant points to a working copy of the QMSTR/qmstr
# repository.
set -e
sudo chown -R $USER:$USER /usr/local/
export GOPATH=$HOME/go
export PATH=$GOPATH/bin:/opt/go/bin:$PATH
cd /vagrant
./ci/test-documentation.sh


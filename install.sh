#!/bin/bash
set -e

function downloadqmstr() {
    go get -u -v github.com/QMSTR/qmstr || true
}

function installdeps() {
    go get -u -v github.com/golang/dep/cmd/dep
    pushd $GOPATH/src/github.com/QMSTR/qmstr
    dep ensure
    echo "install the go protobuf comiler plugin"
    pushd vendor/github.com/golang/protobuf/protoc-gen-go
    go install
    popd
    popd
}

function installqmstr() {
    go generate github.com/QMSTR/qmstr/cmd/qmstr-master
    go install -v github.com/QMSTR/qmstr/cmd/qmstr-master
    go install -v github.com/QMSTR/qmstr/cmd/qmstr-wrapper
    go install -v github.com/QMSTR/qmstr/cmd/qmstrctl
    go install -v github.com/QMSTR/qmstr/cmd/qmstr
    go install -v github.com/QMSTR/qmstr/cmd/analyzers/pkg-analyzer
    go install -v github.com/QMSTR/qmstr/cmd/analyzers/spdx-analyzer
    go install -v github.com/QMSTR/qmstr/cmd/analyzers/scancode-analyzer
    go install -v github.com/QMSTR/qmstr/cmd/analyzers/test-analyzer
    go install -v github.com/QMSTR/qmstr/cmd/analyzers/git-analyzer
    go install -v github.com/QMSTR/qmstr/cmd/reporters/test-reporter
    go install -v github.com/QMSTR/qmstr/cmd/reporters/qmstr-reporter-html
    pushd $GOPATH/src/github.com/QMSTR/qmstr
    make clean
    make install_python_modules
    popd
}

function usage() {
    echo "install.sh -q     # download qmstr sources"
    echo "install.sh -d     # download and install dependencies"
    echo "install.sh -i     # install qmstr only"
    echo "install.sh        # download and install"
    exit 1
}

TASK=''

while getopts ":diq" opt; do
  case $opt in
    d)
      echo "installing deps" >&2
      [[ -n "$TASK" ]] && usage || TASK='deps'
      ;;
    q)
      echo "download qmstr sources" >&2
      [[ -n "$TASK" ]] && usage || TASK='download'
      ;;
    i)
      echo "installing qmstr only" >&2
      [[ -n "$TASK" ]] && usage || TASK='install'
      ;;
    \?)
      echo "Invalid option: -$OPTARG" >&2
      usage
      ;;
  esac
done

case $TASK in
    download)
         downloadqmstr
         ;;
    deps)
         installdeps
         ;;
    install)
         installdeps
         installqmstr
         ;;
    '')
         echo "full clean install"
         downloadqmstr
         installdeps
         installqmstr
         ;;
esac

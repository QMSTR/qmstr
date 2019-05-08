#!/bin/sh
set -e

cd json-c
echo "Building and documenting JSON-C..."
sh autogen.sh
./configure
make
echo "Building JSON-C completed..."

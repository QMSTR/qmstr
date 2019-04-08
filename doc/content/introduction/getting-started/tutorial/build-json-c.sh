#!/bin/sh
set -e

echo "Building and documenting JSON-C..."
sh autogen.sh
./configure
qmstrctl create package:json-c
make
echo "Building JSON-C completed..."

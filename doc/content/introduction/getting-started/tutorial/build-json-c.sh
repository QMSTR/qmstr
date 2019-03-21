#!/bin/sh
set -e

echo "Building and documenting JSON-C..."
sh autogen.sh
./configure
make
echo "Building JSON-C completed..."

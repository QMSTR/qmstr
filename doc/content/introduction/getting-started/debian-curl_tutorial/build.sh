#!/bin/sh
set -e

cd curl
echo "Building Debian cURL..."
# Make sure to run as root
sudo dpkg-buildpackage -B -us -uc

echo "Building Debian cURL completed..."
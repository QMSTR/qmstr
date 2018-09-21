#!/bin/sh
#
# Test the Quartermaster documentation (written in Markdown) using the
# shelldoc tool (https://github.com/endocode/shelldoc) in the current
# environment.
set -e

# test that we are in the correct location
if [ ! -f docker/test-documentation.sh ]; then
    echo "Start this script from the repository root!"
    exit 2
fi

# test the Markdown files under doc/
FILES="doc/Installation.md doc/Getting-Started.md"
# test the README:
FILES="$FILES README.md"
# Aaaand go...:
shelldoc $FILES

#!/bin/sh
#
# Test the Quartermaster documentation (written in Markdown) using the
# shelldoc tool (https://github.com/endocode/shelldoc) in the current
# environment.
set -e

# test that we are in the correct location
if [ ! -f ci/test-documentation.sh ]; then
    echo "Start this script from the repository root!"
    exit 2
fi

# test the README:
FILES=README.md
# test the Markdown files under doc/
FILES="$FILES doc/*md"
# Aaaand go...:
shelldoc $FILES

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

# test the Markdown files under doc/
FILES="doc/content/introduction/installation/_index.md doc/content/introduction/getting-started/_index.md doc/content/introduction/getting-started/json-c_tutorial/_index.md doc/content/introduction/getting-started/debian-curl_tutorial/_index.md"
# test the README:
FILES="$FILES README.md"
# Aaaand go...:
shelldoc run $FILES

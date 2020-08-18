#!/bin/bash

set -e
for binary in "qmstr-master test-analyzer test-reporter"; do
    which ${binary} >/dev/null || (echo "Binary ${binary} not found. Exiting." && exit 1)
done

qmstr-master $@


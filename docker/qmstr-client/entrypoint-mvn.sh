#/bin/bash

set -e
cd /home/qmstrclient/buildroot/guava

mvn -pl .,guava clean package


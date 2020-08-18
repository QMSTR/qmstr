#!/bin/bash

set -e
cd /var/qmstr/buildroot/guava

mvn -pl .,guava clean package


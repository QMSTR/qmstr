#/bin/bash

set -e
POM=$(find . -name "pom.xml")
cd $(dirname $POM)
mvn compile


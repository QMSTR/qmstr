#/bin/bash

set -e
declare -rg PROJECT_DIR=${BUILDROOT}/${REPO_NAME}

cd ${PROJECT_DIR}

mvn -pl .,guava clean package


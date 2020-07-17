#/bin/bash

set -e
PROJECT_DIR=${BUILDROOT}/${REPO_NAME}

cd ${PROJECT_DIR}

mvn -pl .,guava clean package


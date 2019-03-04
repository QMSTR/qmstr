#!/bin/bash

qmstr_master=$(docker ps |grep -m1 qmstr/master| awk '{ print $1 }')
if [ -z ${qmstr_master} ]; then
    echo "master container not found" 1>&2
    exit 1
fi

QMSTR_MASTER=$(docker inspect --format '{{range .NetworkSettings.Ports}}{{range .}}{{.HostIp}}:{{.HostPort}}{{end}}{{end}}' ${qmstr_master})
if [ -z ${QMSTR_MASTER} ]; then
    echo "master container port settings not found" 1>&2
    exit 2
fi

export QMSTR_MASTER=${QMSTR_MASTER}
echo "export QMSTR_MASTER=${QMSTR_MASTER}"


#!/bin/bash

if [ -z "$QMSTR_MASTER" ]
then
    echo "QMSTR_MASTER not set" 1>&2
    exit 1
fi

fwd_port=$(echo $QMSTR_MASTER | cut -d ":" -f 2)
if [ -z "${fwd_port##*[!0-9]*}" ]
then
    echo "no port found" 1>&2
fi

mastercontainer=$(docker ps | grep $QMSTR_MASTER | awk '{ print $1 }')
if [ -z "${mastercontainer}" ]
then
    echo "master container not found" 1>&2
fi

network=$(docker inspect --format "{{.HostConfig.NetworkMode}}" ${mastercontainer})
if [ -z "${network}" ]
then
    echo "network not found" 1>&2
fi


docker run --rm -d -e MASTERCONTAINER="${mastercontainer}" -p 8080:8080 -p 8000:8000 --network ${network} qmstr-ratel
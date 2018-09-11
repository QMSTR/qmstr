#!/bin/bash

set -e

# set up forwarding the dgraph port to access dgraph from ratel via localhost 
socat tcp-listen:8080,fork,reuseaddr tcp-connect:${MASTERCONTAINER}:8080 &

exec /usr/local/bin/dgraph-ratel
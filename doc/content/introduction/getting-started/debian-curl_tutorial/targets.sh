#!/bin/bash
set -e

function hashthis() {
  sha1sum $1 | awk '{ print $1 }'
}

# getOldestHash() uses the path and the timestamp of a node
# and returns the hash of the oldest node
function getOldestHash() {
	TIMESTAMP=`qmstrctl describe file:path:$1 | grep $1 | awk '/Timestamp/{print $NF}' | sort -n | head -1`
    TIMESTAMP="/${TIMESTAMP}/"
	qmstrctl describe file:path:$1 | grep $1 | awk -F'[ ,]' '$TIMESTAMP{print $9}' | head -1
}

# create curl targets
CURL_BDIR=curl/debian/curl/usr/
qmstrctl create file:${CURL_BDIR}share/doc/curl/NEWS.Debian.gz --name NEWS.Debian.gz
qmstrctl create file:${CURL_BDIR}share/doc/curl/changelog.gz --name changelog.gz
qmstrctl create file:${CURL_BDIR}share/doc/curl/copyright --name copyright
qmstrctl create file:${CURL_BDIR}share/doc/curl/changelog.Debian.gz --name changelog.Debian.gz
qmstrctl create file:${CURL_BDIR}share/man/man1/curl.1.gz --name curl.1.gz
qmstrctl create file:${CURL_BDIR}share/zsh/vendor-completions/_curl --name _curl

# connect targets to curl package
qmstrctl connect package:curl_7.64.0-3_amd64.deb \
	file:hash:$(hashthis ${CURL_BDIR}bin/curl) \
	file:${CURL_BDIR}share/doc/curl/NEWS.Debian.gz \
	file:${CURL_BDIR}share/doc/curl/changelog.gz \
	file:${CURL_BDIR}share/doc/curl/copyright \
	file:${CURL_BDIR}share/doc/curl/changelog.Debian.gz \
	file:${CURL_BDIR}share/man/man1/curl.1.gz \
	file:${CURL_BDIR}share/zsh/vendor-completions/_curl

# connect missing dependencies
qmstrctl connect file:hash:$(getOldestHash ${CURL_BDIR}bin/curl) \
	file:path:curl/debian/build/src/.libs/curl

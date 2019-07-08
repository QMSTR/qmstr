#!/bin/bash
set -e

function hashthis() {
  sha1sum $1 | awk '{ print $1 }'
}

# create curl targets
CURL_BDIR=curl/debian/curl/usr/
qmstrctl create file:${CURL_BDIR}share/doc/curl/NEWS.Debian.gz
qmstrctl create file:${CURL_BDIR}share/doc/curl/changelog.gz
qmstrctl create file:${CURL_BDIR}share/doc/curl/copyright
qmstrctl create file:${CURL_BDIR}share/doc/curl/changelog.Debian.gz
qmstrctl create file:${CURL_BDIR}share/man/man1/curl.1.gz
qmstrctl create file:${CURL_BDIR}share/zsh/vendor-completions/_curl

# connect targets to curl package
qmstrctl connect package:curl_7.64.0-3_amd64.deb \
	file:hash:$(hashthis ${CURL_BDIR}bin/curl) \
	file:${CURL_BDIR}share/doc/curl/NEWS.Debian.gz \
	file:${CURL_BDIR}share/doc/curl/changelog.gz \
	file:${CURL_BDIR}share/doc/curl/copyright \
	file:${CURL_BDIR}share/doc/curl/changelog.Debian.gz \
	file:${CURL_BDIR}share/man/man1/curl.1.gz \
	file:${CURL_BDIR}share/zsh/vendor-completions/_curl



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

# create libcurl4 targets
LIBCURL_BDIR=curl/debian/libcurl4/usr/
qmstrctl create file:${LIBCURL_BDIR}share/doc/libcurl4/NEWS.Debian.gz
qmstrctl create file:${LIBCURL_BDIR}share/doc/libcurl4/changelog.Debian.gz
qmstrctl create file:${LIBCURL_BDIR}share/doc/libcurl4/changelog.gz
qmstrctl create file:${LIBCURL_BDIR}share/doc/libcurl4/copyright

# connect targets to libcurl4 package
qmstrctl connect package:libcurl4_7.64.0-3_amd64.deb \
	file:hash:$(hashthis ${LIBCURL_BDIR}lib/x86_64-linux-gnu/libcurl.so.4.5.0) \
	file:${LIBCURL_BDIR}share/doc/libcurl4/NEWS.Debian.gz \
	file:${LIBCURL_BDIR}share/doc/libcurl4/changelog.Debian.gz \
	file:${LIBCURL_BDIR}share/doc/libcurl4/changelog.gz \
	file:${LIBCURL_BDIR}share/doc/libcurl4/copyright

# create libcurl3-gnutls targets
GNU_BDIR=curl/debian/libcurl3-gnutls/usr/
qmstrctl create file:${GNU_BDIR}share/lintian/overrides/libcurl3-gnutls

GNU_DOC=${GNU_BDIR}share/doc/libcurl3-gnutls/
qmstrctl create file:${GNU_DOC}NEWS.Debian.gz
qmstrctl create file:${GNU_DOC}changelog.Debian.gz
qmstrctl create file:${GNU_DOC}changelog.gz
qmstrctl create file:${GNU_DOC}copyright

# connect targets to libcurl3-gnutls package
qmstrctl connect package:libcurl3-gnutls_7.64.0-3_amd64.deb \
	file:hash:$(hashthis ${GNU_BDIR}lib/x86_64-linux-gnu/libcurl-gnutls.so.4.5.0) \
	file:${GNU_DOC}NEWS.Debian.gz \
	file:${GNU_DOC}changelog.Debian.gz \
	file:${GNU_DOC}changelog.gz \
	file:${GNU_DOC}copyright \
	file:${GNU_BDIR}share/lintian/overrides/libcurl3-gnutls

# create libcurl3-nss targets
NSS_DIR=curl/debian/libcurl3-nss/usr/
qmstrctl create file:${NSS_DIR}share/lintian/overrides/libcurl3-nss

NSS_DOC=${NSS_DIR}share/doc/libcurl3-nss/
qmstrctl create file:${NSS_DOC}NEWS.Debian.gz
qmstrctl create file:${NSS_DOC}changelog.Debian.gz
qmstrctl create file:${NSS_DOC}changelog.gz
qmstrctl create file:${NSS_DOC}copyright

# connect targets to libcurl3-nss package
qmstrctl connect package:libcurl3-nss_7.64.0-3_amd64.deb \
	file:hash:$(hashthis ${NSS_DIR}lib/x86_64-linux-gnu/libcurl-nss.so.4.5.0) \
	file:${NSS_DOC}NEWS.Debian.gz \
	file:${NSS_DOC}changelog.Debian.gz \
	file:${NSS_DOC}changelog.gz \
	file:${NSS_DOC}copyright \
	file:${NSS_DIR}share/lintian/overrides/libcurl3-nss

# create libcurl4-openssl-dev targets
OPENSSL_BDIR=curl/debian/libcurl4-openssl-dev/usr/
qmstrctl create file:${OPENSSL_BDIR}bin/curl-config

H_FILES=${OPENSSL_BDIR}include/x86_64-linux-gnu/curl/
qmstrctl create file:${H_FILES}curl.h
qmstrctl create file:${H_FILES}curlver.h
qmstrctl create file:${H_FILES}easy.h
qmstrctl create file:${H_FILES}mprintf.h
qmstrctl create file:${H_FILES}multi.h
qmstrctl create file:${H_FILES}stdcheaders.h
qmstrctl create file:${H_FILES}system.h
qmstrctl create file:${H_FILES}typecheck-gcc.h
qmstrctl create file:${H_FILES}urlapi.h

qmstrctl create file:${OPENSSL_BDIR}lib/x86_64-linux-gnu/libcurl.la
qmstrctl create file:${OPENSSL_BDIR}lib/x86_64-linux-gnu/pkgconfig/libcurl.pc

qmstrctl create file:${OPENSSL_BDIR}share/aclocal/libcurl.m4
qmstrctl create file:${OPENSSL_BDIR}share/man/man1/curl-config.1.gz

OPENSSL_DOC=${OPENSSL_BDIR}share/doc/libcurl4-openssl-dev/
qmstrctl create file:${OPENSSL_DOC}NEWS.Debian.gz
qmstrctl create file:${OPENSSL_DOC}changelog.Debian.gz
qmstrctl create file:${OPENSSL_DOC}changelog.gz
qmstrctl create file:${OPENSSL_DOC}copyright


# connect targets to libcurl4-openssl-dev package
qmstrctl connect package:libcurl4-openssl-dev_7.64.0-3_amd64.deb \
	file:${OPENSSL_DOC}NEWS.Debian.gz \
	file:${OPENSSL_DOC}changelog.Debian.gz \
	file:${OPENSSL_DOC}changelog.gz \
	file:${OPENSSL_DOC}copyright \
	file:${OPENSSL_BDIR}share/aclocal/libcurl.m4

qmstrctl connect package:libcurl4-openssl-dev_7.64.0-3_amd64.deb \
	file:hash:$(hashthis ${OPENSSL_BDIR}lib/x86_64-linux-gnu/libcurl.a) \
	file:${OPENSSL_BDIR}lib/x86_64-linux-gnu/libcurl.la \
	file:${OPENSSL_BDIR}lib/x86_64-linux-gnu/pkgconfig/libcurl.pc \
	file:${OPENSSL_BDIR}bin/curl-config \
	file:${OPENSSL_BDIR}share/man/man1/curl-config.1.gz

qmstrctl connect package:libcurl4-openssl-dev_7.64.0-3_amd64.deb \
	file:${H_FILES}curl.h \
	file:${H_FILES}curlver.h \
	file:${H_FILES}easy.h \
	file:${H_FILES}mprintf.h \
	file:${H_FILES}multi.h \
	file:${H_FILES}stdcheaders.h \
	file:${H_FILES}system.h \
	file:${H_FILES}typecheck-gcc.h \
	file:${H_FILES}urlapi.h

# create libcurl4-gnutls-dev targets
GNUTLS_DEV_BDIR=curl/debian/libcurl4-gnutls-dev/usr/
qmstrctl create file:${GNUTLS_DEV_BDIR}bin/curl-config

H_FILES=${GNUTLS_DEV_BDIR}include/x86_64-linux-gnu/curl/
qmstrctl create file:${H_FILES}curl.h
qmstrctl create file:${H_FILES}curlver.h
qmstrctl create file:${H_FILES}easy.h
qmstrctl create file:${H_FILES}mprintf.h
qmstrctl create file:${H_FILES}multi.h
qmstrctl create file:${H_FILES}stdcheaders.h
qmstrctl create file:${H_FILES}system.h
qmstrctl create file:${H_FILES}typecheck-gcc.h
qmstrctl create file:${H_FILES}urlapi.h

qmstrctl create file:${GNUTLS_DEV_BDIR}lib/x86_64-linux-gnu/libcurl-gnutls.la
qmstrctl create file:${GNUTLS_DEV_BDIR}lib/x86_64-linux-gnu/pkgconfig/libcurl.pc

qmstrctl create file:${GNUTLS_DEV_BDIR}share/aclocal/libcurl.m4
qmstrctl create file:${GNUTLS_DEV_BDIR}share/man/man1/curl-config.1.gz

GNUTLS_DOC=${GNUTLS_DEV_BDIR}share/doc/libcurl4-gnutls-dev/
qmstrctl create file:${GNUTLS_DOC}NEWS.Debian.gz
qmstrctl create file:${GNUTLS_DOC}changelog.Debian.gz
qmstrctl create file:${GNUTLS_DOC}changelog.gz
qmstrctl create file:${GNUTLS_DOC}copyright

# connect targets to libcurl3-nss package
qmstrctl connect package:libcurl4-gnutls-dev_7.64.0-3_amd64.deb \
	file:${GNUTLS_DEV_BDIR}bin/curl-config \
	file:${H_FILES}curl.h \
	file:${H_FILES}curlver.h \
	file:${H_FILES}easy.h \
	file:${H_FILES}mprintf.h \
	file:${H_FILES}multi.h \
	file:${H_FILES}stdcheaders.h \
	file:${H_FILES}system.h \
	file:${H_FILES}typecheck-gcc.h \
	file:${H_FILES}urlapi.h

qmstrctl connect package:libcurl4-gnutls-dev_7.64.0-3_amd64.deb \
	file:hash:$(hashthis ${GNUTLS_DEV_BDIR}lib/x86_64-linux-gnu/libcurl-gnutls.a) \
	file:${GNUTLS_DEV_BDIR}lib/x86_64-linux-gnu/libcurl-gnutls.la \
	file:${GNUTLS_DEV_BDIR}lib/x86_64-linux-gnu/pkgconfig/libcurl.pc

qmstrctl connect package:libcurl4-gnutls-dev_7.64.0-3_amd64.deb \
	file:${GNUTLS_DEV_BDIR}share/aclocal/libcurl.m4 \
	file:${GNUTLS_DOC}NEWS.Debian.gz \
	file:${GNUTLS_DOC}changelog.Debian.gz \
	file:${GNUTLS_DOC}changelog.gz \
	file:${GNUTLS_DOC}copyright \
	file:${GNUTLS_DEV_BDIR}share/man/man1/curl-config.1.gz

# create libcurl4-nss-dev targets
NSS_DEV_DIR=curl/debian/libcurl4-nss-dev/usr/
qmstrctl create file:${NSS_DEV_DIR}bin/curl-config

H_FILES=${NSS_DEV_DIR}include/x86_64-linux-gnu/curl/
qmstrctl create file:${H_FILES}curl.h
qmstrctl create file:${H_FILES}curlver.h
qmstrctl create file:${H_FILES}easy.h
qmstrctl create file:${H_FILES}mprintf.h
qmstrctl create file:${H_FILES}multi.h
qmstrctl create file:${H_FILES}stdcheaders.h
qmstrctl create file:${H_FILES}system.h
qmstrctl create file:${H_FILES}typecheck-gcc.h
qmstrctl create file:${H_FILES}urlapi.h

LIB_NSS_DEV=${NSS_DEV_DIR}lib/x86_64-linux-gnu/
qmstrctl create file:${LIB_NSS_DEV}libcurl-nss.la
qmstrctl create file:${LIB_NSS_DEV}pkgconfig/libcurl.pc

qmstrctl create file:${NSS_DEV_DIR}share/aclocal/libcurl.m4
qmstrctl create file:${NSS_DEV_DIR}share/man/man1/curl-config.1.gz

NSS_DEV_DOC=${NSS_DEV_DIR}share/doc/libcurl4-nss-dev/
qmstrctl create file:${NSS_DEV_DOC}NEWS.Debian.gz
qmstrctl create file:${NSS_DEV_DOC}changelog.Debian.gz
qmstrctl create file:${NSS_DEV_DOC}changelog.gz
qmstrctl create file:${NSS_DEV_DOC}copyright


# connect targets to libcurl4-nss-dev package
qmstrctl connect package:libcurl4-nss-dev_7.64.0-3_amd64.deb \
	file:${NSS_DEV_DIR}bin/curl-config \
	file:${H_FILES}curl.h \
	file:${H_FILES}curlver.h \
	file:${H_FILES}easy.h \
	file:${H_FILES}mprintf.h \
	file:${H_FILES}multi.h \
	file:${H_FILES}stdcheaders.h \
	file:${H_FILES}system.h \
	file:${H_FILES}typecheck-gcc.h \
	file:${H_FILES}urlapi.h

qmstrctl connect package:libcurl4-nss-dev_7.64.0-3_amd64.deb \
	file:hash:$(hashthis ${LIB_NSS_DEV}libcurl-nss.a) \
	file:${LIB_NSS_DEV}libcurl-nss.la \
	file:${LIB_NSS_DEV}pkgconfig/libcurl.pc

qmstrctl connect package:libcurl4-nss-dev_7.64.0-3_amd64.deb \
	file:${NSS_DEV_DIR}share/aclocal/libcurl.m4 \
	file:${NSS_DEV_DOC}NEWS.Debian.gz \
	file:${NSS_DEV_DOC}changelog.Debian.gz \
	file:${NSS_DEV_DOC}changelog.gz \
	file:${NSS_DEV_DOC}copyright \
	file:${NSS_DEV_DIR}share/man/man1/curl-config.1.gz


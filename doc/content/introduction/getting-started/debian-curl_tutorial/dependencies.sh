#!/bin/sh
set -e

echo "Downloading Debian cURL dependencies"
# Make sure to run as root
sudo apt-get update && apt-get install -y dpkg-dev debhelper libtool pkgconf \
            libssh2-1-dev python dh-exec groff-base libgnutls28-dev \
            libidn2-0-dev libkrb5-dev libldap2-dev libnghttp2-dev libnss3-dev \
            libpsl-dev librtmp-dev openssh-server quilt stunnel4 \
            libssl1.0-dev

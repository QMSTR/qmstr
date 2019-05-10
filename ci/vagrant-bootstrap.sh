#!/bin/sh
echo "Updating APT repositories..."
apt-get update -qq
echo "Installing protobuf compiler..."
apt-get -yqq install protobuf-compiler
echo "Installing build enviroment..."
apt-get -yqq install make autoconf build-essential libtool git
echo "Installing JSON-C dependencies..."
apt-get -yqq build-dep json-c
echo "Installing Docker..."
apt-get -yqq install docker.io
echo "Adding vagrant user to the docker group..."
addgroup vagrant docker
echo "Installing golang..."
export $(grep -v '^#' /vagrant/versions.env | xargs)
cd /opt
curl -s -o /opt/go.tar.gz https://dl.google.com/go/go${GO_VERSION}.linux-amd64.tar.gz
tar -xf go.tar.gz
export PATH=/opt/go/bin:$PATH
echo export PATH=$PATH >> /etc/bash.bashrc
echo "Installing shelldoc..."
git clone https://github.com/endocode/shelldoc
cd shelldoc
make
mv -f cmd/shelldoc/shelldoc /usr/local/bin/
cd -
echo "Done."

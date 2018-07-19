#!/bin/sh
echo "Updating APT repositories..."
apt-get update -qq
echo "Installing golang..."
apt-get -yqq install golang
echo "Installing protobuf compiler..."
apt-get -yqq install protobuf-compiler
echo "Installing Go protobuf interface..."
apt-get -yqq install golang-goprotobuf-dev
echo "Installing shelldoc..."
go get -u github.com/endocode/shelldoc/cmd/shelldoc && mv -f ~/go/bin/shelldoc /usr/local/bin/
echo "Done."

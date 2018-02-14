FROM golang:1.9 as builder

# install build deps
RUN apt-get update && \
    apt-get install -y git protobuf-compiler && \
    rm -rf /var/lib/apt/lists/*

# install golang tools
RUN go get -u github.com/golang/protobuf/protoc-gen-go && \
    go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p ${GOPATH}/src/github.com/QMSTR/qmstr && \
    mkdir ${GOPATH}/src/github.com/QMSTR/qmstr/cmd && \
    mkdir ${GOPATH}/src/github.com/QMSTR/qmstr/pkg
ADD cmd ${GOPATH}/src/github.com/QMSTR/qmstr/cmd
ADD pkg ${GOPATH}/src/github.com/QMSTR/qmstr/pkg
COPY Gopkg.* ${GOPATH}/src/github.com/QMSTR/qmstr/

RUN cd $GOPATH/src/github.com/QMSTR/qmstr && \
    dep ensure
RUN go generate github.com/QMSTR/qmstr/cmd/qmstr-master
RUN cd $GOPATH/src/github.com/QMSTR/qmstr && \
    go test ./...
RUN go install github.com/QMSTR/qmstr/cmd/qmstr-master && \
    go install github.com/QMSTR/qmstr/cmd/qmstr-wrapper && \
    go install github.com/QMSTR/qmstr/cmd/qmstr-cli

FROM ubuntu:17.10

# install runtime deps
RUN apt-get update && \
    apt-get install -y software-properties-common curl autoconf git libio-captureoutput-perl python python-pip protobuf-compiler && \
    rm -rf /var/lib/apt/lists/*

# install ninka
RUN mkdir /ninka && \
    git clone https://github.com/dmgerman/ninka.git /ninka && \
    cd /ninka/comments && make && make install && \
    rm /usr/local/man/man1 && \
    cd /ninka && perl Makefile.PL && make && make install && \
    rm -fr /ninka

# install scancode
RUN pip install scancode-toolkit

# install dgraph
RUN curl https://get.dgraph.io -sSf | bash

# copy qmstr installation from previous stage
RUN mkdir -p /go/bin
COPY --from=builder /go/bin/qmstr-master /go/bin/qmstr-master

EXPOSE 50051
# dgraph ratel web ui
EXPOSE 8000

COPY ci/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT [ "/docker-entrypoint.sh" ]

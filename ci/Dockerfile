# Create a container that only builds the software to be installed into the master container:
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
RUN go install github.com/QMSTR/qmstr/cmd/qmstr-master && \
    go install github.com/QMSTR/qmstr/cmd/qmstr-wrapper && \
    go install github.com/QMSTR/qmstr/cmd/qmstr-cli

# separate stage just to run the master unit tests:
FROM builder as master_unit_tests
CMD cd $GOPATH/src/github.com/QMSTR/qmstr && go test ./...

# the runtime stage contains all the elements needed to run the master and the analysis tools:
FROM ubuntu:17.10 as runtime

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
COPY ci/scancode.requirements /tmp/scancode.requirements
RUN pip install -r /tmp/scancode.requirements
RUN rm /tmp/scancode.requirements

# install dgraph
RUN curl https://get.dgraph.io -sSf | bash

# copy qmstr installation from previous stage
RUN mkdir -p /go/bin
COPY --from=builder /go/bin/qmstr-master /go/bin/qmstr-master

EXPOSE 50051
# dgraph ratel web ui
EXPOSE 8000

VOLUME /buildroot

# release master container, based on the runtime stage:
FROM runtime as master
COPY ci/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh
ENTRYPOINT [ "/docker-entrypoint.sh" ]

# development container for the master and other tools, based in the runtime stage:
# The $GOROOT/src directory can be passed in as a volume, to allow for testing local changes.
FROM runtime as dev
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:/usr/lib/go-1.9/bin:$PATH

# install golang 1.9
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:gophers/archive && \
    apt-get update && \
    apt-get install -y curl golang-1.9-go autoconf git libio-captureoutput-perl python python-pip protobuf-compiler

EXPOSE 2345

# install go deps
RUN go get -u -v github.com/golang/protobuf/protoc-gen-go && \
    go get -v github.com/dgraph-io/dgo && \
    go get -u -v github.com/derekparker/delve/cmd/dlv && \
    go get -v github.com/spf13/pflag

RUN go get -v github.com/magefile/mage
RUN go get -d -v github.com/gohugoio/hugo
RUN (cd $GOPATH/src/github.com/gohugoio/hugo && mage vendor && mage install)

VOLUME /go/src

COPY ci/dev-entrypoint.sh /dev-entrypoint.sh
RUN chmod +x /dev-entrypoint.sh
ENTRYPOINT [ "/dev-entrypoint.sh" ]

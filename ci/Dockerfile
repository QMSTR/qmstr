# Versioning
ARG DGRAPH_VERSION=1.0.5
ARG HUGO_VERSION=0.39
ARG SCANCODE_VERSION=2.9.1

# Create a container that only builds the software to be installed into the master container:
FROM ubuntu:18.04 as builder
LABEL org.qmstr.image builder
ENV GOPATH /go
ENV PATH ${GOPATH}/bin:/usr/lib/go-1.9/bin:$PATH
# install golang 1.9
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:gophers/archive && \
    apt-get update && \
    apt-get install -y curl golang-1.9-go autoconf git libgit2-dev libio-captureoutput-perl python-dev python-virtualenv protobuf-compiler

# install build deps
RUN apt-get update && \
    apt-get install -y git protobuf-compiler python-virtualenv python-dev && \
    rm -rf /var/lib/apt/lists/*

# install golang tools
RUN mkdir -p ${GOPATH}/src/github.com/QMSTR/qmstr && \
    mkdir ${GOPATH}/src/github.com/QMSTR/qmstr/cmd && \
    mkdir ${GOPATH}/src/github.com/QMSTR/qmstr/pkg
ADD cmd ${GOPATH}/src/github.com/QMSTR/qmstr/cmd
ADD pkg ${GOPATH}/src/github.com/QMSTR/qmstr/pkg
ADD proto ${GOPATH}/src/github.com/QMSTR/qmstr/proto
ADD python ${GOPATH}/src/github.com/QMSTR/qmstr/python
ADD Makefile ${GOPATH}/src/github.com/QMSTR/qmstr/Makefile
COPY Gopkg.* ${GOPATH}/src/github.com/QMSTR/qmstr/

WORKDIR ${GOPATH}/src/github.com/QMSTR/qmstr
RUN make install_qmstr_all

# separate stage just to run the master unit tests:
FROM builder as master_unit_tests

LABEL org.qmstr.image master_unit_tests

WORKDIR ${GOPATH}/src/github.com/QMSTR/qmstr
CMD make gotest

# the runtime stage contains all the elements needed to run the master and the analysis tools:
FROM ubuntu:18.04 as runtime

LABEL org.qmstr.image runtime

# install runtime deps
RUN apt-get update && \
    apt-get install -y software-properties-common curl autoconf git libgit2-dev libio-captureoutput-perl python python-pip protobuf-compiler && \
    rm -rf /var/lib/apt/lists/*

ARG HUGO_VERSION
RUN curl -L --output /tmp/hugo.deb https://github.com/gohugoio/hugo/releases/download/v${HUGO_VERSION}/hugo_${HUGO_VERSION}_Linux-64bit.deb
RUN dpkg -i /tmp/hugo.deb
RUN rm /tmp/hugo.deb

# init html reporter data
COPY --from=builder /go/src/github.com/QMSTR/qmstr /tmp/qmstr
RUN /tmp/qmstr/cmd/reporters/qmstr-reporter-html/setup.sh /usr/share/qmstr /tmp/qmstr
RUN rm -fr /tmp/qmstr

# install ninka
RUN mkdir /ninka && \
    git clone https://github.com/dmgerman/ninka.git /ninka && \
    cd /ninka/comments && make && make install && \
    rm /usr/local/man/man1 && \
    cd /ninka && perl Makefile.PL && make && make install && \
    rm -fr /ninka

# install scancode
ARG SCANCODE_VERSION
RUN pip install -U scancode-toolkit==${SCANCODE_VERSION}

# install dgraph
ARG DGRAPH_VERSION
RUN curl -L --output /tmp/dgraph.sha256 https://github.com/dgraph-io/dgraph/releases/download/v"${DGRAPH_VERSION}"/dgraph-checksum-linux-amd64.sha256
RUN curl -L --output /tmp/dgraph.tar.gz https://github.com/dgraph-io/dgraph/releases/download/v"${DGRAPH_VERSION}"/dgraph-linux-amd64.tar.gz
RUN temp_dir=$(mktemp -d 2>/dev/null) && \
    tar -C $temp_dir -xzf /tmp/dgraph.tar.gz && \
    sed -e s#/usr/local/bin#${temp_dir}# /tmp/dgraph.sha256 | sha256sum -c - && \
    mv ${temp_dir}/* /usr/local/bin/ && \
    rmdir ${temp_dir}
RUN rm /tmp/dgraph*

# copy qmstr installation from previous stage
RUN mkdir -p /go/bin
COPY --from=builder /usr/local/bin/* /usr/local/bin/ 

EXPOSE 50051
# dgraph ratel web ui
EXPOSE 8000

VOLUME /buildroot

# release master container, based on the runtime stage:
FROM runtime as master

LABEL org.qmstr.image master 

ENV GOPATH /go
ENV PATH ${GOPATH}/bin:$PATH

COPY ci/docker-entrypoint.sh /docker-entrypoint.sh
RUN chmod +x /docker-entrypoint.sh

ENTRYPOINT [ "/docker-entrypoint.sh" ]

# development container for the master and other tools, based on the runtime stage:
# The $GOROOT/src directory can be mounted in as a volume, to allow testing of local changes.
FROM runtime as dev

LABEL org.qmstr.image dev 

ENV GOPATH /go
ENV PATH ${GOPATH}/bin:/usr/lib/go-1.9/bin:$PATH

# install golang 1.9
RUN apt-get update && \
    apt-get install -y software-properties-common && \
    add-apt-repository ppa:gophers/archive && \
    apt-get update && \
    apt-get install -y curl golang-1.9-go autoconf git libio-captureoutput-perl python-dev python-virtualenv protobuf-compiler

EXPOSE 2345

# install additional go deps
RUN go get -u -v github.com/derekparker/delve/cmd/dlv

VOLUME /go/src

# cache wheels
COPY --from=builder /root/.pex/build /root/.pex/build
RUN rm /root/.pex/build/*qmstr*

COPY ci/dev-entrypoint.sh /dev-entrypoint.sh
RUN chmod +x /dev-entrypoint.sh
ENTRYPOINT [ "/dev-entrypoint.sh" ]

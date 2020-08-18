################################################################################
# STAGE 1: Build master binaries                                               #
################################################################################
ARG GOBUILDER_IMAGE=golang:1.12-buster

FROM $GOBUILDER_IMAGE as gobuilder
ENV GOPROXY="https://proxy.golang.org"

RUN set -e && \
  apt update && \
  apt install -y protobuf-compiler && \
  mkdir /root/qmstr/

WORKDIR /root/qmstr

COPY ./ ./

RUN set -e && \
  go build -o ./out/qmstr-master ./masterserver/

################################################################################
# STAGE 2: Build master container image (deploy)                               #
################################################################################

FROM debian:buster-slim as deploy

# Required QMSTR directories
ENV QMSTR_DIRS="/var/qmstr/ /var/cache/qmstr/ /var/lib/qmstr/"

# Copy binaries from build stage
COPY --from=gobuilder /root/qmstr/out/* /usr/local/bin/

RUN set -e && \
  addgroup qmstr && \
  adduser --system qmstr --ingroup qmstr && \
  mkdir -p ${QMSTR_DIRS} && \
  chown -R qmstr ${QMSTR_DIRS}

WORKDIR /home/qmstr
USER qmstr

EXPOSE 50051

VOLUME /home/qmstr/config
VOLUME /var/qmstr/buildroot

ENTRYPOINT ["/usr/local/bin/qmstr-master"]

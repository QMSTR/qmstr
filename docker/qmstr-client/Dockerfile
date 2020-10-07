################################################################################
# STAGE 1a: Build client binaries                                              #
################################################################################
ARG GOBUILDER_IMAGE=golang:1.12-buster

FROM $GOBUILDER_IMAGE AS gobuilder
ENV GOPROXY="https://proxy.golang.org"

RUN set -e && \
  apt update && \
  apt install -y protobuf-compiler && \
  mkdir /root/qmstr/

WORKDIR /root/qmstr

COPY ./ ./

RUN set -e && \
  go test ./clients/qmstr && \
  go build -o ./out/qmstr ./clients/qmstr && \
  go test ./clients/qmstrctl && \
  go build -o ./out/qmstrctl ./clients/qmstrctl

################################################################################
# STAGE 1b: Build base container image                                         #
################################################################################

FROM debian:buster-slim AS base
ARG UID=1000

COPY --from=gobuilder /root/qmstr/out/* /usr/local/bin/

RUN set -e && \
  mkdir -p /var/qmstr/ && \
  addgroup qmstrclient && \
  adduser --system qmstrclient --ingroup qmstrclient --uid $UID && \
  chown -R qmstrclient:qmstrclient /var/qmstr/

WORKDIR /home/qmstrclient
USER qmstrclient

VOLUME /var/qmstr/buildroot

ENTRYPOINT qmstrctl version && qmstr --help

################################################################################
# STAGE 2a: Build Maven plugin dependencies                                    #
################################################################################

FROM openjdk:11-slim-buster AS javabuilder

COPY ./ ./qmstr

RUN set -e && \
  mkdir -p /usr/share/man/man1 && \
  apt update && \
  apt install -y --no-install-recommends maven && \
  cd ./qmstr/lib/java-qmstr && \
  ./gradlew install && \
  cd ../../modules/builders/qmstr-maven-plugin && \
  mvn install

################################################################################
# STAGE 2b: Build client container image w/ Maven plugin                       #
################################################################################

FROM base AS mvn

ENV M2_HOME /maven
USER root

RUN set -e && \
  mkdir -p /usr/share/man/man1 && \
  apt update  && \
  apt install -y --no-install-recommends openjdk-11-jdk openjfx maven && \
  mkdir -p ${M2_HOME}/conf

ENV JAVA_HOME="/usr/lib/jvm/java-1.11.0-openjdk-amd64"
ENV PATH="${JAVA_HOME}/bin:${PATH}"

ADD ./docker/qmstr-client/settings.xml /usr/share/maven/conf/settings.xml

COPY --from=javabuilder --chown=qmstrclient /root/.m2/repository ${M2_HOME}/repo

WORKDIR /var/qmstr/buildroot/project

USER qmstrclient

ENTRYPOINT mvn package

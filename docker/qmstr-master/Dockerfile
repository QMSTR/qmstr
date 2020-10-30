# Scancode version and info
ARG SCANCODE_VERSION="3.2.3"
ARG SCANCODE_RELEASE_CANDIDATE=""
ARG SCANCODE_PROJECT_NAME="scancode-toolkit"
ARG SCANCODE_INSTALLATION_FOLDER="scancode"

################################################################################
# STAGE 1a: Build master binaries                                              #
################################################################################

FROM golang:1.12-buster as gobuilder
ENV GOPROXY="https://proxy.golang.org"

RUN set -e && \
  apt update && \
  apt install -y protobuf-compiler && \
  mkdir /root/qmstr/

WORKDIR /root/qmstr

COPY ./ ./

RUN set -e && \
  go build -o ./out/scancode-analyzer ./modules/analyzers/scancode-analyzer && \
  go build -o ./out/qmstr-master ./masterserver/

################################################################################
# STAGE 1b: Downloading scancode                                               #
################################################################################

FROM alpine as scancode

ARG SCANCODE_VERSION
ARG SCANCODE_RELEASE_CANDIDATE
ARG SCANCODE_PROJECT_NAME
ARG SCANCODE_INSTALLATION_FOLDER
ENV SCANCODE_ZIP_LINK="https://github.com/nexB/${SCANCODE_PROJECT_NAME}/releases/download/v${SCANCODE_VERSION}${SCANCODE_RELEASE_CANDIDATE}/${SCANCODE_PROJECT_NAME}-${SCANCODE_VERSION}${SCANCODE_RELEASE_CANDIDATE}.zip"
ENV SCANCODE_ZIP_NAME="scancode.zip"

RUN wget ${SCANCODE_ZIP_LINK} --output-document ${SCANCODE_ZIP_NAME} && \
  mkdir ${SCANCODE_INSTALLATION_FOLDER} && \
  unzip -q ${SCANCODE_ZIP_NAME} -d ${SCANCODE_INSTALLATION_FOLDER}

################################################################################
# STAGE 2: Build master container image (deploy)                               #
################################################################################

FROM python:3.6-slim-buster as deploy

# Required QMSTR directories
ENV QMSTR_DIRS="/var/qmstr/ /var/cache/qmstr/ /var/lib/qmstr/"

# Copy binaries from build stage
COPY --from=gobuilder /root/qmstr/out/* /usr/local/bin/

# Copying scancode binary
ARG SCANCODE_VERSION
ARG SCANCODE_RELEASE_CANDIDATE
ARG SCANCODE_PROJECT_NAME
ARG SCANCODE_INSTALLATION_FOLDER
ENV SCANCODE_LOCAL_DIRECTORY="/usr/local/bin/scancode/"
RUN mkdir -p ${SCANCODE_LOCAL_DIRECTORY}
COPY --from=scancode /${SCANCODE_INSTALLATION_FOLDER}/${SCANCODE_PROJECT_NAME}-${SCANCODE_VERSION}${SCANCODE_RELEASE_CANDIDATE}/ ${SCANCODE_LOCAL_DIRECTORY}
ENV PATH="${SCANCODE_LOCAL_DIRECTORY}:${PATH}"

# Installing scancode dependencies
RUN apt update && apt install -y libgomp1 libbz2-1.0 xz-utils zlib1g libxml2-dev libxslt1-dev

RUN set -e && \
  addgroup qmstr && \
  adduser --system qmstr --ingroup qmstr && \
  mkdir -p ${QMSTR_DIRS} && \
  chown -R qmstr ${QMSTR_DIRS}

# Configuring scancode (must be launched at least once)
RUN chown -R qmstr ${SCANCODE_LOCAL_DIRECTORY} && \
  scancode --help

WORKDIR /home/qmstr
USER qmstr

EXPOSE 50051

VOLUME /home/qmstr/config
VOLUME /var/qmstr/buildroot

ENTRYPOINT ["/usr/local/bin/qmstr-master"]

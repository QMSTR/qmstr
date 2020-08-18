FROM debian:latest

ARG UID=1000
RUN set -e && \
  apt update && \
  apt install -y git && \
  addgroup git && \
  adduser --system git --ingroup git --uid $UID

WORKDIR /home/git
USER git

VOLUME /home/git/buildroot

ENTRYPOINT [ "git", "clone" ]

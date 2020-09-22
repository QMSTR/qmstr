FROM python:latest
ARG UID=1000

RUN set -e && \
  addgroup configure && \
  adduser --uid $UID --system configure --ingroup configure

WORKDIR /home/configure
USER configure

COPY ./lib/pyqmstr/pom-patch/ ./

RUN set -e && \
  pip install -r ./requirements.txt

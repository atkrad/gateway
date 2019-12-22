FROM haproxy:2.1.1-alpine

ENV WAIT4X_VERSION=v0.1.0

LABEL maintainer="mohammad.abdolirad@snapp.cab"

RUN apk --no-cache --update add git bash curl rsync openssh-client openssl \
    && curl -SLO https://github.com/atkrad/wait4x/releases/download/$WAIT4X_VERSION/wait4x-linux-amd64 \
    && mv wait4x-linux-amd64 /usr/local/bin/wait4x \
    && chmod a+x /usr/local/bin/wait4x

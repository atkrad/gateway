FROM haproxy:2.1.1-alpine

ENV HAPROXY_DATA_PLANE_API_VERSION=v1.2.4

LABEL maintainer="m.abdolirad@gmail.com"

RUN apk --no-cache --update add curl \
    && curl -SLO https://github.com/haproxytech/dataplaneapi/releases/download/$HAPROXY_DATA_PLANE_API_VERSION/dataplaneapi \
    && mv dataplaneapi /usr/local/bin/dataplaneapi \
    && chmod a+x /usr/local/bin/dataplaneapi

COPY haproxy.cfg /usr/local/etc/haproxy/haproxy.cfg
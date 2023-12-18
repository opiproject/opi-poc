# SPDX-License-Identifier: Apache-2.0
# Copyright (C) 2023 Network Plumping Working Group
# Copyright (C) 2023 Nordix Foundation.

FROM golang:1.20-alpine3.19 as builder

COPY . /usr/src/evpn-gw-cni

ENV HTTP_PROXY $http_proxy
ENV HTTPS_PROXY $https_proxy

WORKDIR /usr/src/evpn-gw-cni
RUN apk add --no-cache --virtual build-dependencies build-base=~0.5 && \
    make clean && \
    make build

FROM alpine:3.19
COPY --from=builder /usr/src/evpn-gw-cni/build/evpn-gw /usr/bin/
WORKDIR /

LABEL io.k8s.display-name="EVPN GW CNI"

COPY ./images/entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
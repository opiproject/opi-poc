FROM golang:alpine as builder

COPY . /usr/src/xpu-cni

ENV HTTP_PROXY $http_proxy
ENV HTTPS_PROXY $https_proxy

WORKDIR /usr/src/xpu-cni
RUN apk add --no-cache --virtual build-dependencies build-base=~0.5 && \
    make clean && \
    make build

FROM alpine:3
COPY --from=builder /usr/src/xpu-cni/build/xpu /usr/bin/
WORKDIR /

LABEL io.k8s.display-name="XPU CNI"

COPY ./images/entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]
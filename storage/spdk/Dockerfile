# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Dell Inc, or its subsidiaries.

FROM docker.io/library/fedora:36 as build

ARG TAG=v22.09
# Pick an arch that has at least sse 4.2 but does not require newer avx
# See https://gcc.gnu.org/onlinedocs/gcc/x86-Options.html
ARG ARCH=x86-64-v2

WORKDIR /root
RUN dnf install -y git rpm-build diffutils procps-ng pip && dnf clean all

# hadolint ignore=DL3003
RUN git clone https://github.com/spdk/spdk --branch ${TAG} --depth 1 && \
    cd spdk && git submodule update --init --depth 1 && scripts/pkgdep.sh --rdma

# hadolint ignore=DL3003
RUN cd spdk && ./rpmbuild/rpm.sh --target-arch=${ARCH} --without-uring --without-crypto \
    --without-fio --with-raid5f --with-vhost --without-pmdk --without-rbd \
    --with-rdma --with-shared --with-iscsi-initiator --without-vtune

FROM docker.io/library/fedora:36

WORKDIR /root
RUN mkdir -p /root/rpmbuild
COPY --from=build /root/rpmbuild/ /root/rpmbuild/
RUN dnf install -y /root/rpmbuild/rpm/x86_64/*.rpm && dnf clean all
EXPOSE 9009
HEALTHCHECK CMD curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{\"id\": 1, \"method\": \"bdev_get_bdevs\"}' http://localhost:9009 || exit 1
CMD ["/usr/libexec/spdk/scripts/rpc_http_proxy.py", "0.0.0.0", "9009", "spdkuser", "spdkpass"]

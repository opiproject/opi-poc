FROM fedora:36 as build

ARG TAG=v22.05
ARG ARCH=native

WORKDIR /root
RUN dnf install -y git rpm-build diffutils procps-ng && dnf clean all

# hadolint ignore=DL3003
RUN git clone https://github.com/spdk/spdk --branch ${TAG} --depth 1 && \
    cd spdk && git submodule update --init --depth 1 && scripts/pkgdep.sh --rdma

# hadolint ignore=DL3003
RUN cd spdk && CFLAGS=-mno-avx512f ./rpmbuild/rpm.sh --without-uring --without-crypto \
    --without-fio --with-raid5 --with-vhost --without-pmdk --without-rbd \
    --with-rdma --with-shared --with-iscsi-initiator --without-vtune --without-isal

FROM fedora:36

WORKDIR /root
RUN mkdir -p /root/rpmbuild
COPY --from=build /root/rpmbuild/ /root/rpmbuild/
RUN dnf install -y /root/rpmbuild/rpm/x86_64/*.rpm && dnf clean all
EXPOSE 9009
HEALTHCHECK CMD curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{\"id\": 1, \"method\": \"bdev_get_bdevs\"}' http://localhost:9009 || exit 1
CMD ["/usr/libexec/spdk/scripts/rpc_http_proxy.py", "0.0.0.0", "9009", "spdkuser", "spdkpass"]

# Dockerfile

FROM ubuntu:24.04 AS builder

ENV DEBIAN_FRONTEND=noninteractive

# Install iproute2 and other networking tools
# hadolint ignore=DL3008
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
			iproute2 iputils-ping net-tools git tar wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /usr/local
RUN wget -q https://go.dev/dl/go1.25.1.linux-arm64.tar.gz && \
    tar -xvf go1.25.1.linux-arm64.tar.gz && \
    rm go1.25.1.linux-arm64.tar.gz
ENV PATH=$PATH:/usr/local/go/bin/
RUN go version

# Clone OPI Strongswan Bridge
WORKDIR /build
RUN git clone https://github.com/opiproject/opi-strongswan-bridge.git
WORKDIR /build/opi-strongswan-bridge

# Build OPI Strongswan Bridge
RUN go mod download && \
    go build -v -o /opi-vici-bridge ./cmd/...

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# Install iproute2 and other networking tools
# hadolint ignore=DL3008
RUN apt-get update -y && \
    apt-get install -y --no-install-recommends \
			iproute2 iputils-ping net-tools git tar \
			strongswan strongswan-pki libcharon-extra-plugins \
			libcharon-extauth-plugins libstrongswan-extra-plugins \
			strongswan-swanctl vim tcpdump ethtool && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*


COPY --from=builder /opi-vici-bridge /
COPY --from=docker.io/fullstorydev/grpcurl:v1.8.9-alpine /bin/grpcurl /usr/local/bin/
EXPOSE 50051

# Copy Strongswan config files and setup
RUN rm -rf /etc/swanctl /etc/strongswan.conf /etc/strongswan.d
COPY strongswan.conf /etc/strongswan.conf
COPY strongswan.d /etc/strongswan.d
COPY swanctl /etc/swanctl
COPY run.sh /run.sh

# Set up Strongswan
ENTRYPOINT ["bash", "/run.sh"]

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# Install iproute2 and other networking tools
# hadolint ignore=DL3008
RUN apt-get update -y && \
	apt-get install -y --no-install-recommends \
			iproute2 iputils-ping net-tools git tar \
			vim tcpdump && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY run.sh /
ENTRYPOINT ["bash", "/run.sh"]

FROM ubuntu:24.04

ENV DEBIAN_FRONTEND=noninteractive

# Install iproute2 and other networking tools
# hadolint ignore=DL3008
RUN apt-get update -y && \
	apt-get install -y --no-install-recommends \
			iproute2 iputils-ping net-tools git tar gettext-base \
			strongswan strongswan-pki libcharon-extra-plugins \
			libcharon-extauth-plugins libstrongswan-extra-plugins \
			strongswan-swanctl vim tcpdump && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN rm -rf /etc/swanctl /etc/strongswan.conf /etc/strongswan.d
COPY strongswan.conf /etc/strongswan.conf
COPY strongswan.d /etc/strongswan.d
COPY swanctl /etc/swanctl

RUN ipsec restart

COPY run.sh /run.sh
CMD [ "bash", "/run.sh" ]

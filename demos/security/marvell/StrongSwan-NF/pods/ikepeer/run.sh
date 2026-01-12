#!/bin/bash
set -e

# Set default values for environment variables
export HOST_SUBNET="${HOST_SUBNET:-10.56.217.0/24}"
export EXTERNAL_SUBNET="${EXTERNAL_SUBNET:-20.20.20.0/24}"

echo "Configuration:"
echo "  HOST_SUBNET: ${HOST_SUBNET}"
echo "  EXTERNAL_SUBNET: ${EXTERNAL_SUBNET}"

# Process swanctl configuration templates
echo "Processing configuration templates..."
if [ -f /etc/swanctl/swanctl.conf.template ]; then
    envsubst < /etc/swanctl/swanctl.conf.template > /etc/swanctl/swanctl.conf
fi
if [ -f /etc/swanctl/psk.conf.template ]; then
    envsubst < /etc/swanctl/psk.conf.template > /etc/swanctl/psk.conf
fi
if [ -f /etc/swanctl/eap.conf.template ]; then
    envsubst < /etc/swanctl/eap.conf.template > /etc/swanctl/eap.conf
fi
if [ -f /etc/swanctl/rw.conf.template ]; then
    envsubst < /etc/swanctl/rw.conf.template > /etc/swanctl/rw.conf
fi

echo "Starting IPsec..."
ipsec restart || { echo "Failed to start ipsec"; exit 1; }

echo "IPsec started successfully. Container running..."
sleep infinity

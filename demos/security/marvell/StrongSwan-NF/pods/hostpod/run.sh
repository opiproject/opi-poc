#!/bin/bash
set -e

# Set default values for environment variables
export EXTERNAL_SUBNET="${EXTERNAL_SUBNET:-20.20.20.0/24}"
export HOST_GATEWAY="${HOST_GATEWAY:-10.56.217.1}"

echo "Configuration:"
echo "  EXTERNAL_SUBNET: ${EXTERNAL_SUBNET}"
echo "  HOST_GATEWAY: ${HOST_GATEWAY}"

echo "Adding route to IPsec network..."
ip r add "${EXTERNAL_SUBNET}" via "${HOST_GATEWAY}" || echo "Warning: Route may already exist"

echo "Host pod ready. Container running..."
sleep infinity

#!/bin/bash
set -e

# Set default values for environment variables
export NF_HOST_IP="${NF_HOST_IP:-10.56.217.1/24}"
export NF_EXTERNAL_IP="${NF_EXTERNAL_IP:-20.20.20.1/24}"

echo "Configuration:"
echo "  NF_HOST_IP: ${NF_HOST_IP}"
echo "  NF_EXTERNAL_IP: ${NF_EXTERNAL_IP}"

# Configure network interfaces
echo "Configuring network interfaces..."
ifconfig net1 "${NF_HOST_IP}" || echo "Warning: net1 configuration failed"
ifconfig net2 "${NF_EXTERNAL_IP}" || echo "Warning: net2 configuration failed"

ethtool -K net1 rx off tx off
ethtool -K net2 rx off tx off

# Start IPsec
echo "Starting IPsec..."
ipsec restart

# Wait for VICI socket to be available
echo "Waiting for VICI socket..."
SOCKET_PATH="/var/run/charon.vici"
while [ ! -S "$SOCKET_PATH" ]; do
    sleep 0.5
done
echo "VICI socket is ready"

swanctl --load-all

# Start OPI StrongSwan Bridge
echo "Starting OPI StrongSwan Bridge on port 50151..."
/opi-vici-bridge -port=50151

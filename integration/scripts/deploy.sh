#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation
#
# This script is used to deploy the OPI PoC into one of two
# environments:
#
# * Developer emulated xPU setup
# * Split setup between an actual xPU and a host/VM/container
#
# This script assumes you have Docker installed on the xPU and host
# instances.
#

set -eo pipefail

# Enable debug if running in CI
if [ -n "${CI}" ] && [ "${CI}" == "true" ]
then
    set -x
fi

# DOCKER_COMPOSE setup
DC=docker-compose

if [ "$(command -v ${DC})" == "" ]
then
    DC="docker compose"
fi

usage() {
    echo ""
    echo "Usage: $0 -m [dev | xpu] -b [BMC IP address] -i [IP address of the host] -x [IP address of the xPU]"
    echo ""
    exit 1
}

deploy_dev() {
    echo "Deploying development emulated environment"

    # Run the integration script to start things up
    ./scripts/integration.sh start
}

deploy_xpu() {
    echo "Deploying xPU environment"
    echo "BMC IP address: ${BMC_IP_ADDRESS}"
    echo "Host IP address: ${HOST_IP_ADDRESS}"
    echo "xPU IP address: ${XPU_IP_ADDRESS}"

    # Deploy to host/VM/container
    export COMPOSE_FILE=docker-compose.pxe.yml:docker-compose.spdk.yml
    export DOCKER_HOST="ssh://user@${HOST_IP_ADDRESS}"
    bash -c "${DC} up -d"

    # Deploy to XPU
    export COMPOSE_FILE=docker-compose.xpu.yml
    export DOCKER_HOST="ssh://user@${XPU_IP_ADDRESS}"
    bash -c "${DC} up -d"
}

# Default mode is dev
# Valid values: "dev" or "xpu"
MODE=dev

# IP addresses, only used for xpu mode
XPU_IP_ADDRESS=
HOST_IP_ADDRESS=
BMC_IP_ADDRESS=

while getopts b:i:m:x: option
do
    case "${option}"
        in
        b)
            BMC_IP_ADDRESS="${OPTARG}"
            ;;
        i)
            HOST_IP_ADDRESS="${OPTARG}"
            ;;
        m)
            MODE="${OPTARG}"
            ;;
        x)
            XPU_IP_ADDRESS="${OPTARG}"
            ;;
        *)
            usage
            ;;
    esac
done

if [ "${MODE}" != "dev" ] && [ "${MODE}" != "xpu" ]
then
    usage
fi

echo "Selected mode ${MODE}"

if [ "${MODE}" == "xpu" ]
then
    deploy_xpu
elif [ "${MODE}" == "dev" ]
then
    deploy_dev
fi

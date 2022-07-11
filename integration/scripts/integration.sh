#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation

set -euxo pipefail

INT_BUILD=build
INT_BUILD_START=build-start
INT_RUN_TESTS=run-tests
INT_LOGS=logs
INT_STOP=stop

# DOCKER_COMPOSE setup
DC=docker-compose

if [ "$(which ${DC})" == "" ]
then
    DC="docker compose"
fi

usage() {
    echo ""
    echo "Usage: integration.sh [${INT_BUILD} | ${INT_BUILD_START} | ${INT_RUN_TESTS} | ${INT_LOGS} | ${INT_STOP}]"
    echo ""
}

build_containers() {
    bash -c "${DC} -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml build --parallel"
}

start_containers() {
    bash -c "${DC} -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml down"
    docker network prune --force
    bash -c "${DC} -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml up -d"
}

run_integration_tests() {
    bash -c "${DC} -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml ps"
    echo wait 5s... && sleep 5s
    curl --fail http://127.0.0.1:8001/redfish/v1/Systems/437XR1138R2
    curl --fail http://127.0.0.1:8002/redfish/v1/Systems/437XR1138R2
    curl --fail http://127.0.0.1:8082/var/lib/tftpboot/
    curl --fail http://127.0.0.1:9090/
    curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{"id": 1, "method": "bdev_get_bdevs"}' http://127.0.0.1:9004
    curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{"id": 1, "method": "bdev_get_bdevs"}' http://127.0.0.1:9009
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2210 host@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2208  bmc@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2207  xpu@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2209  bmc@127.0.0.1 hostname
    bash -c "${DC} -f docker-compose.pxe.yml exec -T pxe nmap --script broadcast-dhcp-discover"
    bash -c "${DC} -f docker-compose.pxe.yml exec -T pxe nmap --script broadcast-dhcp-discover" | grep "Server Identifier: 10.127.127.3"
    bash -c "${DC} -f docker-compose.pxe.yml exec -T pxe curl --fail http://10.127.127.3:8082/var/lib/tftpboot/"
    bash -c "${DC} -f docker-compose.pxe.yml exec -T pxe tftp 10.127.127.3 -v -c get grubx64.efi"
    bash -c "${DC} -f docker-compose.spdk.yml exec -T spdk-target /root/spdk/build/examples/identify -r 'traddr:127.0.0.1 trtype:TCP adrfam:IPv4 trsvcid:4420'"
    bash -c "${DC} -f docker-compose.xpu.yml  exec -T xpu-spdk /root/spdk/build/examples/identify -r 'traddr:10.127.127.4 trtype:TCP adrfam:IPv4 trsvcid:4420'"
}

acquire_logs() {
    bash -c "${DC} -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml ps -a"
    bash -c "${DC} -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml logs" || true
    netstat -an || true
    ifconfig -a || true
    docker inspect bash -c "${DC} compose -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml ps -aq" || true
}

stop_containers() {
    bash -c "${DC} -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml down"
}

if [ "$#" -lt 1 ]
then
    usage
    exit 1
fi

if [ "$1" == "${INT_BUILD}" ]
then
    build_containers
elif [ "$1" == "${INT_BUILD_START}" ]
then
    build_containers
    start_containers
elif [ "$1" == "${INT_RUN_TESTS}" ]
then
    run_integration_tests
elif [ "$1" == "${INT_LOGS}" ]
then
    acquire_logs
elif [ "$1" == "${INT_STOP}" ]
then
    stop_containers
else
    echo "Invalid argument: $1"
    usage
    exit 1
fi

#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation

set -euxo pipefail

INT_BUILD_START=build-start
INT_RUN_TESTS=run-tests
INT_LOGS=logs
INT_STOP=stop

usage() {
    echo ""
    echo "Usage: integration.sh [${INT_BUILD_START} | ${INT_RUN_TESTS} | ${INT_LOGS} | ${INT_STOP}]"
    echo ""
}

build_and_start_containers() {
    docker-compose -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml down
    docker network prune --force
    docker-compose -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml up -d
}

run_integration_tests() {
    docker-compose -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml ps
    echo wait 5s... && sleep 5s
    curl --fail http://127.0.0.1:8001/redfish/v1/Systems/437XR1138R2
    curl --fail http://127.0.0.1:8002/redfish/v1/Systems/437XR1138R2
    curl --fail http://127.0.0.1:8082/var/lib/tftpboot/
    curl --fail http://127.0.0.1:9090/
    curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{\"id\": 1, \"method\": \"bdev_get_bdevs\"}' http://127.0.0.1:9004 || true
    curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{\"id\": 1, \"method\": \"bdev_get_bdevs\"}' http://127.0.0.1:9009 || true
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2210 host@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2208  bmc@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2207  xpu@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2209  bmc@127.0.0.1 hostname
    docker-compose -f docker-compose.pxe.yml exec -T pxe dnf install -y nmap tftp
    docker-compose -f docker-compose.pxe.yml exec -T pxe nmap --script broadcast-dhcp-discover
    docker-compose -f docker-compose.pxe.yml exec -T pxe nmap --script broadcast-dhcp-discover | grep "Server Identifier: 10.127.127.3"
    docker-compose -f docker-compose.pxe.yml exec -T pxe curl --fail http://10.127.127.3:8082/var/lib/tftpboot/
    docker-compose -f docker-compose.pxe.yml exec -T pxe tftp 10.127.127.3 -v -c get grubx64.efi
}

acquire_logs() {
    docker-compose logs || true
    netstat -an || true
    ifconfig -a || true
    docker inspect "$(docker-compose -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml ps -aq)" || true
}

stop_containers() {
    docker-compose -f docker-compose.yml -f docker-compose.xpu.yml -f docker-compose.otel.yml -f docker-compose.spdk.yml -f docker-compose.pxe.yml down
}

if [ "$#" -lt 1 ]
then
    usage
    exit 1
fi

if [ "$1" == "${INT_BUILD_START}" ]
then
    build_and_start_containers
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

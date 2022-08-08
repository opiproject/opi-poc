#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation

set -euxo pipefail

INT_BUILD=build
INT_START=start
INT_TESTS=tests
INT_LOGS=logs
INT_STOP=stop

# DOCKER_COMPOSE setup
DC=docker-compose

if [ "$(command -v ${DC})" == "" ]
then
    DC="docker compose"
fi

usage() {
    echo ""
    echo "Usage: integration.sh [${INT_BUILD} | ${INT_START} | ${INT_TESTS} | ${INT_LOGS} | ${INT_STOP}]"
    echo ""
}

build_containers() {
    bash -c "${DC} build --parallel"
}

start_containers() {
    bash -c "${DC} down"
    docker network prune --force
    bash -c "${DC} pull"
    # Workaround for running on servers without AVX512
    if [ -n "${BUILD_SPDK:-}" ]; then
        bash -c "${DC} build spdk-target"
    fi
    bash -c "${DC} up -d"
}

run_integration_tests() {
    bash -c "${DC} ps"
    # shellcheck disable=SC2046
    uniq -c <<< "$(sort <<< "$(docker inspect --format='{{json .State.Health.Status}}' $(${DC} ps -q))")"

    # TODO: replace sleep with timeouted-wait for all services to become healthy
    echo wait 5s... && sleep 5s
    curl --fail http://127.0.0.1:8001/redfish/v1/Systems/437XR1138R2
    curl --fail http://127.0.0.1:8002/redfish/v1/Systems/437XR1138R2
    curl --fail http://127.0.0.1:8082/var/lib/tftpboot/
    curl --fail http://127.0.0.1:9091/
    for i in $(seq 1 20)
    do
        echo "${i}"
	if [[ "$(curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{"id": 1, "method": "spdk_get_version"}' http://127.0.0.1:9004)" ]]
	then
            break
	else
            sleep 1
        fi
    done
    for i in $(seq 1 20)
    do
        echo "$i"
        if [[ "$(curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{"id": 1, "method": "spdk_get_version"}' http://127.0.0.1:9009)" ]]
        then
            break
        else
            sleep 1
        fi
    done
    curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{"id": 1, "method": "bdev_get_bdevs"}' http://127.0.0.1:9004
    curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{"id": 1, "method": "bdev_get_bdevs"}' http://127.0.0.1:9009
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2210 host@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2208  bmc@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2207  xpu@127.0.0.1 hostname
    sshpass -p 123456 ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2209  bmc@127.0.0.1 hostname
    bash -c "${DC} exec -T pxe nmap --script broadcast-dhcp-discover"
    bash -c "${DC} exec -T pxe cat /var/lib/dhcpd/dhcpd.leases"
    bash -c "${DC} exec -T pxe nmap --script broadcast-dhcp-discover" | grep "Server Identifier: 10.127.127.3"
    bash -c "${DC} exec -T pxe curl --fail http://10.127.127.16:8082/var/lib/tftpboot/"
    bash -c "${DC} exec -T pxe bash -c 'tftp 10.127.127.3 -v -c get grubx64.efi && diff ./grubx64.efi /var/lib/tftpboot/grubx64.efi'"
    bash -c "${DC} exec -T sztp ./run-sztpd-test.sh"
    bash -c "${DC} exec -T spdk-target /usr/local/bin/identify -r 'traddr:10.129.129.4 trtype:TCP adrfam:IPv4 trsvcid:4420'"
    bash -c "${DC} exec -T xpu-spdk /usr/local/bin/identify    -r 'traddr:10.129.129.4 trtype:TCP adrfam:IPv4 trsvcid:4420'"
    bash -c "${DC} exec -T spdk-target /usr/local/bin/perf     -r 'traddr:10.129.129.4 trtype:TCP adrfam:IPv4 trsvcid:4420' -c 0x1 -q 1 -o 4096 -w randread -t 10"
    bash -c "${DC} exec -T xpu-spdk /usr/local/bin/perf         -r 'traddr:10.129.129.4 trtype:TCP adrfam:IPv4 trsvcid:4420' -c 0x1 -q 1 -o 4096 -w randread -t 10"
    NETWORK_CLIENT_NAME=$(${DC} ps | grep example-network-client | awk '{print $1}')
    NETWORK_CLIENT_RC=$(docker inspect --format '{{.State.ExitCode}}' "${NETWORK_CLIENT_NAME}")
    if [ "${NETWORK_CLIENT_RC}" != "0" ]; then
        echo "example-network-client failed:"
        docker logs "${NETWORK_CLIENT_NAME}"
        exit 1
    fi
    STORAGE_CLIENT_NAME=$(${DC} ps | grep example-storage-client | awk '{print $1}')
    STORAGE_CLIENT_RC=$(docker inspect --format '{{.State.ExitCode}}' "${STORAGE_CLIENT_NAME}")
    if [ "${STORAGE_CLIENT_RC}" != "0" ]; then
        echo "example-storage-client failed:"
        docker logs "${STORAGE_CLIENT_NAME}"
        exit 1
    fi
    bash -c "${DC} exec -T strongswan swanctl --stats"
    bash -c "${DC} exec -T strongswan swanctl --list-sas"



    # This should be last
    bash -c "${DC} ps"
    # shellcheck disable=SC2046
    uniq -c <<< "$(sort <<< "$(docker inspect --format='{{json .State.Health.Status}}' $(${DC} ps -q))")"
}

acquire_logs() {
    bash -c "${DC} ps -a"
    bash -c "${DC} logs" || true
    netstat -an || true
    ifconfig -a || true
}

stop_containers() {
    bash -c "${DC} down --volumes"
}

if [ "$#" -lt 1 ]
then
    usage
    exit 1
fi

if [ "$1" == "${INT_BUILD}" ]
then
    build_containers
elif [ "$1" == "${INT_START}" ]
then
    start_containers
elif [ "$1" == "${INT_TESTS}" ]
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

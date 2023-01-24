#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation

set -euxo pipefail

INT_BUILD=build
INT_START=start
INT_TESTS=tests
INT_LOGS=logs
INT_STOP=stop

# docker compose plugin
command -v docker-compose || { shopt -s expand_aliases && alias docker-compose='docker compose'; }

usage() {
    echo ""
    echo "Usage: integration.sh [${INT_BUILD} | ${INT_START} | ${INT_TESTS} | ${INT_LOGS} | ${INT_STOP}]"
    echo ""
}

build_containers() {
    docker-compose build --parallel
}

start_containers() {
    docker-compose down
    docker network prune --force
    docker-compose pull
    docker-compose up -d
}

run_integration_tests() {
    docker-compose ps
    # shellcheck disable=SC2046
    uniq -c <<< "$(sort <<< "$(docker inspect --format='{{json .State.Health.Status}}' $(docker-compose ps -q))")"

    # TODO: replace sleep with timeouted-wait for all services to become healthy
    echo wait 5s... && sleep 5s
    curl --fail http://127.0.0.1:8001/redfish/v1/Systems/437XR1138R2
    curl --fail http://127.0.0.1:8002/redfish/v1/Systems/437XR1138R2
    curl --fail http://127.0.0.1:8082/var/lib/misc/

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

    docker-compose exec -T dhcp cat /var/lib/dhcpd/dhcpd.leases
    docker-compose run nmap
    docker-compose run nmap | grep "Server Identifier: 10.127.127.3"
    # docker-compose exec -w /tmp/sztpd-simulator -T bootstrap ./run-sztpd-test.sh
    docker-compose exec -T redirecter curl --fail -H 'Accept:application/yang-data+json' http://127.0.0.1:1080/.well-known/host-meta
    docker-compose exec -T bootstrap curl --fail -H 'Accept:application/yang-data+json' http://127.0.0.1:1080/.well-known/host-meta

    docker-compose exec -T spdk-target /usr/local/bin/identify -r 'traddr:10.129.129.4 trtype:TCP adrfam:IPv4 trsvcid:4420'
    docker-compose exec -T xpu-spdk /usr/local/bin/identify    -r 'traddr:10.129.129.4 trtype:TCP adrfam:IPv4 trsvcid:4420'
    docker-compose exec -T spdk-target /usr/local/bin/perf     -r 'traddr:10.129.129.4 trtype:TCP adrfam:IPv4 trsvcid:4420' -c 0x1 -q 1 -o 4096 -w randread -t 10
    docker-compose exec -T xpu-spdk /usr/local/bin/perf         -r 'traddr:10.129.129.4 trtype:TCP adrfam:IPv4 trsvcid:4420' -c 0x1 -q 1 -o 4096 -w randread -t 10

    curl --fail http://127.0.0.1:9091/api/v1/query?query=mem_free | grep mem_free
    curl --fail http://127.0.0.1:9091/api/v1/query?query=cpu_usage_user | grep cpu_usage_user
    curl --fail http://127.0.0.1:9091/api/v1/query?query=xpu_num_blocks | grep xpu_num_blocks
    curl --fail http://127.0.0.1:9091/api/v1/query?query=net_bytes_recv | grep net_bytes_recv
    curl --fail http://127.0.0.1:9091/api/v1/query?query=redfish_thermal_fans_reading_rpm | grep redfish_thermal_fans_reading_rpm

    SZTP_AGENT_NAME=$(docker-compose ps | grep agent | awk '{print $1}')
    SZTP_AGENT_RC=$(docker wait "${SZTP_AGENT_NAME}")
    if [ "${SZTP_AGENT_RC}" != "0" ]; then
        echo "${SZTP_AGENT_NAME} failed:"
        docker logs "${SZTP_AGENT_NAME}"
        exit 1
    fi

    NETWORK_CLIENT_NAME=$(docker-compose ps | grep example-network-client | awk '{print $1}')
    NETWORK_CLIENT_RC=$(docker inspect --format '{{.State.ExitCode}}' "${NETWORK_CLIENT_NAME}")
    if [ "${NETWORK_CLIENT_RC}" != "0" ]; then
        echo "example-network-client failed:"
        docker logs "${NETWORK_CLIENT_NAME}"
        exit 1
    fi
    STORAGE_CLIENT_NAME=$(docker-compose ps | grep example-storage-client | awk '{print $1}')
    STORAGE_CLIENT_RC=$(docker inspect --format '{{.State.ExitCode}}' "${STORAGE_CLIENT_NAME}")
    if [ "${STORAGE_CLIENT_RC}" != "0" ]; then
        echo "example-storage-client failed:"
        docker logs "${STORAGE_CLIENT_NAME}"
        exit 1
    fi
    docker-compose exec -T strongswan swanctl --stats
    docker-compose exec -T strongswan swanctl --list-sas



    # This should be last
    docker-compose ps
    # shellcheck disable=SC2046
    uniq -c <<< "$(sort <<< "$(docker inspect --format='{{json .State.Health.Status}}' $(docker-compose ps -q))")"
}

acquire_logs() {
    docker-compose ps -a
    docker-compose logs || true
    netstat -an || true
    ifconfig -a || true
}

stop_containers() {
    docker-compose down --volumes
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

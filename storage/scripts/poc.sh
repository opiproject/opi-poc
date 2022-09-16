#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation
# Copyright (c) 2022 Dell Inc, or its subsidiaries.

set -euxo pipefail

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

# DOCKER_COMPOSE setup
DC=docker-compose

if [ "$(command -v $DC)" == "" ]
then
    DC="docker compose"
fi

usage() {
    echo "Usage: poc.sh [start|stop|tests|logs]"
}

tests_poc() {
    if [ "$(uname -m)" == "x86_64" ]; then
        # Show x86-64-v level.
        "${SCRIPT_DIR}"/x86v.sh
    fi
    $DC ps -a
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
    curl --fail --insecure --user spdkuser:spdkpass -X POST -H 'Content-Type: application/json' -d '{"id": 1, "method": "bdev_get_bdevs"}' http://127.0.0.1:9009
    $DC run opi-spdk-client

    # check reflection
    grpc_cli=(docker run --network=storage_opi --rm namely/grpc-cli)
    "${grpc_cli[@]}" ls opi-spdk-server:50051
    "${grpc_cli[@]}" ls opi-spdk-server:50051 opi.storage.v1.NVMeControllerService -l
    "${grpc_cli[@]}" ls opi-spdk-server:50051 opi.storage.v1.NVMeNamespaceService -l
    "${grpc_cli[@]}" ls opi-spdk-server:50051 opi.storage.v1.NVMeSubsystemService -l
    "${grpc_cli[@]}" ls opi-spdk-server:50051 opi.storage.v1.NVMfRemoteControllerService -l
    "${grpc_cli[@]}" ls opi-spdk-server:50051 opi.storage.v1.VirtioBlkService -l
    "${grpc_cli[@]}" ls opi-spdk-server:50051 opi.storage.v1.VirtioScsiControllerService -l
    "${grpc_cli[@]}" ls opi-spdk-server:50051 opi.storage.v1.VirtioScsiLunService -l

    # this is last line
    $DC ps -a
}

stop_poc() {
    $DC down --volumes
    docker network prune --force
}

start_poc() {
    $DC up --detach
}

logs_poc() {
    $DC ps -a
    $DC logs || true
    docker inspect "$($DC ps -q)" || true
    netstat -an || true
    ifconfig -a || true
}

stop_containers() {
    $DC down --volumes
}

if [ $# -eq 0 ]
  then
	usage
	exit 1
fi

case ${1} in
    start)
        echo "Starting PoC"
        stop_poc
        start_poc
        ;;
    stop)
        echo "Stopping PoC"
        stop_poc
        ;;
    tests)
        echo "Testing PoC"
        tests_poc
        ;;
    logs)
        echo "Logs PoC"
        logs_poc
        ;;
    *)
        usage
        ;;
esac

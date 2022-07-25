#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
# Copyright (c) 2022 Intel Corporation
# Copyright (c) 2022 Dell Inc, or its subsidiaries.

set -euxo pipefail

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

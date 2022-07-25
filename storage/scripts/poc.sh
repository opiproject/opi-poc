#!/bin/bash
#
# Copyright (c) 2022 Intel Corporation
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

usage() {
	echo "Usage: poc.sh [start|stop|tests|logs]"
}

tests_poc() {
	docker-compose ps -a
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
	docker-compose ps -a
}

stop_poc() {
	# Bring the containers down
	docker-compose down
}

start_poc() {
    docker-compose down
    docker network prune --force
    docker-compose up -d
}

logs_poc() {
    docker-compose ps -a
    docker-compose logs || true
    netstat -an || true
    ifconfig -a || true
}

stop_containers() {
    bash -c "${DC} down --volumes"
}

case ${1} in
	start)
		echo "Starting PoC"
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

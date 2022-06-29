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
	echo "Usage: poc.sh [start|stop]"
}

start_poc() {
	# Bringup the containers
	docker compose up -d

	# Verify things work
	ping -c 3 192.168.55.2
	ping -c 3 192.168.65.2
}

stop_poc() {
	# Bring the containers down
	docker compose down

	# Stop ipdk-plugin
	sudo killall ipdk-plugin

	# Stop IPDK container
	ipdk stop
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
	*)
		usage
		;;
esac

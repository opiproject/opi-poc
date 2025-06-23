#!/usr/bin/env bash
set -euo pipefail
#This file checks the status of nginx pod running on the IPU 
#once for DURATION once every INTERVAL

# Command to run
CMD='curl --cacert /root/summit/server.crt https://172.16.3.200/nginx_status'

# Total duration (in seconds) and interval
DURATION=$((60 * 60))   # 60 minutes
INTERVAL=2              # every 2 seconds

START_TIME=$SECONDS

while (( SECONDS - START_TIME < DURATION )); do
  echo "=== $(date +'%Y-%m-%d %H:%M:%S') ==="
  $CMD
  echo ""
  sleep $INTERVAL
done

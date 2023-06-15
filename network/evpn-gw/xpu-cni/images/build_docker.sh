#!/usr/bin/env bash
set -e

## Build docker image
docker build -t xpu-cni -f ../Dockerfile  ../
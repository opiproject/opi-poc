# API GW

This directory contains the Dockerfile and configuration for an API GW.
Currently, we are using [Kong](https://github.com/Kong/kong) as the API GW.
Kong was selected due to it's ease of use, and it's configurability. We
use the [docker-compose](https://github.com/Kong/docker-kong/blob/master/compose/docker-compose.yml)
file from the Kong repository, which is integrated into the xPU file.

As Kong is configured to front the API endpoints on the xPU, we will add
notes here.

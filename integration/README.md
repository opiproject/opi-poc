# Integration

## Diagram

![DPU SW Components](xPU-Integration-Blocks.png)

## Start

```bash
docker-compose -f docker-compose.yml -f docker-compose.otel.yml -f docker-compose.pxe.yml up
```

## Test

* Check Prometheus at <http://0.0.0.0:9090>
* Check Platform/Host BMC redfish server <http://0.0.0.0:8001/redfish/v1>
* Check NIC/DPU/IPU BMC redfish server <http://0.0.0.0:8002/redfish/v1>
* Check PXE server <http://0.0.0.0:8082/var/lib/tftpboot>
* Log into Host CPU: `ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2210 host@127.0.0.1`
* Log into Host BMC: `ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2208 bmc@127.0.0.1`
* Log into  xPU CPU: `ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2207 xpu@127.0.0.1`
* Log into  xPU BMC: `ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2209 bmc@127.0.0.1`

## Stop

```bash
docker-compose -f docker-compose.yml -f docker-compose.otel.yml -f docker-compose.pxe.yml down --remove-orphans
```

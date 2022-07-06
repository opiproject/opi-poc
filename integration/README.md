# Integration

## Diagram

![DPU SW Components](xPU-Integration-Blocks.png)

## Start

```bash
docker-compose -f docker-compose.yml -f docker-compose.telegraf.yml -f docker-compose.pxe.yml up
```

## Test

* Check Prometheus at <http://0.0.0.0:9090>
* Check Platform/Host BMC redfish server <http://0.0.0.0:8001/redfish/v1>
* Check NIC/DPU/IPU BMC redfish server <http://0.0.0.0:8002/redfish/v1>
* Check PXE server <http://0.0.0.0:8082/var/lib/tftpboot>

## Stop

```bash
docker-compose -f docker-compose.yml -f docker-compose.telegraf.yml -f docker-compose.pxe.yml down --remove-orphans
```

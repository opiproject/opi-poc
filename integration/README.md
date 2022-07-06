# Integration

## Diagram

![DPU SW Components](xPU-Integration-Blocks.png)

## Start

```bash
docker-compose -f docker-compose.yml -f docker-compose.telegraf.yml -f docker-compose.pxe.yml up
```

## Stop

```bash
docker-compose -f docker-compose.yml -f docker-compose.telegraf.yml -f docker-compose.pxe.yml down --remove-orphans
```

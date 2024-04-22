# Management server

as of now this is merged with dh1 until new server arives

## Roles

- DHCP server for the lab
- GitHub self hosted runners
- OTEL, Prometheus and Grafana servers for the lab monitoring
- SZTP bootstrap server for the lab
- network gateway between private network 172.22 and the internet
- Other containers and vms for remaining services

## BIOS settings

- change power to always on

## DHCP server

run bare metal or via docker:

```bash
sudo docker run --rm -it --network=host --privileged --restart always -v "$(PWD)"/fs/etc/dhcp/dhcpd.conf:/etc/dhcp/dhcpd.conf docker.io/networkboot/dhcpd:1.3.0
```

or [compose](https://github.com/opiproject/sztp/blob/0addb57154332e7ecdebe4ba18b2633278857ef5/docker-compose.yml#L51-L68)

## Monitoring & Telemetry

Run `docker compose up -d` or `docker-compose up -d` from [here](../../otel)

1. [OTEL Gateway Collector](https://opentelemetry.io/docs/collector/deployment/gateway/) to aggregate telemetry from all DPUs and IPUs.
2. [Prometheus](https://prometheus.io/) Monitoring system & time series database
3. [Grafana](https://grafana.com/) Open source analytics & monitoring solution for every database.

## Paswordless

for all servers, dpu and ipus from [here](../../ips.md) run:

```bash
 ssh-copy-id root@172.22.X.X
```

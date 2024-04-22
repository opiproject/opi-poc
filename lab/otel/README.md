# Monitoring & Telemetry

Took from <https://github.com/opiproject/otel>

![OPI Telemetry Deploy Option](https://github.com/opiproject/otel/blob/main/doc/dpu-otel.png)

## On DPUs and IPUs

### Configuration

Create `telegraf.conf` file, see example [here](./config/telegraf.conf)

- change `outputs.opentelemetry` to the management server name/ip
- change `192.168.240.1` to the internal DPU/IPU AMC/BMC for redfish collection
- make sure [SPDK](https://spdk.io/) app and [spdk_rpc_http_proxy.py](https://github.com/spdk/spdk/blob/v24.01.x/scripts/rpc_http_proxy.py) script are running to collect `storage` statistics

### Service

Run telegraf container:

```bash
sudo docker run -d --restart=always --network=host -v ./config/telegraf.conf:/etc/telegraf/telegraf.conf docker.io/library/telegraf:1.29
```

### Optional Temperature

For regular Servers, add to your config file:

```ini
[[inputs.temp]]
  # no configuration
```

For `Nvidia BlueField` cards, to monitor temperature, add to your config file:

```ini
[[inputs.file]]
  files = ["/run/emu_param/bluefield_temp"]
  name_override = "temp"
  value_field_name="temp"
  data_format = "value"
  data_type = "integer"
  file_tag = "sensor"
```

and add to your docker run command:

```text
-v /run/emu_param:/run/emu_param
```

For `Intel MEV` cards the temperature is on the ICC chip, no easy access to it:

```ini
[[inputs.exec]]
  commands = ["iset-cli get-temperature"]
  name_override = "temp"
  data_format = "json"
```

## On Management server

See management server details [here](../hardware/mgmt)

### Run

Run `docker compose up -d` or `docker-compose up -d`

:exclamation: `docker-compose` is deprecated. For details, see [Migrate to Compose V2](https://docs.docker.com/compose/migrate/).

This will start those services:

1. [OTEL Gateway Collector](https://opentelemetry.io/docs/collector/deployment/gateway/) to aggregate telemetry from all DPUs and IPUs.
2. [Prometheus](https://prometheus.io/) Monitoring system & time series database
3. [Grafana](https://grafana.com/) Open source analytics & monitoring solution for every database.

### Prometheus

1. Open <http://172.22.0.1:9091> to explore Prometheus UI
2. or via API examples:

```bash
curl --fail http://172.22.0.1:9091/api/v1/query?query=mem_free | grep mem_free
curl --fail http://172.22.0.1:9091/api/v1/query?query=cpu_usage_user | grep cpu_usage_user
curl --fail http://172.22.0.1:9091/api/v1/query?query=spdk_num_read_ops | grep spdk_num_read_ops
curl --fail http://172.22.0.1:9091/api/v1/query?query=nstat_TcpActiveOpens | grep nstat_TcpActiveOpens
curl --fail http://172.22.0.1:9091/api/v1/query?query=redfish_power_powercontrol_interval_in_min | grep redfish_power_powercontrol_interval_in_min
```

### Grafana

1. Open <http://172.22.0.1:3000> to explore Grafana UI
2. or via API examples:

```bash
curl -s http://172.22.0.1:3000/api/datasources | jq .
```

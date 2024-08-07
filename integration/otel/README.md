# Open Telemetry

Telegraf (agent) still runs inside DPU/IPU and connects to OTEL gateway. Later on replace it with OTEL agent.

OTEL gateway collector runs outside of DPU/IPU and connects to backends like ELK, Jaeger, Zipkin and Prometheus back-ends...

## Docs

* <https://opentelemetry.io/docs/collector/getting-started/>
* <https://github.com/open-telemetry/opentelemetry-collector-contrib/blob/main/receiver/hostmetricsreceiver>

* Dockers used
  * <https://hub.docker.com/r/otel/opentelemetry-collector>
  * <https://hub.docker.com/r/prom/prometheus>

## Getting started

:exclamation: `docker-compose` is deprecated. For details, see [Migrate to Compose V2](https://docs.docker.com/compose/migrate/).

Run `docker-compose -f docker-compose.otel.yml up`

## Example

* Prometheus at <http://0.0.0.0:9091>

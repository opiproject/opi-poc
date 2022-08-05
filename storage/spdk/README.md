# OPI Storage SPDK container

This directory contains an example SPDK app to demonstrate the OPI Storage APIs integration.

## Docs

* [JSON RPC Proxy](https://spdk.io/doc/jsonrpc_proxy.html)
* [SPDK SMA](https://github.com/spdk/spdk/tree/master/python/spdk/sma)
* [SPDK CSI](https://github.com/spdk/spdk-csi/blob/master/deploy/spdk/Dockerfile)
* [CSI Spec](https://github.com/container-storage-interface/spec/blob/master/spec.md)

## Getting started

Run `docker-compose up spdk`

## SPDK RPC proxy

```text
$ curl -k --user spdkuser:spdkpass -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "bdev_get_bdevs", "params": {"name": "Malloc0"}}' http://127.0.0.1:9009/
{"jsonrpc":"2.0","id":1,"result":[{"name":"Malloc0","aliases":["f1c5d95a-b235-40af-9e4d-2c0b3320de80"],"product_name":"Malloc disk","block_size":512,"num_blocks":131072,"uuid":"f1c5d95a-b235-40af-9e4d-2c0b3320de80","assigned_rate_limits":{"rw_ios_per_sec":0,"rw_mbytes_per_sec":0,"r_mbytes_per_sec":0,"w_mbytes_per_sec":0},"claimed":false,"zoned":false,"supported_io_types":{"read":true,"write":true,"unmap":true,"write_zeroes":true,"flush":true,"reset":true,"nvme_admin":false,"nvme_io":false},"driver_specific":{}}]}
```

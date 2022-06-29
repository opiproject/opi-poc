# Storage

## Docs

* https://spdk.io/doc/jsonrpc_proxy.html
* https://github.com/spdk/spdk/tree/master/python/spdk/sma
* https://github.com/spdk/spdk-csi/blob/master/deploy/spdk/Dockerfile
* https://github.com/container-storage-interface/spec/blob/master/spec.md

## Getting started:

Run `docker-compose up -d`

## SPDK RPC proxy

```
curl -k --user spdkuser:spdkpass -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "bdev_get_bdevs", "params": {"name": "Malloc0"}}' http://127.0.0.1:9009/
```


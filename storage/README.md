# Storage

This is a simple spdk based storage PoC

## Docs

* [JSON RPC Proxy](https://spdk.io/doc/jsonrpc_proxy.html)
* [SPDK SMA](https://github.com/spdk/spdk/tree/master/python/spdk/sma)
* [SPDK CSI](https://github.com/spdk/spdk-csi/blob/master/deploy/spdk/Dockerfile)
* [CSI Spec](https://github.com/container-storage-interface/spec/blob/master/spec.md)

## Getting started

Run `docker-compose up -d`

## Test SPDK is up

```bash
curl -k --user spdkuser:spdkpass -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "bdev_get_bdevs", "params": {"name": "Malloc0"}}' http://127.0.0.1:9009/
```

## SPDK gRPC example

Optionally if you need to download modules

```bash
docker run --rm -it -v `pwd`:/app -w /app golang:alpine go get all
docker run --rm -it -v `pwd`:/app -w /app golang:alpine go get github.com/opiproject/opi-api/storage/proto
```

Run example server (not for production) manually

```bash
   docker run --rm -it -v `pwd`:/app -w /app -p 50051:50051 golang:alpine go run jsonrpc.go server.go
```

Run example client (not for production) manually

```bash
   docker run --net=host --rm -it -v  `pwd`:/app -w /app golang:alpine go run client.go
```

Run both examples client and server via compose (not for production)

```bash
   docker-compose up opi-spdk-client
```

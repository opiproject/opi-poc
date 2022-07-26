# Storage

This is a simple spdk based storage PoC

## Docs

* [JSON RPC Proxy](https://spdk.io/doc/jsonrpc_proxy.html)
* [SPDK SMA](https://github.com/spdk/spdk/tree/master/python/spdk/sma)
* [SPDK CSI](https://github.com/spdk/spdk-csi/blob/master/deploy/spdk/Dockerfile)
* [CSI Spec](https://github.com/container-storage-interface/spec/blob/master/spec.md)

## Getting started

Run `docker-compose up -d`

## SPDK RPC proxy

```bash
curl -k --user spdkuser:spdkpass -X POST -H "Content-Type: application/json" -d '{"id": 1, "method": "bdev_get_bdevs", "params": {"name": "Malloc0"}}' http://127.0.0.1:9009/
```

## SPDK gRPC example

Compile protobufs

```bash
   docker run -v $PWD:/defs namely/protoc-all -d proto -l go -o ./proto/  --go-source-relative
```

Run example server (not for production)

```bash
   docker run --rm -it -v `pwd`:`pwd` -w `pwd` -p 50051:50051 golang:alpine go run server.go
```

Run example client (not for production)

```bash
   docker run --net=host --rm -it -v  `pwd`:`pwd` -w `pwd` golang:alpine go run client.go
```

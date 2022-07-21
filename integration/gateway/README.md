# Gateway

This example of gRPC API gataway based on `nginx`
See <https://www.nginx.com/blog/deploying-nginx-plus-as-an-api-gateway-part-3-publishing-grpc-services/>

And two example services: network (in python) and storage (in go)

And two example clients that send gRPC request via `nginx` gRPC API gataway to network and storage services.

## Run

```text
docker-compose up --build example-storage-client example-network-client
```

## Scale

TBD: how to update ngnix configuration.

```text
docker-compose up --scale example-network=3 --scale example-storage=2 gateway
```

## Security

TBD <https://www.nginx.com/blog/deploying-nginx-plus-as-an-api-gateway-part-1/>

## Authentication

TBD <https://www.nginx.com/blog/deploying-nginx-plus-as-an-api-gateway-part-1/>

## Monitoring

TBD via OTEL

# Storage demo

## Video

tbd

## API

[API protobufs](https://github.com/opiproject/opi-api/tree/main/storage)

## Server

see <https://hub.docker.com/r/opiproject/opi-nvidia-bridge>

```Shell
$ docker run --rm -it -v /var/tmp/:/var/tmp/ -p 50051:50051 ghcr.io/opiproject/opi-nvidia-bridge:main
2022/11/29 00:03:55 plugin serevr is &{{}}
2022/11/29 00:03:55 server listening at [::]:50051
```

## Client

see <https://github.com/opiproject/pydpu> and <https://github.com/opiproject/godpu>

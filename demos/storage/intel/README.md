# Storage demo

## Video

See <https://www.youtube.com/watch?v=9EHZ_1ARCiY>

## API

[API protobufs](https://github.com/opiproject/opi-api/tree/main/storage)

## Server

see <https://hub.docker.com/r/opiproject/opi-intel-bridge>

```Shell
$ docker run --rm -it -v /var/tmp/:/var/tmp/ -p 50051:50051 ghcr.io/opiproject/opi-intel-bridge:main
2022/11/29 00:03:55 plugin serevr is &{{}}
2022/11/29 00:03:55 server listening at [::]:50051
```

## Client

see <https://github.com/opiproject/pydpu> and <https://github.com/opiproject/godpu>

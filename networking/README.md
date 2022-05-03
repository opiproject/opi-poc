# Software Firewall PoC

The Proof of Concept is made up of three types of containers:

- traffic sources
- traffic targets
- proxies

The initial PoC will be:

- IPDK P4-eBPF container for networking
- Docker-based container setup using Docker Compose
- Server containers hosting traffic sources
- Traffic target containers

# Running the PoC

The OPI PoC uses IPDK with p4-ebpf for networking. You will want to install
both [IPDK p4-ebpf](https://github.com/ipdk-io/ipdk/tree/main/build/networking_ebpf)
as well as the [IPDK Docker CNI](https://github.com/mestery/ipdk-plugin).

## Installing Dependencies

You will need to have both Docker engine and Docker Compose installed
to run the demo.

Follow the [README](https://github.com/ipdk-io/ipdk/blob/main/build/networking_ebpf/README_DOCKER.md)
to install the IPDK p4-ebpf container.

Then, follow the [instructions](https://github.com/mestery/ipdk-plugin) to
install the ipdk-plugin to manage docker networks.

Next, start both ipdk and ipdk-plugin:

```
$ ipdk start -d --link-namespace --ulimit
$ sudo ~/go/src//ipdk-plugin &
```

## Running the PoC

You can now run the demo setup as follows:

```
$ docker compose up -d
```

## Verify the PoC

The next step is to verify you have connectivity across the containers via the
networks IPDK p4-ebpf is providing. The following shows an example of pinging
from clients to the nginx container:

```
$ docker exec -it client1 ping -c 3 192.168.55.2
PING 192.168.55.2 (192.168.55.2): 56 data bytes
64 bytes from 192.168.55.2: seq=0 ttl=63 time=0.094 ms
64 bytes from 192.168.55.2: seq=1 ttl=63 time=0.109 ms
64 bytes from 192.168.55.2: seq=2 ttl=63 time=0.132 ms

--- 192.168.55.2 ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 0.094/0.111/0.132 ms
$ docker exec -it client2 ping -c 3 192.168.65.2
PING 192.168.65.2 (192.168.65.2): 56 data bytes
64 bytes from 192.168.65.2: seq=0 ttl=63 time=0.075 ms
64 bytes from 192.168.65.2: seq=1 ttl=63 time=0.095 ms
64 bytes from 192.168.65.2: seq=2 ttl=63 time=0.103 ms

--- 192.168.65.2 ping statistics ---
3 packets transmitted, 3 packets received, 0% packet loss
round-trip min/avg/max = 0.075/0.091/0.103 ms
$
```

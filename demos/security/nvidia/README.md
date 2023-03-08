# IPSec demo

## hardware

- server with Ubuntu 22.04
- Nvidia BlueField2
- Keysight CloudStorm 100G

## configuration

### host (the server holding the DPU)

- install Ubuntu 22.04 server
- install `MLNX_OFED_LINUX-5.8-1.1.2.1-ubuntu22.04-x86_64.iso`
- install `DOCA_1.5.1_BSP_3.9.3_Ubuntu_20.04-4.2211-LTS.signed.bfb` on the BlueField2
- set BlueField2 in SEPARATED_HOST mode to make things easier

```Shell
mlxconfig -d /dev/mst/mt41686_pciconf0 s INTERNAL_CPU_MODEL=0
mlxconfig -d /dev/mst/mt41686_pciconf0.1 s INTERNAL_CPU_MODEL=0
```

### dpu (BlueField2)

- enable ip forwarding `sysctl net.ipv4.ip_forward=1`
- set ip addresses on the phisical links

```Shell
ip addr add 200.0.0.1/24 dev enp3s0f0s0
ip addr add 201.0.0.1/24 dev enp3s0f1s0
```
  
- install Docker (see [Docker manual](https://docs.docker.com/engine/install/ubuntu/) )
- set the secrets (TODO: see if it is posible to do it via opi-api)

```Shell
cat /etc/swanctl/conf.d/k.swanctl.conf
secrets {
        ike-tun1_0_0 {
                secret = "ipsec"
        }
        ike-tun1_0_1 {
                secret = "ipsec"
        }
```

- start OPI server

see <https://hub.docker.com/r/opiproject/opi-strongswan-bridge>

```Shell
docker pull ghcr.io/opiproject/opi-strongswan-bridge:main

sudo sed -i "s/# socket = unix:\/\/\${piddir}\/charon\.vici/socket = unix:\/\/\/var\/run\/charon\.vici/g" /etc/strongswan.d/charon/vici.conf

docker run --rm --network host \
    --mount src=/var/run,target=/var/run,type=bind \
    -d ghcr.io/opiproject/opi-strongswan-bridge:main /opi-vici-bridge -port=50151
```

- to debug or see what is happening open 2 ssh connections to the DPU

```Shell
# in SSH1, instead of -d use -it flag when starting the container
docker run --rm --network host \
    --mount src=/var/run,target=/var/run,type=bind \
    -it ghcr.io/opiproject/opi-strongswan-bridge:main /opi-vici-bridge -port=50151

# in SSH2, it will show how ipsec tunnels are initiated and terminated
swanctl --log
```

### anywhere (on your pc/laptop on the host server.....)

- run the OPI client application done for ipsec demo

see <https://github.com/opiproject/pydpu> and <https://github.com/opiproject/godpu>

```Shell
# add python api lib to path
git clone https://github.com/opiproject/opi-api.git
export PYTHONPATH=${PYTHONPATH}:${HOME}/opi-api/security/v1/gen/python/

pip3 install grpcio

python3

# run the ipsec-config.py from this folder
```

### ixload

- load `opi-ipsec-demo-1.rxf`
- assign the ports
- click start button

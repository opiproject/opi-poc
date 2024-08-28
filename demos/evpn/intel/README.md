# EVPN demo

## hardware

- server with Ubuntu 22.04
- Intel DPU
- Keysight NOVUS 100G

## configuration

### host (the server holding the DPU)

- install Ubuntu 22.04 server or any other OS since host is not used
- Intel setup info: install SDK 1.4+

### IPU - IMC

- bring links up, (not mandatory but got links up faster)

```Shell
iset-cli set-phy-capabilities --port=1 --phy-type=100gbase_cr4
iset-cli set-phy-capabilities --port=0 --phy-type=100gbase_cr4
```

### IPU - ACC

- enable ip forwarding

```Shell
sysctl net.ipv4.ip_forward=1`
iptables -P FORWARD ACCEPT
```

- set ip addresses and Vlan and other required configuration 

```Shell
ip addr add 200.0.0.1/24 dev enp0s1f0d1
ip link add link  enp0s1f0d2 name  enp0s1f0d2.10 type vlan id 10 
ip link set  enp0s1f0d2 up  
ip link set enp0s1f0d2.10 up 

ip link add name lo0 type dummy 
ifconfig lo0 2.2.2.2 netmask 255.255.255.255 up 
 
ip link add br10 type bridge 
ip link set enp0s1f0d1 up 
ip link add vni10 type vxlan local 2.2.2.2 dstport 4789 id 10 
ip link set vni10 master br10 

ip link set vni10 up 
ip link set br10 up 

ip address add dev br10 172.16.0.1/24 
ip link set enp0s1f0d2.10 master br10 
ip link add link  enp0s1f0d2 name  enp0s1f0d2.20 type vlan id 20 

ip link set enp0s1f0d2.20 up 
ip link add br20 type bridge 
ip link add vni20 type vxlan local 2.2.2.2 dstport 4789 id 20 
ip link set vni20 master br20 

ip link set vni20 up 
ip link set br20 up 

ip address add dev br20 172.16.1.1/24 
ip link set enp0s1f0d2.20 master br20 
```

- install Docker (see [Docker manual](https://docs.docker.com/engine/install/ubuntu/) )
- run required containers
- Default FRR Conf
```Shell
frr version 8.5_git
frr defaults datacenter
hostname dut
no ipv6 forwarding
service integrated-vtysh-config
!
password opi
enable password opi
exit
!
```

```Shell
docker run --privileged --rm --network host \
     --mount src=/work/opi-evpn-bridge/conf/frr.conf,target=/etc/frr/frr.conf,type=bind \
     --cap-add=NET_ADMIN \
     --cap-add=SYS_ADMIN \
     --cap-add=SYS_MODULE \
     --cap-add=CAP_NET_RAW \
     --cap-add=SYS_TIME \
     --cap-add=SYS_PTRACE \
     -it quay.io/frrouting/frr:9.1.0 bash  
```

- Setup vtysh inside the FRR container

```Shell
docker exec -it $FRRContainerID bash
touch /etc/frr/vtysh.conf  
sed -i "s/bgpd=no/bgpd=yes/g" /etc/frr/daemons  
sed -i "s/127.0.0.1/0.0.0.0/g" /etc/frr/daemons  
ip link add name lo0 type dummy  
ifconfig lo0 2.2.2.2 netmask 255.255.255.255 up  
/etc/init.d/frr stop  
/usr/lib/frr/watchfrr -d -F traditional zebra bgpd staticd  
sleep infinity
```

- if watchfrr gives out of memory error
  
```Shell
echo 1 > /proc/sys/vm/overcommit_memory
```

- load vtysh config

```Shell
vtysh
config t
router bgp 65100  
bgp router-id 2.2.2.2 
bgp bestpath as-path multipath-relax 
neighbor 200.0.0.2 remote-as 102 
! 

address-family ipv4 unicast 
  network 2.2.2.2/32 
  redistribute connected 
exit-address-family 
! 
address-family l2vpn evpn 
  neighbor 200.0.0.2 activate 
  advertise-all-vni 
  vni 10 
   route-target import 102:10 
  exit-vni 
  advertise-svi-ip 
  advertise ipv4 unicast 
exit-address-family 
exit 
! 
```
### IxNetwork

- load `intel_l2_evpn_2_vni.ixncfg`
- assign the ports
- click start button

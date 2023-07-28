# Demos

## IPSec offload on Marvell DPU with gRPC

This demo shows configuration of Marvell DPU for IPsec offload using gRPC.
The data path is based on HW accelerated VPP and leverages Ligato for off-chip
configuration.

### hardware

- server with Ubuntu 18.04 (host)
- Marvell Octeon CN10K (DPU)
- external server

### configuration

#### host

- Configure Mellanox 100G nic towards DPU

#### external server

- Run etcd
- Run external Agent on external server

### dpu

- Run IPSec offload container
  This will launch vpp and Ligato (vpp-agent)

### Video Recording

see <https://marvell.zoom.us/rec/share/aiA94wryZcAul5HnOpe8xzCH_LVUOItBraNrjUMsopFsnMiWtDdRIClTSHuwLSb6.BL6WoMiOBX2Vb97L>

---------------------------------------------------------------------

## Strongswan integration into Marvell DPU using opi-strongswan-bridge

This demo shows strongswan configuration on Marvell DPU using
opi-strongswan-bridge. opi-strongswan-bridge is a secure server and allows
IPSec off-chip configuration using gRPC based OPI security APIs. It will write
this configuration to Strongswan via vici socket interface. Strongswan in Linux
control plane handles IKE negotiation. IKE packets received in data path by VPP
are transferred to Linux control plane using lcp plugin in VPP. After IKE
negotiation, VPP imports the ip xfrm config from Linux and encrypts traffic.

### remote client

- Run ipsec-config.py from opi-poc on external server

### dpu applications

- Run opi-strongswan-bridge
- Run VPP
- Run Strongswan

### Demo Link

see <https://marvell.zoom.us/rec/share/J9G--oOZyB7WdZo3Xbp5gXO8pegHYOid0uB3Ujm9l19R4FRzZs97kq530yr48lOV.3dhm7jLYQopko5sI>

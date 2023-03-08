# IPSec offload on Marvell DPU with gRPC

This demo shows configuration of Marvell DPU for IPsec offload using gRPC.
The data path is based on HW accelerated VPP and leverages Ligato for off-chip
configuration.

## hardware

- server with Ubuntu 18.04 (host)
- Marvell Octeon CN10K (DPU)
- external server

## configuration

### host

- Configure Mellanox 100G nic towards DPU

### external server

- Run etcd
- Run external Agent on external server

### dpu

- Run IPSec offload container
  This will launch vpp and Ligato (vpp-agent)

## Video Recording

see <https://marvell.zoom.us/rec/share/aiA94wryZcAul5HnOpe8xzCH_LVUOItBraNrjUMsopFsnMiWtDdRIClTSHuwLSb6.BL6WoMiOBX2Vb97L>

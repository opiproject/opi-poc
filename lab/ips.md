# IP Allocations

Proposal is to use 172.22.0.0/16 network for all ip addresing in the lab, and have 1 gateway that will handle internet access and be entry point for access in the lab.

## Device types ip allocation

| IR Range            | Purpose                                                     |
|-----------------    |---------------------------------------------------------    |
| 172.22.0.0/24       | All management devices like, gateway, kvm, ts, pdu, ...     |
| 172.22.1.0/24       | All servers, dpu hosts, test, ...                           |
| 172.22.2.0/24       | All iDRAC, iLO, BMC...                                      |
| 172.22.3.0/24       | All DPUs                                                    |
| 172.22.4.0/24       | All BMC of the DPU                                          |
| 172.22.222.0/24     | dhcp allocated to any unknown device                        |
| 172.22.254.0/24     | VPN Lease Pool -- IPs dynamically assigned to VPN clients   |

## Per Device IP allocation

| IPv4                          | IPv6     | DEVICE           |
|---------------------------    |------    |----------------- |
| public ip / 172.22.0.1        |          | gateway/mgmt     |
| 172.22.0.5                    |          | ToR              |
| 172.22.0.10                   |          | Network PDU      |
| 172.22.0.20                   |          | KVM              |
| 172.22.0.30                   |          | TS               |
| 172.22.0.226                  |          | F5 VPN Floating  |
| 172.22.0.227                  |          | F5 BIG-IP-01     |
| 172.22.0.228                  |          | F5 VPN Floating  |
|                               |          |                  |
| 172.22.1.1                    |          | DPU Host 1       |
| 172.22.2.1                    |          | DH1 iDRAC        |
| 172.22.3.1                    |          | AMD DSC200       |
|                               |          |                  |
| 172.22.1.2                    |          | DPU Host 2       |
| 172.22.2.2                    |          | DH2 iLO          |
| 172.22.3.2                    |          | Nvidia BF2       |
| 172.22.4.2                    |          | Nvidia BF2 bmc   |
|                               |          |                  |
| 172.22.1.3                    |          | DPU Host 3       |
| 172.22.2.3                    |          | DH3 iLO          |
| 172.22.3.3                    |          | Marvell CN106    |
| 172.22.4.3                    |          | Marvell CN106 bmc |
|                               |          |                  |
| 172.22.1.4                    |          | DPU Host 4       |
| 172.22.2.4                    |          | DH4 iDRAC        |
| 172.22.4.4                    |          | Intel DYP imc    |
|                               |          |                  |
| 172.22.1.100                  |          | Tgen1            |
| 172.22.2.100                  |          | Tgen1 bmc        |
|                               |          |                  |
| 172.22.1.250                  |          | Arista 100G      |

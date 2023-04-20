# Lab Bill of Materials

[Spending Guidelines](https://github.com/opiproject/opi/blob/main/Policies/spending_guidelines.md)

## Phase 1 testbed diagram

![xPU Rack phase 1](./images/opi-rack-phase1.svg)

## Phase 1 bill of materials

| Item                       | Vendor     | qty | price      | references |
|----------------------------|------------|-----|------------|------------|
| AMD/Pensando DSC100        | Dell       | 1   | $14000     |            |
| Intel Mount Evans          | Intel      | 1   | donation ? |            |
| Marvel ?                   | ?          | 1   | ?          |            |
| Nvidia Blue Field 2        | HPE        | 1   | $14000     |            |
| Test Server                | Supermicro | 1-4?   | $10,262    |            |
| 100G switch                | Arista     | 1   | $20,000    |            |
| Network PDU                | APC        | 2   | $1,170     |            |
| KVM server                 | Vertiv     | 1   | $8,942     |            |
| KVM dongle                 | Vertiv     | 5   | $143       |            |
| KVM to dongle cable 5m     |            | 10  | $10        |            |
| Terminal server            | Vertiv     | 1   | $5,647     |            |
| Serial cable               |            | 5   | $10        |            |
| 3M QSFP28 DAC cable        | Molex      | 12  | $125       |            |
|                            |            |     |            |            |                                                                                                                |

## Details:

### AMD/Pensando DSC100

* Dell PowerEdge R650 Server 210-AYJZ
* Trusted Platform Module 2.0 V3 461-AAIG
* 2.5" Chassis with up to 8 NVMe Drives, RAID Config, 3 PCIe Slots
* 2x Intel Xeon Gold 6326 2.9G, 16C/32T, 11.2GT/s, 24M Cache, Turbo
* Performance Optimized 370-AAIP
* iDRAC9 Datacenter 15G 528-CRVW
* 16x 16GB RDIMM, 3200MT/s, Dual Rank 370-AEVQ
* 1.92TB Enterprise NVMe Read Intensive AG Drive U.2 Gen with carrier 400-BKGW
* Dell DPU, Pensando 100GbE Dual Port QSFP56, DSC100, PCIe
* [Vendor: Dell](https://www.dell.com)

can be aquaired from both Dell and HPE at similar price in a qualified server.

### Intel Mount Evans

* TODO: get details
* donation ?
* details TBD

### Marvel ?

* TODO: get details
* details TBD

### Nvidia Blue Field 2 MBF2H536C-CECOT

* HPE DL360 Gen10+
* 2x Intel Xeon Gold 6326 2.9G, 16C/32T, 11.2GT/s, 24M Cache, Turbo
* 16x 16GB 2Rx8 PC4-3200AA-R Smart Kit
* 1.6TB NVMe MU SFF BC U.3 PM1735a
* HPE DL36x Gen10+ High Perf Fan Kit
* HPE iLO Adv 1-svr Lic 3yr Support
* HPE DPU Nvidia Blue Field 2 MBF2H536C-CECOT 100GbE 2p SFP56 BMC Spl Adptr
* [Vendor: HPE](https://www.hpe.com)

can be aquaired from both Dell and HPE at similar price in a qualified server.

### Test Server RAX XS4-11S3-10G

* 1U
* Intel C621A Chipset - 4x NVMe/SATA - 1x M.2 - Dual Intel 10-Gigabit Ethernet (RJ45) - 600W Power Supply
* 2x PCIE 4.0 x16
* Intel Xeon Gold 6314U Processor 32-Core 2.3GHz 36MB Cache (185W)
* 8 x 16GB PC4-25600 3200MHz DDR4 ECC RDIMM
* 960GB Micron 7450 PRO Series M.2 PCIe 4.0 x4 NVMe Solid State Drive (80mm)
* Connect X6 2x100G nic
* Thinkmate Server Manager (Datacenter Management Package)
* Thinkmate 3 Year Depot Warranty (Return for Repair)
* [Vendor: Thinkmate](https://www.thinkmate.com/quotation-request?a=YToxOntzOjI6ImlkIjtpOjY0MjIwNDt9)

### 100G switch

* [Arista DCS-7280CR2-60-F](https://www.arista.com/assets/data/pdf/Datasheets/7280R-DataSheet.pdf)
* [Vendor: Hula Networks](https://www.hulanetworks.com)

### Network PDU

* [APC 885-1935 APDU9941](https://www.apc.com/us/en/product/APDU9941/apc-rack-pdu-9000-switched-0u-30a-200v-and-208v-21-c13-and-c15-3-c19-and-c21-sockets/?range=61799-netshelter-switched-rack-pdus&selected-node-id=27602435913) 
* [Vendor: Provantage](https://www.provantage.com/apc-apdu9941~7AMP987M.htm)
* TODO: validate that it can be controled from a python script

### KVM server

* [MergePoint Unity 8032](https://www.vertiv.com/en-us/products-catalog/monitoring-control-and-management/ip-kvm/avocent-mergepoint-unity-digital-kvm-switches)
* [Vendor: Provantage](https://www.provantage.com/vertiv-mpu8032dac-400~7LBRT80Q.htm)
* [KVM dongle](https://www.provantage.com/vertiv-mpuiq-vmchs-g01~7AVOE04X.htm)                                                                                                          

### Terminal Server

* [Avocent ACS8000 48p serial](https://www.vertiv.com/en-us/products-catalog/monitoring-control-and-management/serial-consoles-and-gateways/avocent-acs-8000-serial-consoles)
* [Vendor: Amazon](https://www.amazon.com/Vertiv-Avocent-48-port-Console-ACS8048DAC-400/dp/B01N64R35P?th=1)        

### 3M QSFP28 DAC cable

* [Molex 1002971301](https://www.molex.com/molex/products/part-detail/cable_assemblies/1002971301)
* [Vendor: Digikey](https://www.digikey.com/en/products/detail/molex/1002971301/5361444)
* [Vendor: Mouser](https://www.mouser.com/ProductDetail/Molex/100297-1301?qs=G1LhLIAbs1yJUxZ4qPgMwg%3D%3D)

### Miscelanious cables

* serial cables
* KVM dongle cables
* managemnt network cables
* [Vendor: Monoprice](https://www.monoprice.com)
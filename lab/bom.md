# Lab Bill of Materials

[Spending Guidelines](https://github.com/opiproject/opi/blob/main/Policies/spending_guidelines.md)

No vendor selection has happened yet, bellow table is for illustration/estimate only.

| Item                             | Manufacturer | Part Number        | Quantity | MSRP     | Street  | Budget   | references                                                                                                                                                                                                                                                  |
|----------------------------------|--------------|--------------------|----------|----------|---------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 4 post rack                      |              |                    | 1        |          | $800    | $800     |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| DPUs/IPUs                        |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Intel Mount Evans                | Intel        |                    | 2        | $4,000   | $4,000  | $8,000   |                                                                                                                                                                                                                                                             |
| Intel IPU E2000                  | Intel        |                    | 2        | $4,000   | $4,000  | $8,000   |                                                                                                                                                                                                                                                             |
| Marvell Octeon 10                | Marvell      |                    | 2        | $4,000   | $4,000  | $8,000   |                                                                                                                                                                                                                                                             |
| Nvidia BlueField 2               | Nvidia       |                    | 2        | $4,000   | $4,000  | $8,000   | [Product](https://www.nvidia.com/en-us/networking/products/data-processing-unit/)                                                                                                                                                                                      |
| Nvidia BlueField 3               | Nvidia       |                    | 2        | $4,000   | $4,000  | $8,000   | [Product](https://www.nvidia.com/en-us/networking/products/data-processing-unit/)                                                                                                                                                                                     |
| AMD/Pensando DSC200              | AMD          |                    | 2        | $4,000   | $4,000  | $8,000   | [Product](https://www.amd.com/system/files/documents/pensando-dsc-200-product-brief.pdf)                                                                                                                                                                               |
| need full list of supported _PUs |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Servers                          |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| VM Hypervisor                    | Dell         | RAX QH12-22E4-2GPU | 1        |          | $26,020 | $26,020  | [Shop](https://www.thinkmate.com/quotation-request?a=YToxOntzOjI6ImlkIjtpOjYxOTQ2Njt9)                                                                                                                                                                              |
| DPU Host (1 or 2 _PUs)           | Supermicro   | RAX XS4-11S3-10G   | 10       |          | $5,131  | $51,310  | [Shop](https://www.thinkmate.com)                                                                                                                                                                                                                                   |                                                                                                                                                                                                                                                         |
| Storage Host                     | ASUS         | ESC8000A-E11       | 1        |          | $14,000 | $14,000  | [Product](https://www.asus.com/commercial-servers-workstations/esc8000a-e11/)                                                                                                                                                                                          |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Traffic generators               |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Server with 100G NICs            | Supermicro   | RAX XS4-11S3-10G   | 2        |          | $5,131  | $10,262  | [Shop](https://www.thinkmate.com)                                                                                                                                                                                                                                   |
| hardware traffic generator       |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Switches                         |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| ToR 1G with 10G uplink           | Cisco        |                    | 1        |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| _PU to tgen                      | Arista       | DCS-7280CR2-60-F   | 2        | $149,995 | $20,000 | $40,000  | [Product](https://www.arista.com/assets/data/pdf/Datasheets/7280R-DataSheet.pdf)                                                                                                                                                                                       |
| 400G ?                           |              |                    |          |          |         |          |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| PDUs                             |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| APC 885-1935                     | APC          | APDU9941           | 2        | $1,960   | $1,170  | $2,340   | [Product](https://www.apc.com/us/en/product/APDU9941/apc-rack-pdu-9000-switched-0u-30a-200v-and-208v-21-c13-and-c15-3-c19-and-c21-sockets/?range=61799-netshelter-switched-rack-pdus&selected-node-id=27602435913) [Shop](https://www.provantage.com/apc-apdu9941~7AMP987M.htm) |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| KVMs                             |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| MergePoint Unity 8032            | Vertiv       | MPU8032DAC-400     | 1        | $14,601  | $8,942  | $8,942   | [Product](https://www.vertiv.com/en-us/products-catalog/monitoring-control-and-management/ip-kvm/avocent-mergepoint-unity-digital-kvm-switches/) [Shop](https://www.provantage.com/vertiv-mpu8032dac-400~7LBRT80Q.htm)                                                          |
| KVM dongle                       | Vertiv       | MPUIQ-VMCHS        | 6        | $220     | $143    | $2,002   | [Product](https://www.provantage.com/vertiv-mpuiq-vmchs-g01~7AVOE04X.htm)                                                                                                                                                                                              |
| KVM to dongle cable 5m           |              |                    | 6        | $10      | $10     | $140     |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Terminal servers                 |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Avocent ACS8000 48p serial       | Vertiv       | ACS8048DAC-400     | 1        | $7,720   | $5,647  | $5,647   | [Product](https://www.vertiv.com/en-us/products-catalog/monitoring-control-and-management/serial-consoles-and-gateways/avocent-acs-8000-serial-consoles/) [Shop](https://www.amazon.com/Vertiv-Avocent-48-port-Console-ACS8048DAC-400/dp/B01N64R35P?th=1)                       |
| RJ45 5m serial cables            |              |                    | 14       | $10      | $10     | $140     |                                                                                                                                                                                                                                                             |
| DB9 to RJ45 addapter             |              |                    | 5        | $10      | $10     | $50      |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          |         | $0       |                                                                                                                                                                                                                                                             |
| Cables                           |              |                    |          |          |         | $0       |                                                                                                                                                                                                                                                             |
| 3M QSFP28 DAC cable              | Molex        | 1002971301         | 40       | $135     | $125    | $5,000   |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          |         | $0       |                                                                                                                                                                                                                                                             |
| Licenses                         |              |                    |          |          |         | 0        |                                                                                                                                                                                                                                                             |
| ESXi 8.0 license for 128 cores   | VmWare       | VS8-STD-C          | 4        |          | $1,100  | $4,400   |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          |         |          |                                                                                                                                                                                                                                                             |
| TOTAL                            |              |                    |          |          |         | $237,327 |                                                                                                                                                                                                                                                             |                                                                                                                                                                                                                                                         |

another table probably needed to include the actual deployment info:

* rack space for each device
* power used by each device
* power to heat to HVAC math

## Phase 1 testbed diagram

![xPU Rack phase 1](./images/opi-rack-phase1.svg)

## Phase 1 bill of materials

| Item                             | Manufacturer | Part Number        | Quantity | MSRP     | Street  | Budget   | references                                                                                                                                                                                                                                                  |
|----------------------------------|--------------|--------------------|----------|----------|---------|----------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 4 post rack                      |              |                    | 1        |          | $800    | $800     |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| DPUs/IPUs                        |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Intel Mount Evans                | Intel        |                    | 1        | $4,000   | $4,000  | $4,000   |                                                                                                                                                                                                                                                             |
| Marvell Octeon 10                | Marvell      |                    | 1        | $4,000   | $4,000  | $4,000   |                                                                                                                                                                                                                                                             |
| Nvidia BlueField 2               | Nvidia       |                    | 1        | $4,000   | $4,000  | $4,000   | [Product](https://www.nvidia.com/en-us/networking/products/data-processing-unit/)                                                                                                                                                                                      |
| AMD/Pensando DSC200              | AMD          |                    | 1        | $4,000   | $4,000  | $4,000   | [Product](https://www.amd.com/system/files/documents/pensando-dsc-200-product-brief.pdf)                                                                                                                                                                               |
| need full list of supported _PUs |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Servers                          |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| VM Hypervisor                    | Dell         | RAX QH12-22E4-2GPU | 1        |          | $26,020 | $26,020  | [Shop](https://www.thinkmate.com/quotation-request?a=YToxOntzOjI6ImlkIjtpOjYxOTQ2Njt9)                                                                                                                                                                              |
| DPU Host (1 or 2 _PUs)           | Supermicro   | RAX XS4-11S3-10G   | 3        |          | $5,131  | $15,393  | [Shop](https://www.thinkmate.com)                                                                                                                                                                                                                                   |                                                                                                                                                                                                                                                         |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Traffic generators               |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Server with 100G NICs            | Supermicro   | RAX XS4-11S3-10G   | 1        |          | $10,262  | $10,262  | [Shop](https://www.thinkmate.com)                                                                                                                                                                                                                                   |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Switches                         |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| ToR 1G with 10G uplink           | Cisco        |                    | 1        |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| _PU to tgen                      | Arista       | DCS-7280CR2-60-F   | 1        | $149,995 | $20,000 | $20,000  | [Product](https://www.arista.com/assets/data/pdf/Datasheets/7280R-DataSheet.pdf)                                                                                                                                                                                       |
| 400G ?                           |              |                    |          |          |         |          |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| PDUs                             |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| APC 885-1935                     | APC          | APDU9941           | 2        | $1,960   | $1,170  | $2,340   | [Product](https://www.apc.com/us/en/product/APDU9941/apc-rack-pdu-9000-switched-0u-30a-200v-and-208v-21-c13-and-c15-3-c19-and-c21-sockets/?range=61799-netshelter-switched-rack-pdus&selected-node-id=27602435913) [Shop](https://www.provantage.com/apc-apdu9941~7AMP987M.htm) |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| KVMs                             |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| MergePoint Unity 8032            | Vertiv       | MPU8032DAC-400     | 1        | $14,601  | $8,942  | $8,942   | [Product](https://www.vertiv.com/en-us/products-catalog/monitoring-control-and-management/ip-kvm/avocent-mergepoint-unity-digital-kvm-switches/) [Shop](https://www.provantage.com/vertiv-mpu8032dac-400~7LBRT80Q.htm)                                                          |
| KVM dongle                       | Vertiv       | MPUIQ-VMCHS        | 6        | $220     | $143    | $858     | [Product](https://www.provantage.com/vertiv-mpuiq-vmchs-g01~7AVOE04X.htm)                                                                                                                                                                                              |
| KVM to dongle cable 5m           |              |                    | 6        | $10      | $10     | $60      |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Terminal servers                 |              |                    |          |          | $0      | $0       |                                                                                                                                                                                                                                                             |
| Avocent ACS8000 48p serial       | Vertiv       | ACS8048DAC-400     | 1        | $7,720   | $5,647  | $5,647   | [Product](https://www.vertiv.com/en-us/products-catalog/monitoring-control-and-management/serial-consoles-and-gateways/avocent-acs-8000-serial-consoles/) [Shop](https://www.amazon.com/Vertiv-Avocent-48-port-Console-ACS8048DAC-400/dp/B01N64R35P?th=1)                       |
| RJ45 5m serial cables            |              |                    | 7        | $10      | $10     | $70      |                                                                                                                                                                                                                                                             |
| DB9 to RJ45 addapter             |              |                    | 5        | $10      | $10     | $50      |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          |         | $0       |                                                                                                                                                                                                                                                             |
| Cables                           |              |                    |          |          |         | $0       |                                                                                                                                                                                                                                                             |
| 3M QSFP28 DAC cable              | Molex        | 1002971301         | 20       | $135     | $125    | $2,500   |                                                                                                                                                                                                                                                             |
|                                  |              |                    |          |          |         | $0       |                                                                                                                                                                                                                                                             |
| Licenses                         |              |                    |          |          |         | 0        |                                                                                                                                                                                                                                                             |
| ESXi 8.0 license for 128 cores   | VmWare       | VS8-STD-C          | 4        |          | $1,100  | $4,400   |       ?                                                                                                                                                                                                                                                      |
|                                  |              |                    |          |          |         |          |                                                                                                                                                                                                                                                             |
| TOTAL                            |              |                    |          |          |         | $113,342 |                                                                                                                                                                                                                                                             |                                                                                                                                                                                                                                                         |

## Selection criteria

| Item                            | how                                                                                                                   |
|---------------------------------|-----------------------------------------------------------------------------------------------------------------------|
| features                        | must do the job required, expected, have the feature set                                                              |
| waranty                         | 1 year standard waranty as a minimum, can be extended one more year in the next year budged                           |
| support                         | no or standard support for 1 year, can be extended one more year in the next year budged                              |
| end of life date                | product must not be on an end of life, end of support anouncement                                                     |
| api availability                | functionality must be accesible via API for devices part of CI pipeline, not mandatory for manual testbed but desired |
| lab management tool integration | TBD                                                                                                                   |
| availability                    | must be available for purchase with max 30 days lead time                                                             |
| provenience                     | must be a genuine device                                                                                              |
| vendor / mfg                    | both vendor and product manufacturer must have a good past history                                                    |
| price                           | best price that meets the above criterias                                                                             |

## Server details

* cost to host each _PU
* many in a server or each in a unique server
* mandatory features to test interaction with the host

### DPU HOST - RAX XS4-11S3-10G

* 1U
* Intel C621A Chipset - 4x NVMe/SATA - 1x M.2 - Dual Intel 10-Gigabit Ethernet (RJ45) - 600W Power Supply
* 2x PCIE 4.0 x16
* Intel Xeon Gold 6312U Processor 24-Core 2.4GHz 36MB Cache (185W)
* 8 x 16GB PC4-25600 3200MHz DDR4 ECC RDIMM
* 960GB Micron 7450 PRO Series M.2 PCIe 4.0 x4 NVMe Solid State Drive (80mm)
* Connect X6 2x100G nic
* Thinkmate Server Manager (Datacenter Management Package)
* Thinkmate 3 Year Depot Warranty (Return for Repair)
* ??? does it has PCIE 4/6/8 pin power adapter since some _PU needs it

### DPU HOST - RAX QS4-11E4

* 1U
* has PCIE 5.0 x16

### DPU HOST - ESC8000A-E11

* 4U
* has 8 slots of PCIE 4.0 x16 with dedicated PCIE power connector for each slot
* a reboot will restart all 8 _PUs

### GPU servers may be more apropiate for _PU hosting

* better cooling
* more room (2 PCIE slots per PCIE connector)
* dedicated PCIE power connector

### VM HOST - RAX QH12-22E4-2GPU

* AMD EPYC 9004 Series - 2U - 2 GPUs - 8x 3.5" SAS/SATA + 4x 3.5" Hybrid - 2x M.2 NVMe - Dual GbE (RJ45) - 2400W Redundant
* 2 x AMD EPYC 9554 Processor 64-core 3.10GHz 256MB Cache (360W)
* 24 x 16GB PC5-38400 4800MHz DDR5 ECC RDIMM
* 2 x 3.84TB Samsung PM9A3 Series M.2 PCIe 4.0 x4 NVMe Solid State Drive
* Trusted Platform Module - TPM 2.0
* Thinkmate Server Management Software Package - IPMI and Redfish Compatible
* 2 x AC Power Cord (North America), C13, NEMA 5-15P, 2.1m CAB-AC
* Thinkmate Standard Rail Kit for 1U/2U Servers (Square Hole) (Included)
* Thinkmate 3 Year Advanced Parts Replacement Warranty (Zone 0)

### Pedestal servers , why ???

### storage server or back up in AWS/Azure/etc

### Cables

* Optical cables
* Transceivers
* DAC cables

### Switches

* Top of Rack Layer 2/3 switches
* Layer 1 switches
* Programmable switches

### PDUs

* PDUs
* C13/C14/C15 power cables

### KVMs

* KVMs
* Dongles

### Terminal servers

### DPUs/IPUs

* Generally available DPUs
* Generally available IPUs

### Software

* OS licenses
* Software required to enable hardware

### Hosting services (University of New Hamshire Interoperatability Lab)

* Rack space
* Electricity
* HVAC
* Services (lab tech to change cables)

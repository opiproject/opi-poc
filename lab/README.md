# OPI Lab

- [Goals and Requirments](goals-and-requirements.md)
- Setup
  - [LAB IPs allocation](./ips.md)
  - [Server Software Setup](server-setup.md)
  - [Bill of Materials](bom.md)
  - [Physical Testbed Setup](physical-testbed.md)
    - [ToR - Arista DCS-7280TR-48C6-R](./hardware/ToR/README.md)
    - [TS - Serial Consoles - Avocent ACS8000](./hardware/TS/README.md)
    - [KVM - KVM Switch - Avocent MPU8032DAC](./hardware/KVM/README.md)
    - [VPN - F5 BIG IP i4000](./hardware/VPN/README.md)
    - [nPDU - Network PDU - Vertiv VP7N3001](./hardware/nPDU/README.md)
    - [DH1 - DPU Host 1 - Dell R650](./hardware/dh1/README.md)
    - [DH2 - DPU Host 2 - HPE DL360](./hardware/dh2/README.md)
    - [DH3 - DPU Host 3 - HP RL300](./hardware/dh3/README.md)
    - [DH4 - DPU Host 4 - Dell R760](./hardware/dh4/README.md)
    - [A100G - 100G Switch - Arista DCS-7280CR](./hardware/A100G/README.md)
    - [TGEN1 - Traffic Generator 1 - Supermicro](./hardware/tgen1/README.md)
    - [Cables/Fibers/transceivers](./hardware/cables.md)
- [Running the Test Cases](running-the-tests.md)

## Passwords

We manage all passwords in <https://opi-team.1password.com>

## Access the LAB

- ask in slack for initial user/password
- log into [vpn.opiproject-lab.org](http://vpn.opiproject-lab.org) via web browser and change first time password
- either allow the browser based Network Access application to be installed and the VPN will be launched
- or if your company restricts browser-based installs you can:
  - If on Windows, install [F5 Access](./images/f5-vpn-msft-app.png) app from [www.microsoft.com/store/productId/9WZDNCRDSFN0](https://www.microsoft.com/store/productId/9WZDNCRDSFN0)
  - If on Mac, install the F5 Access app from the Mac App store from [apps.apple.com/us/app/f5access/id1243219105](https://apps.apple.com/us/app/f5access/id1243219105)
- create new OPI [VPN connection](./images/add-vpn-windows.png) and enter your new credentials
- test it by going to [https://172.22.4.2/redfish/v1/](https://172.22.4.2/redfish/v1/)

## Rack diagram

![lab cabling diagram](./images/opi-lab-cabling.drawio.svg)

# 100G Switch

This device is intended to mold the test traffic between the trafic generators and the devices under test

## HW spec

Arista 7280R, 48x100GbE QSFP and 8x40GbE QSFP+ switch, front to rear air, 2 x AC and 2 x C19-C20 cords

- Part number: Arista DCS-7280CR-48-F
- Hardware version: 11.01
- Software image version: 4.20.10M

```bash
A100G#show version
Arista DCS-7280CR-48-F
Hardware version:    11.01
Serial number:       JPE16491339
System MAC address:  444c.a897.b5cf
 
Software image version: 4.20.10M
Architecture:           i386
Internal build version: 4.20.10M-10040268.42010M
Internal build ID:      5ba83857-5952-4713-b850-6e3c7c79cac3
 
Uptime:                 4 weeks, 6 days, 19 hours and 1 minutes
Total memory:           15992160 kB
Free memory:            13583756 kB
```

## Management

- serial port
- 1G ethernet management port

## gNMI

```bash
enable
configure
management api gnmi
   transport grpc openmgmt
      port 5900
```

and

```bash
A100G#show management api gnmi
Enabled:            Yes
Server:             running on port 5900, in default VRF
SSL Profile:        none
```

## Docs

[7280R-DataSheet](https://www.arista.com/assets/data/pdf/Datasheets/7280R-DataSheet.pdf)

## Config

[Config](arista.config)

## Pictures

![front](front.jpg)

![back](back.png)

# PXE

This a container with DHCP, TFTP and HTTP servers

## Mount points

- /var/lib/tftpboot/ serves (initrd, kernel, grub.cfg)
- /etc/dhcp/dhcpd.conf config file for DHCP

## Environment variables

- NODE_IP_SUBNET - The DHCP subnet
- NODE_IP_NETMASK - The DHCP netmask
- NODE_IP_RANGE_MIN - The DHCP Min address
- NODE_IP_RANGE_MAX - The DHCP Max address
- NODE_IP_ADDRESS - The current node IP address

## Run

```text
docker-compose -f docker-compose.pxe.yml up --build
```

## Attach to any container on the same network

```text
docker-compose -f docker-compose.pxe.yml exec pxe bash
```

## Run DHCP discover and get PXE server IP

```text
[root@805bc8fcb44f /]# nmap --script broadcast-dhcp-discover
Starting Nmap 7.80 ( https://nmap.org ) at 2022-07-06 17:55 UTC
Pre-scan script results:
| broadcast-dhcp-discover:
|   Response 1 of 1:
|     IP Offered: 10.127.127.10
|     DHCP Message Type: DHCPOFFER
|     Server Identifier: 10.127.127.103
|     IP Address Lease Time: 5m00s
|_    Subnet Mask: 255.255.255.0
WARNING: No targets were specified, so 0 hosts scanned.
Nmap done: 0 IP addresses (0 hosts up) scanned in 1.20 seconds
```

## Test HTTP web server

```text
[root@805bc8fcb44f /]# curl --fail http://10.127.127.103:8082/var/lib/tftpboot/
<!DOCTYPE HTML PUBLIC "-//W3C//DTD HTML 4.01//EN" "http://www.w3.org/TR/html4/strict.dtd">
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<title>Directory listing for /var/lib/tftpboot/</title>
</head>
<body>
<h1>Directory listing for /var/lib/tftpboot/</h1>
<hr>
<ul>
<li><a href="grubx64.efi">grubx64.efi</a></li>
</ul>
<hr>
</body>
</html>
```

## Test TFTP web server

```text
[root@805bc8fcb44f /]# tftp 10.127.127.103 -v -c get grubx64.efi
Connected to 10.127.127.103 (10.127.127.103), port 69
getting from 10.127.127.103:grubx64.efi to grubx64.efi [netascii]
Received 1028622 bytes in 0.1 seconds [154505490 bit/s]
```

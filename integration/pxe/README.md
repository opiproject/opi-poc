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

## Run DHCP discover and get PXE server IP

```text
docker run --rm -it --network integration_xpu-cpu instrumentisto/nmap:7.92 --script broadcast-dhcp-discover
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

## DHCP discover with custom options

### DHCP Server config

Add new custom option:

- see ipv4 <https://www.iana.org/assignments/bootp-dhcp-parameters/bootp-dhcp-parameters.xml>
- see ipv6 <https://www.iana.org/assignments/dhcpv6-parameters/dhcpv6-parameters.xhtml>

```bash
[root@ae82d8778616 /]# grep sztp /etc/dhcp/dhcpd.conf
option sztp-redirect-urls code 143  = text;
    option sztp-redirect-urls "http://192.0.2.1/demo.sh";
```

### DHCP Client config

```bash
[root@ae82d8778616 /]# cat <<- EOF > /etc/dhcp/dhclient.conf
option sztp-redirect-urls code 143  = text;
request subnet-mask,
          broadcast-address,
          routers,
          domain-name,
          domain-name-servers,
          bootfile-name,
          tftp-server-name,
          sztp-redirect-urls,
          host-name;
EOF
```

### DHCP Client run

```bash
[root@ae82d8778616 /]# dhclient -d -v
Internet Systems Consortium DHCP Client 4.4.3
Copyright 2004-2022 Internet Systems Consortium.
All rights reserved.
For info, please visit https://www.isc.org/software/dhcp/

RTNETLINK answers: Operation not permitted
Listening on LPF/eth0/02:42:0a:7f:7f:03
Sending on   LPF/eth0/02:42:0a:7f:7f:03
Sending on   Socket/fallback
DHCPDISCOVER on eth0 to 255.255.255.255 port 67 interval 8 (xid=0x68047a4a)
DHCPOFFER of 10.127.127.100 from 10.127.127.3
DHCPREQUEST for 10.127.127.100 on eth0 to 255.255.255.255 port 67 (xid=0x68047a4a)
DHCPACK of 10.127.127.100 from 10.127.127.3 (xid=0x68047a4a)
RTNETLINK answers: Operation not permitted
bound to 10.127.127.100 -- renewal in 298 seconds.
^C
```

### DHCP Client results

```bash
[root@ae82d8778616 /]# cat /var/lib/dhclient/dhclient.leases
lease {
  interface "eth0";
  fixed-address 10.127.127.100;
  filename "grubx64.efi";
  option subnet-mask 255.255.255.0;
  option sztp-redirect-urls "http://192.0.2.1/demo.sh";
  option dhcp-lease-time 600;
  option tftp-server-name "w.x.y.z";
  option bootfile-name "test.cfg";
  option dhcp-message-type 5;
  option dhcp-server-identifier 10.127.127.3;
  renew 1 2022/08/15 17:46:29;
  rebind 1 2022/08/15 17:50:14;
  expire 1 2022/08/15 17:51:29;
}
```

## Test HTTP web server

```text
docker-compose run web curl --fail http://10.127.127.103:8082/var/lib/tftpboot/
docker run --rm -it --network integration_xpu-cpu alpine/curl:3.14 --fail http://10.127.127.103:8082/var/lib/tftpboot/
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
docker-compose run tftp tftp 10.127.127.103 -v -c get grubx64.efi
Connected to 10.127.127.103 (10.127.127.103), port 69
getting from 10.127.127.103:grubx64.efi to grubx64.efi [netascii]
Received 1028622 bytes in 0.1 seconds [154505490 bit/s]
```

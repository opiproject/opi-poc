= PXE

This a DHCP, TFTP and HTTP servers

Mount points:

- /var/lib/tftpboot/ serves (initrd, kernel, grub.cfg)
- /etc/dhcp/dhcpd.conf config file for DHCP

Environment variables:

- NODE_IP_SUBNET - The DHCP subnet
- NODE_IP_NETMASK - The DHCP netmask
- NODE_IP_RANGE_MIN - The DHCP Min address
- NODE_IP_RANGE_MAX - The DHCP Max address
- NODE_IP_ADDRESS - The current node IP address

Run Example:

```text
$ docker-compose -f docker-compose.pxe.yml up --build
```

Test:

```text
$ sudo nmap --script broadcast-dhcp-discover
```

# Secure Zero Touch Provisioning (sZTP)

Took from <https://github.com/opiproject/sztp>

## Run on Management server

Start Bootstrap and Web servers from [compose](./docker-compose.yml):

```bash
docker compose up -d
```

Add SZTP options to your DHCP server [config](./hardware/mgmt/fs/etc/dhcp/dhcpd.conf), [for example](https://github.com/opiproject/sztp/blob/main/dhcp/dhcpd.conf.template):

```bash
$ grep sztp /etc/dhcp/dhcpd.conf
option sztp-redirect-urls code 143  = text;
option sztp-redirect-urls "https://bootstrap:8080/restconf/operations/ietf-sztp-bootstrap-server:get-bootstrapping-data";
```

Extract certificates from Bootstrap server:

```bash
docker compose cp bootstrap:/opi.pem /tmp/opi.pem
docker compose cp bootstrap:/tmp/sztpd-simulator/pki/client/end-entity/my_cert.pem /tmp/opi_cert.pem
docker compose cp bootstrap:/tmp/sztpd-simulator/pki/client/end-entity/private_key.pem /tmp/opi_private_key.pem
```

Copy extracted certificates to DPUs:

```bash
scp /tmp/opi*.pem root@172.22.3.2:/mnt/
```

## Run on DPUs

Add sztp option to the dhcp client, [example](https://github.com/opiproject/sztp/blob/main/dhcp/dhclient.conf):

```bash
root@bf2:~# grep sztp /etc/dhcp/dhclient.conf
option sztp-redirect-urls code 143  = text;
request subnet-mask, broadcast-address, time-offset, routers, sztp-redirect-urls,
```

Make sure lease file received the correct option:

```bash
root@bf2:~# DHCLIENT_LEASE_FILE=/var/lib/NetworkManager/dhclient-aa93b667-6aac-3804-91e9-4958e07fdb2f-oob_net0.lease
root@bf2:~# grep sztp ${DHCLIENT_LEASE_FILE}
  option sztp-redirect-urls "https://bootstrap:8080/restconf/operations/ietf-sztp-bootstrap-server:get-bootstrapping-data";
```

Add static hostname resolution per certificate dns limitation:

```bash
# add hostnames per certificate dns generation
root@bf2:~# grep bootstrap /etc/hosts
172.22.0.1      bootstrap web
```

Finally run SZTP agent/client:

```bash
docker run --rm -it --network=host -v /mnt/:/mnt \
  --mount type=bind,source=${DHCLIENT_LEASE_FILE},target=/var/lib/dhclient/dhclient.leases \
  ghcr.io/opiproject/opi-sztp-client:main \
  /opi-sztp-agent daemon --bootstrap-trust-anchor-cert /mnt/opi.pem --device-end-entity-cert /mnt/opi_cert.pem --device-private-key /mnt/opi_private_key.pem
```

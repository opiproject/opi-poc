# Host

Host is the server where xPU is pugged into.

Real user Apps/Vms/Containers are runnig on the host.

In our case we runing iperf for networking and fio for storage, for exmaple.

## Host Main CPU

Runs nvme-cli and fio and some DB, example:

```text
$ nvme list
Node                  SN                   Model                                    Namespace Usage                      Format           FW Rev
--------------------- -------------------- ---------------------------------------- --------- -------------------------- ---------------- --------
/dev/nvme0n1          71T0A0CETC88         Dell Ent NVMe CM6 MU 1.6TB               1          14.79  GB /   1.60  TB    512   B +  0 B   2.1.3
/dev/nvme1n1          71T0A0BGTC88         Dell Ent NVMe CM6 MU 1.6TB               1         888.31  MB /   1.60  TB    512   B +  0 B   2.1.3
/dev/nvme2n1          71T0A0C3TC88         Dell Ent NVMe CM6 MU 1.6TB               1         524.29  kB /   1.60  TB    512   B +  0 B   2.1.3
/dev/nvme3n1          71T0A0C9TC88         Dell Ent NVMe CM6 MU 1.6TB               1         524.29  kB /   1.60  TB    512   B +  0 B   2.1.3
```

fio example:

```text
sudo fio --direct=1 --prio=0 --norandommap=1 --group_reporting --cpus_allowed_policy=split --ioengine=libaio --rw=randrw --rwmixread=100 --bs=4096 --runtime=100 --numjobs=16 --iodepth=4 --name=filename1 --filename=/dev/nvme0n1 --name=filename2 --filename=/dev/nvme0n2 --name=filename3 --filename=/dev/nvme0n3
```

ssh example:

```text
$ ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2210 host@127.0.0.1
Warning: Permanently added '[127.0.0.1]:2210' (ECDSA) to the list of known hosts.
host@127.0.0.1's password:
Welcome to OpenSSH Server

host-cpu:~$ exit
logout
Connection to 127.0.0.1 closed.
```

## Host/Platform BMC

Runs redfish server, example:

```text
$ curl --fail http://127.0.0.1:8001/redfish/v1                                                                                                             {
    "@odata.id": "/redfish/v1/",
    "@odata.type": "#ServiceRoot.v1_6_0.ServiceRoot",
    "AccountService": {
        "@odata.id": "/redfish/v1/AccountService"
    },
    "CertificateService": {
        "@odata.id": "/redfish/v1/CertificateService"
    },
    "Chassis": {
        "@odata.id": "/redfish/v1/Chassis"
    },
    "EventService": {
        "@odata.id": "/redfish/v1/EventService"
    },
    "Id": "RootService",
    "Links": {
        "Sessions": {
            "@odata.id": "/redfish/v1/SessionService/Sessions"
        }
    },
    "Managers": {
        "@odata.id": "/redfish/v1/Managers"
    },
    "Name": "Root Service",
    "Oem": {},
    "RedfishVersion": "1.6.0",
    "Registries": {
        "@odata.id": "/redfish/v1/Registries"
    },
    "SessionService": {
        "@odata.id": "/redfish/v1/SessionService"
    },
    "Systems": {
        "@odata.id": "/redfish/v1/Systems"
    },
    "Tasks": {
        "@odata.id": "/redfish/v1/TaskService"
    },
    "UUID": "92384634-2938-2342-8820-489239905423",
    "UpdateService": {
        "@odata.id": "/redfish/v1/UpdateService"
    }
}
```

ssh example:

```text
$ ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -p 2208 bmc@127.0.0.1
Warning: Permanently added '[127.0.0.1]:2208' (ECDSA) to the list of known hosts.
bmc@127.0.0.1's password:
Welcome to OpenSSH Server

host-bmc:~$ exit
logout
Connection to 127.0.0.1 closed.
```

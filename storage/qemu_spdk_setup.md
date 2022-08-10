
# OPI QEMU SPDK setup

## Docs

From <https://spdk.io/doc/vhost.html>

## Virtualization support

Make sure that VT-x/AMD-v support is enabled in BIOS

```bash
$ lscpu | grep -i virtualization
Virtualization:                  VT-x
```

and that kvm modules are loaded

```bash
$ lsmod | grep -i kvm
kvm_intel             217088  0
kvm                   614400  1 kvm_intel
irqbypass              16384  1 kvm
```

## Tool installation

### wget

Installation on Fedora

```bash
sudo dnf install wget
```

or on Ubuntu

```bash
sudo apt install wget
```

### qemu

From <https://www.qemu.org/download/>

```bash
sudo dnf install qemu-kvm
```

or

```bash
sudo apt install qemu-system
```

### libguestfs-tools

```bash
sudo dnf install libguestfs-tools-c
```

or

```bash
sudo apt install libguestfs-tools
```

## Huge Pages

```bash
echo 4096 | sudo tee /proc/sys/vm/nr_hugepages
```

## Run SPDK

```bash
./spdk/build/bin/spdk_tgt  -S /var/tmp -s 1024 -m 0x3
```

## Download guest image

```bash
wget -O guest_os_image.qcow2  https://download.fedoraproject.org/pub/fedora/linux/releases/36/Cloud/x86_64/images/Fedora-Cloud-Base-36-1.5.x86_64.qcow2
```

## Change password

```bash
cat <<- EOF > meta-data
instance-id: iid-local01;
local-hostname: fed21;
EOF

cat <<- EOF > user-data
#cloud-config
password: fedora
chpasswd: { expire: False }
ssh_pwauth: True
EOF

genisoimage -output init.iso -volid cidata -joliet -rock user-data meta-data
```

## virtio-blk

### Configure SPDK virtio-blk

```bash
./spdk/scripts/rpc.py spdk_get_version
./spdk/scripts/rpc.py bdev_malloc_create 64 512 -b Malloc1
./spdk/scripts/rpc.py vhost_create_blk_controller --cpumask 0x2 vhost.1 Malloc1
```

### Run qemu with predefined virtio-blk device

```bash
taskset -c 2,3 /usr/libexec/qemu-kvm \
  -cpu host -smp 2 \
  -cdrom init.iso \
  -m 1G -object memory-backend-file,id=mem0,size=1G,mem-path=/dev/hugepages,share=on -numa node,memdev=mem0 \
  -drive file=guest_os_image.qcow2,if=none,id=disk \
  -device ide-hd,drive=disk,bootindex=0 \
  -chardev socket,id=spdk_vhost_blk0,path=/var/tmp/vhost.1 \
  -device vhost-user-blk-pci,chardev=spdk_vhost_blk0,num-queues=2 \
  --nographic
```

Login using fedora/fedora and run few tests

```bash
[fedora@fed21 ~]$ dmesg | grep virtio
[    1.464079] virtio_blk virtio0: [vda] 131072 512-byte logical blocks (67.1 MB/64.0 MiB)

[fedora@fed21 ~]$ ls -l /sys/class/block | grep virtio
lrwxrwxrwx. 1 root root 0 Aug  9 19:00 vda -> ../../devices/pci0000:00/0000:00:04.0/virtio0/block/vda

[fedora@fed21 ~]$ lsblk /dev/vda
NAME MAJ:MIN RM SIZE RO TYPE MOUNTPOINTS
vda  252:0    0  64M  0 disk

[fedora@fed21 ~]$ sudo dd of=/dev/null if=/dev/vda bs=4096 count=4
4+0 records in
4+0 records out
16384 bytes (16 kB, 16 KiB) copied, 0.000137329 s, 119 MB/s

[fedora@fed21 ~]$ sudo dd if=/dev/urandom of=/dev/vda bs=4096 count=4
4+0 records in
4+0 records out
16384 bytes (16 kB, 16 KiB) copied, 0.00232535 s, 7.0 MB/s
```

### Run qemu with HOT PLUG virtio-blk

Start without virtio-blk now but adding QMP management

```bash
taskset -c 2,3 /usr/libexec/qemu-kvm \
  -cpu host -smp 2 \
  -cdrom init.iso \
  -m 1G -object memory-backend-file,id=mem0,size=1G,mem-path=/dev/hugepages,share=on -numa node,memdev=mem0 \
  -drive file=guest_os_image.qcow2,if=none,id=disk \
  -device ide-hd,drive=disk,bootindex=0 \
  -qmp tcp:localhost:4444,server,wait=off \
  --nographic
```

Login using fedora/fedora and verify no virtio-blk devices present

```bash
[fedora@fed21 ~]$ lsblk
NAME   MAJ:MIN RM  SIZE RO TYPE MOUNTPOINTS
sda      8:0    0    5G  0 disk
├─sda1   8:1    0    1M  0 part
├─sda2   8:2    0 1000M  0 part /boot
├─sda3   8:3    0  100M  0 part /boot/efi
├─sda4   8:4    0    4M  0 part
└─sda5   8:5    0  3.9G  0 part /home
                                /
sr0     11:0    1  366K  0 rom
zram0  252:0    0  962M  0 disk [SWAP]
[fedora@fed21 ~]$ dmesg | tail
[    3.894735] fbcon: bochs-drmdrmfb (fb0) is primary device
[    3.898381] Console: switching to colour frame buffer device 128x48
[    3.901540] bochs-drm 0000:00:02.0: [drm] fb0: bochs-drmdrmfb frame buffer device
[    3.941718] RAPL PMU: API unit is 2^-32 Joules, 0 fixed counters, 10737418240 ms ovfl timer
[    3.975023] e1000 0000:00:03.0 eth0: (PCI:33MHz:32-bit) 52:54:00:12:34:56
[    3.975825] e1000 0000:00:03.0 eth0: Intel(R) PRO/1000 Network Connection
[    5.241498] ISO 9660 Extensions: Microsoft Joliet Level 3
[    5.241897] ISO 9660 Extensions: RRIP_1991A
[    5.542201] e1000: eth0 NIC Link is Up 1000 Mbps Full Duplex, Flow Control: RX
[    5.543632] IPv6: ADDRCONF(NETDEV_CHANGE): eth0: link becomes ready
[fedora@fed21 ~]$
```

Hotplug add new virtio-blk device

```bash
[root@Client-3-3Z78MH3 ~]# telnet localhost 4444
Trying ::1...
Connected to localhost.
Escape character is '^]'.
{"QMP": {"version": {"qemu": {"micro": 0, "minor": 2, "major": 6}, "package": "qemu-kvm-6.2.0-11.module+el8.6.0+14707+5aa4b42d"}, "capabilities": ["oob"]}}

{ "execute": "qmp_capabilities" }
{"return": {}}

{ "execute": "query-commands" }
{"return": [{"name": "device_add"}, {"name": "query-pci"}, {"name": "query-acpi-ospm-status"}, {"name": "query-sgx-capabilities"}, {"name": "query-sgx"}, {"n}

{ "execute": "query-pci" }
{"return": [{"bus": 0, "devices": [{"irq_pin": 0, "bus": 0, "qdev_id": "", "slot": 0, "class_info": {"class": 1536, "desc": "Host bridge"}, "id": {"device": }

{"execute": "chardev-add", "id": 3, "arguments": {"id": "spdk_vhost_blk0", "backend": {"type": "socket", "data":{ "addr": {"type": "unix", "data": {"path": "/var/tmp/vhost.1"} } , "server": false } } }}
{"return": {}, "id": 3}

{"execute": "device_add", "id": 4, "arguments": { "driver": "vhost-user-blk-pci", "chardev": "spdk_vhost_blk0"  } }
{"return": {}, "id": 4}
```

See the devices now magically appear

```bash
[   85.303925] pci 0000:00:04.0: [1af4:1001] type 00 class 0x010000
[   85.304701] pci 0000:00:04.0: reg 0x10: [io  0x0000-0x007f]
[   85.305380] pci 0000:00:04.0: reg 0x14: [mem 0x00000000-0x00000fff]
[   85.306221] pci 0000:00:04.0: reg 0x20: [mem 0x00000000-0x00003fff 64bit pref]
[   85.307944] pci 0000:00:04.0: BAR 4: assigned [mem 0x100000000-0x100003fff 64bit pref]
[   85.308898] pci 0000:00:04.0: BAR 1: assigned [mem 0x40000000-0x40000fff]
[   85.309683] pci 0000:00:04.0: BAR 0: assigned [io  0x1000-0x107f]
[   85.310494] virtio-pci 0000:00:04.0: enabling device (0000 -> 0003)
[   85.334818] ACPI: \_SB_.LNKD: Enabled at IRQ 10
[   85.340987] virtio_blk virtio0: [vda] 131072 512-byte logical blocks (67.1 MB/64.0 MiB)
```

Run same tests again

```bash
[fedora@fed21 ~]$ ls -l /sys/class/block | grep virtio
lrwxrwxrwx. 1 root root 0 Aug 10 09:15 vda -> ../../devices/pci0000:00/0000:00:04.0/virtio0/block/vda

[fedora@fed21 ~]$ lsblk /dev/vda
NAME MAJ:MIN RM SIZE RO TYPE MOUNTPOINTS
vda  251:0    0  64M  0 disk

[fedora@fed21 ~]$ sudo dd of=/dev/null if=/dev/vda bs=4096 count=4
4+0 records in
4+0 records out
16384 bytes (16 kB, 16 KiB) copied, 0.00100402 s, 16.3 MB/s

[fedora@fed21 ~]$ sudo dd if=/dev/urandom of=/dev/vda bs=4096 count=4
4+0 records in
4+0 records out
16384 bytes (16 kB, 16 KiB) copied, 0.0034683 s, 4.7 MB/s
```

## virtio-scsi

### Configure SPDK virtio-scsi

```bash
./spdk/scripts/rpc.py spdk_get_version
./spdk/scripts/rpc.py bdev_malloc_create 64 512 -b Malloc2
./spdk/scripts/rpc.py bdev_malloc_create 64 512 -b Malloc3
./spdk/scripts/rpc.py vhost_create_scsi_controller --cpumask 0x1 vhost.0
./spdk/scripts/rpc.py vhost_scsi_controller_add_target vhost.0 0 Malloc2
./spdk/scripts/rpc.py vhost_scsi_controller_add_target vhost.0 1 Malloc3
```

### Run qemu with predefined virtio-scsi device

```bash
taskset -c 2,3 /usr/libexec/qemu-kvm \
  -cpu host -smp 2 \
  -cdrom init.iso \
  -m 1G -object memory-backend-file,id=mem0,size=1G,mem-path=/dev/hugepages,share=on -numa node,memdev=mem0 \
  -drive file=guest_os_image.qcow2,if=none,id=disk \
  -device ide-hd,drive=disk,bootindex=0 \
  -chardev socket,id=spdk_vhost_scsi0,path=/var/tmp/vhost.0 \
  -device vhost-user-scsi-pci,id=scsi0,chardev=spdk_vhost_scsi0,num_queues=2 \
  --nographic
```

Login using fedora/fedora and run few tests

```bash
[fedora@fed21 ~]$ dmesg | grep -i scsi
[    0.314135] SCSI subsystem initialized
[    0.672208] Block layer SCSI generic (bsg) driver version 0.4 loaded (major 244)
[    0.709481] scsi host0: ata_piix
[    0.710172] scsi host1: ata_piix
[    0.874637] scsi 1:0:0:0: CD-ROM            QEMU     QEMU DVD-ROM     2.5+ PQ: 0 ANSI: 5
[    0.876720] sr 1:0:0:0: [sr0] scsi3-mmc drive: 4x/4x cd/rw xa/form2 tray
[    0.894639] sr 1:0:0:0: Attached scsi CD-ROM sr0
[    0.894861] sr 1:0:0:0: Attached scsi generic sg0 type 5
[    0.899092] scsi 1:0:1:0: Direct-Access     ATA      QEMU HARDDISK    2.5+ PQ: 0 ANSI: 5
[    0.903908] sd 1:0:1:0: Attached scsi generic sg1 type 0
[    0.915100] sd 1:0:1:0: [sda] Attached SCSI disk
[    1.473297] scsi host2: Virtio SCSI HBA
[    1.482949] scsi 2:0:0:0: Direct-Access     INTEL    Malloc disk      0001 PQ: 0 ANSI: 5
[    1.484331] scsi 2:0:1:0: Direct-Access     INTEL    Malloc disk      0001 PQ: 0 ANSI: 5
[    1.497173] sd 2:0:0:0: Attached scsi generic sg2 type 0
[    1.499692] sd 2:0:1:0: Attached scsi generic sg3 type 0
[    1.511162] sd 2:0:0:0: [sdb] Attached SCSI disk
[    1.517151] sd 2:0:1:0: [sdc] Attached SCSI disk

[fedora@fed21 ~]$ ls -l /sys/class/block | grep virtio
lrwxrwxrwx. 1 root root 0 Aug 10 15:26 sdb -> ../../devices/pci0000:00/0000:00:04.0/virtio0/host2/target2:0:0/2:0:0:0/block/sdb
lrwxrwxrwx. 1 root root 0 Aug 10 15:26 sdc -> ../../devices/pci0000:00/0000:00:04.0/virtio0/host2/target2:0:1/2:0:1:0/block/sdc

lsblk --output "NAME,KNAME,MODEL,HCTL,SIZE,VENDOR,SUBSYSTEMS" /dev/sdb /dev/sdc

[fedora@fed21 ~]$ lsblk --output "NAME,KNAME,MODEL,HCTL,SIZE,VENDOR,SUBSYSTEMS" /dev/sdb /dev/sdc
NAME KNAME MODEL       HCTL       SIZE VENDOR   SUBSYSTEMS
sdb  sdb   Malloc disk 2:0:0:0     64M INTEL    block:scsi:virtio:pci
sdc  sdc   Malloc disk 2:0:1:0     64M INTEL    block:scsi:virtio:pci

[fedora@fed21 ~]$ sudo dd of=/dev/null if=/dev/sdc bs=4096 count=4
4+0 records in
4+0 records out
16384 bytes (16 kB, 16 KiB) copied, 0.000636195 s, 25.8 MB/s

[fedora@fed21 ~]$ sudo dd if=/dev/urandom of=/dev/sdc bs=4096 count=4
4+0 records in
4+0 records out
16384 bytes (16 kB, 16 KiB) copied, 0.0131856 s, 1.2 MB/s
```

### Run qemu with HOT PLUG virtio-scsi

tbd

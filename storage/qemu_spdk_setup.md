
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

## Configure SPDK virtio-blk

```bash
./spdk/scripts/rpc.py spdk_get_version
./spdk/scripts/rpc.py bdev_malloc_create 64 512 -b Malloc1
./spdk/scripts/rpc.py vhost_create_blk_controller --cpumask 0x2 vhost.1 Malloc1
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

## Run qemu

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

# Prepare Server

## Automation

Then run the [playbook](./ansible) to automate all the lab steps:

```bash
ansible-playbook -i inventory setup.yml
```

## Manual

* Install Ubuntu[^1] 22.04 x64 on the server. ([ubuntu-22.04.1-live-server-amd64.iso](https://releases.ubuntu.com/22.04/))
  * select all default options (unless otherwise noted bellow)
  * on disk setup: disable LVM (optional)
  * on profile setup: put name, servername, username, password all as `opi` for example purposes
  * on ssh setup: enable `install OpenSSH server`
* Install Ubuntu prerequisites

```Shell
sudo apt -y update
sudo apt -y upgrade
sudo apt -y autoremove
sudo apt -y install \
    python3 \
    python3-pip \
    python-is-python3 \
    net-tools \
    curl \
    git \
    make
sudo apt -y install ubuntu-desktop (TODO: remove this depedency)
```

* install Docker (all credits to [Docker manual](https://docs.docker.com/engine/install/ubuntu/) )

```Shell
sudo apt-get -y remove docker docker-engine docker.io containerd runc
sudo apt-get update
sudo apt-get -y install \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg-agent \
    gnupg \
    lsb-release \
    software-properties-common
sudo mkdir -p /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
echo \
    "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
    $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get -y install docker-ce docker-ce-cli containerd.io docker-compose-plugin
sudo docker run hello-world
```

* add your user to docker group

```Shell
sudo usermod -aG docker $USER
```

* enable root (optional)

```Shell
    sudo sed -i "s/#PermitRootLogin prohibit-password/PermitRootLogin yes/g" /etc/ssh/sshd_config
    echo 'root:opi' | sudo chpasswd
    sudo systemctl restart sshd
```

* setup management port configuration using this sample `/etc/netplan/00-installer-config.yaml`:

```code
  ---
  network:
    ethernets:
      ens160:
        dhcp4: false
        dhcp6: false
    bridges:
      br1:
        interfaces: [ens160]
        addresses: [10.36.118.210/24]
        routes:
    - to: default
      via: 10.36.118.1
        mtu: 1500
        nameservers:
          addresses: [4.4.4.4, 8.8.8.8]
        parameters:
    stp: false
    forward-delay: 0
    max-age: 0
        dhcp4: false
        dhcp6: false
    version: 2
```

* check the yaml file is ok (optional)

```Shell
sudo apt -y install yamllint
yamllint /etc/netplan/00-installer-config.yaml
```

* reboot
  * ensure networking is ok
  * this is needed also for the permissions to be update, otherwise next step will fail

## Others

* install ansible

```bash
sudo apt -y install ansible-core
```

* make sure dhcp is running

```bash
systemctl status isc-dhcp-server
```

* Fix pasword-less access

```bash
ssh-keygen
ssh-copy-id 172.22.X.X
```

## Testing

* clone the `opiproject/testing` repository into your working directory:

```Shell
# TBD right now work is happening under opi-poc repo
git clone https://github.com/opiproject/testing
```

* build container

```Shell
docker build --no-cache --tag opi/test-framework:latest ./testing/framework
docker tag opi/test-framework:latest opi/test-framework:1.0.0 # we can chose a versioning schema for the containers
```

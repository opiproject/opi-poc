# StrongSwan Network Function Demo

This directory contains the StrongSwan network function demonstration for Marvell DPUs.

## Demo Videos

### Live Demo Recording

- **Marvell Demo Session**: [Watch Recording](https://marvell.zoom.us/rec/share/urJviRAA8b79PkNy_u42AvR-BMLMWqnwOCtk6cDIAIPeyhTvPtr5ryyyRbhCDHY.jjo_eYZ3rSsJsU1k?startTime=1766088504000)

### OPI Summit Presentation

- **OPI Summit 2024**: [StrongSwan IPsec Offload Demo](https://www.youtube.com/watch?v=SUniRV2Bzr4&t=2497s)
  - Demo section: 41:37 - 1:08:47
  - Full presentation covers OPI architecture and StrongSwan integration with DPU offload

## Related Projects

This demo integrates several open-source projects:

- **[DPU Operator](https://github.com/openshift/dpu-operator)** - OpenShift operator for managing DPU resources and workloads
- **[OPI StrongSwan Bridge](https://github.com/opiproject/opi-strongswan-bridge)** - gRPC bridge for StrongSwan configuration via OPI APIs
  - [Configuration Examples](https://github.com/opiproject/opi-strongswan-bridge/tree/main/conf) - Reference StrongSwan configuration files
- **[godpu](https://github.com/opiproject/godpu)** - Go-based OPI client for IPsec testing and configuration

## Setup Diagram

![StrongSwan NF Setup](setup.jpg)

## Directory Structure

- **pods/** - Kubernetes pod definitions and manifests
- **configs/** - StrongSwan configuration files
- **client/** - Client scripts and utilities

## Overview

StrongSwan NF leverages the OPI StrongSwan Bridge to configure StrongSwan over gRPC. The demo showcases a cloud-native approach to IPSec configuration while utilizing underlying hardware capabilities for offloading cryptographic operations.

## Setup

### Hardware Requirements

**OCP Master Node:**

- Platform: Dell Inc. PowerEdge T550 (Serial: J2LSGR3)
- OS: RHEL with OCP

**OCP Worker/Host:**

- Platform: Supermicro SYS-740P-TR (Serial: A438921X2712857)
- OS: RHCOS with OCP

**DPU:**

- Platform: Marvell crb106-pcie (WA-CN106-A1-PCIE-2P100-R1-029)
- OS: RHEL with microshift

### Network Configuration

The demo uses two IP subnets that can be customized via environment variables:

- **HOST_SUBNET** (default: `10.56.217.0/24`) - Network for host pods
- **EXTERNAL_SUBNET** (default: `20.20.20.0/24`) - External test network

To customize these values, modify the environment variables in the deployment manifests:

**For host pod** ([conf/host-pod.yaml](conf/host-pod.yaml)):

```yaml
env:
- name: EXTERNAL_SUBNET
  value: "20.20.20.0/24"
- name: HOST_GATEWAY
  value: "10.56.217.1"
```

**For network function**, the IP addresses are configured in the run script with default values (`10.56.217.1/24` for host-side and `20.20.20.1/24` for external-side). To customize, modify [pods/nf/run.sh](pods/nf/run.sh).

**For IKE peer** (when running with `podman run`):

```bash
podman run -d --rm --name ike-peer --network host --privileged \
  -e HOST_SUBNET="10.56.217.0/24" \
  -e EXTERNAL_SUBNET="20.20.20.0/24" \
  ike-peer:latest
```

## Build Required Containers

Before building, ensure Docker is installed and you have access to a container registry.

### Container Registry Configuration

The demo requires pushing built images to a container registry. Set your registry address:

```bash
export REGISTRY="<your-registry-address>"
# Example: export REGISTRY="quay.io/myuser"
# Example: export REGISTRY="docker.io/myorg"
```

### Build Host Pod Container

Navigate to the hostpod directory and build the container image:

```bash
cd pods/hostpod/
docker build -t host-pod:latest .
```

Tag and push the image to your container registry:

```bash
# Tag for your registry
docker tag host-pod:latest ${REGISTRY}/host-pod:latest

# Push to registry
docker push ${REGISTRY}/host-pod:latest
```

### Build IKE Peer Container

Navigate to the ikepeer directory and build the container image:

```bash
cd pods/ikepeer/
docker build -t ike-peer:latest .
```

Tag and push the image to your container registry:

```bash
# Tag for your registry
docker tag ike-peer:latest ${REGISTRY}/ike-peer:latest

# Push to registry
docker push ${REGISTRY}/ike-peer:latest
```

### Build Network Function Container

Navigate to the nf directory and build the container image:

```bash
cd pods/nf/
docker build -t strongswan-nf:latest .
```

Tag and push the image to your container registry:

```bash
# Tag for your registry
docker tag strongswan-nf:latest ${REGISTRY}/strongswan-nf:latest

# Push to registry
docker push ${REGISTRY}/strongswan-nf:latest
```

### Update Deployment Manifests

After building and pushing images, update the YAML manifests with your registry:

```bash
# Update host-pod.yaml
sed -i "s|{{.registry}}|${REGISTRY}|g" conf/host-pod.yaml

# Update sfc.yaml
sed -i "s|{{.registry}}|${REGISTRY}|g" conf/sfc.yaml
```

Verify the changes:

```bash
grep "image:" conf/host-pod.yaml conf/sfc.yaml
```

## Demo Execution

### Step 1: Create Crypto VF on DPU

Create a Virtual Function for crypto operations on the DPU:

```bash
# On DPU - Create crypto VF
echo 1 > /sys/bus/pci/devices/0002\:20\:00.0/sriov_numvfs

# Verify VF creation
lspci | grep Enc
```

Expected output:

```text
0002:20:00.0 Encryption controller: Cavium, Inc. Octeon 10 CPT Cryptographic Accelerator, Physical function (rev 51)
0002:20:00.1 Encryption controller: Cavium, Inc. Octeon 10 CPT Cryptographic Accelerator, Virtual function (rev 51)
```

### Step 2: Deploy Pods

#### Deploy OPI Network Function Pod on DPU

```bash
KUBECONFIG=/root/kubeconfig.microshift oc apply -f conf/sfc.yaml
```

**Note:** Ensure the `sfc.yaml` file references the `strongswan-nf:latest` image from your registry.

Verify the OPI NF is running:

```bash
KUBECONFIG="/root/kubeconfig.microshift" oc get pod -n openshift-dpu-operator opi-nf
```

#### Deploy Host Pod on OCP Cluster

```bash
KUBECONFIG=/root/kubeconfig.ocpcluster oc apply -f conf/host-pod.yaml
```

**Note:** Ensure the `host-pod.yaml` file references the `host-pod:latest` image from your registry.

Verify the host pod is running:

```bash
KUBECONFIG="/root/kubeconfig.ocpcluster" oc get pod
```

#### Start StrongSwan Peer Container on External Test System

```bash
# On External Test System
podman run -d --rm --name ike-peer --network host --privileged ike-peer:latest
```

### Step 3: Verify Initial State

Before configuration, verify that IPSec is not yet configured on any system.

#### Check IPSec State on StrongSwan Peer (External Test System)

```bash
# Check IPSec states
podman exec -it ike-peer ip x s

# Check IPSec policies
podman exec -it ike-peer ip x p
```

Both commands should return empty results.

#### Check IPSec State on OPI Network Function

```bash
# Check IPSec states
KUBECONFIG="/root/kubeconfig.microshift" oc exec -it -n openshift-dpu-operator opi-nf -- ip x s

# Check IPSec policies
KUBECONFIG="/root/kubeconfig.microshift" oc exec -it -n openshift-dpu-operator opi-nf -- ip x p
```

Both commands should return empty results.

#### Verify Routes

Check routing tables on all systems:

```bash
# On External Test System
ip route

# On OPI NF
KUBECONFIG="/root/kubeconfig.microshift" oc exec -it -n openshift-dpu-operator opi-nf -- ip route

# On Host Pod
KUBECONFIG="/root/kubeconfig.ocpcluster" oc exec -it host-pod -- ip route
```

### Step 4: Configure IPSec Tunnel

#### Monitor StrongSwan Logs

Start monitoring logs on both OPI NF and the peer:

```bash
# On OPI NF
KUBECONFIG="/root/kubeconfig.microshift" oc exec -it -n openshift-dpu-operator opi-nf -- swanctl --log

# On External Test System (in separate terminal)
podman exec -it ike-peer swanctl --log
```

#### Run OPI Client to Configure StrongSwan

```bash
# Navigate to client directory
cd client/

# Execute configuration via gRPC with default values
go run main.go --addr 20.20.20.1:50151 ipsec test

# Or specify custom traffic selectors to match your network configuration
go run main.go --addr 20.20.20.1:50151 ipsec test \
  --local-ts "10.56.217.0/24" \
  --remote-ts "20.20.20.0/24"
```

See [client/README.md](client/README.md) for detailed client setup instructions.

### Step 5: Verify IPSec Configuration

#### Verify IPSec State on StrongSwan Peer

```bash
# Check IPSec states
podman exec -it ike-peer ip x s

# Check IPSec policies
podman exec -it ike-peer ip x p
```

#### Verify IPSec State on OPI Network Function

```bash
# Check IPSec states
KUBECONFIG="/root/kubeconfig.microshift" oc exec -it -n openshift-dpu-operator opi-nf -- ip x s

# Check IPSec policies
KUBECONFIG="/root/kubeconfig.microshift" oc exec -it -n openshift-dpu-operator opi-nf -- ip x p

# List active Security Associations
KUBECONFIG="/root/kubeconfig.microshift" oc exec -it -n openshift-dpu-operator opi-nf -- swanctl --list-sas
```

### Step 6: Validate Hardware Crypto Offload

#### Check Crypto Hardware Stats Before Traffic Test

```bash
# On DPU
cat /sys/kernel/debug/cn10k/cpt/cpt_pc | grep "CPT instruction requests"
```

Note the current count of CPT (Crypto Processing Tool) instruction requests.

#### Monitor Network Traffic

Start packet capture on the OPI NF to observe encrypted traffic:

```bash
# On OPI NF
KUBECONFIG="/root/kubeconfig.microshift" oc exec -it -n openshift-dpu-operator opi-nf -- tcpdump -i net2
```

#### Get Host Pod IP Address

```bash
KUBECONFIG=/root/kubeconfig.ocpcluster oc get pod host-pod -o json | \
    jq -r '.metadata.annotations["k8s.v1.cni.cncf.io/network-status"]' | \
    jq -r '.[] | select(.interface == "net1") | .ips[0]'
```

#### Send Test Traffic Through IPSec Tunnel

From the External Test System (IP: 20.20.20.20), ping the host pod:

```bash
# On External Test System
ping -c 4 -w 2 -i 0.01 10.56.217.13
```

This traffic should travel through the IPSec tunnel to the OPI NF on the DPU, get decrypted with hardware acceleration, and reach the host pod.

#### Check Crypto Hardware Stats After Traffic Test

```bash
# On DPU
cat /sys/kernel/debug/cn10k/cpt/cpt_pc | grep "CPT instruction requests"
```

Compare this count with the previous value. An increase indicates that crypto operations were offloaded to the DPU hardware accelerator.

## Expected Results

- IPSec tunnel successfully established between External Test System and OPI NF
- Traffic encrypted/decrypted using hardware crypto acceleration on Marvell DPU
- Increased CPT instruction request count after traffic test
- Successful ping responses between External Test System and Host Pod through the tunnel

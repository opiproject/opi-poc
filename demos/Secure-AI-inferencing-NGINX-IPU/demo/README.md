# Demo for Running NGINX on Intel IPU with Red Hat OpenShift

This demo runs in the [OPI lab][opilab], which installs an OCP cluster and an IPU host
with RHEL and microshift. It then proceeds to showcase a NGINX loadbalancer on the IPU
that forwards traffic to AI services on the host.

The setup is similar to Red Hat Summit 2025, see the [demo video][demovideo]. See
also the [blog article about the setup][labsetup].

## IPU Demo Setup

In the [opilab][opilab], we configure an Intel IPU inside host `dh4`. The IPU runs
RHEL 9.6 with microshift and [dpu-operator][dpuop].

There is also the Openshift 4.19 cluster "opicluster".  The host `tgen1` runs 3
virtual machines `opicluster-master-{1,2,3}` as master nodes.  The host `dh4`
is a worker node. The cluster has also dpu-operator installed.

DPU Operator 4.19 is installed via [Operatorhub][rhcatalog].

The cluster was installed and configured manually and using console.redhat.com.
See also the [blog article about the setup][labsetup].

[opilab]: https://github.com/opiproject/lab
[demovideo]: https://www.linkedin.com/posts/ppindell_update-revealed-at-red-hat-summit-last-year-activity-7333338464651792384-e1yW?utm_source=share&utm_medium=member_ios&rcm=ACoAAAARCJUBkVc4lx8FfmYUF_zaKRbmhvyimTk
[labsetup]: https://access.redhat.com/articles/7120276
[dpuop]: https://github.com/openshift/dpu-operator/
[rhcatalog]: https://catalog.redhat.com/en/search?gs=&q=dpu&searchType=all

## Usage

### Demo Script

See the [demo.sh](./demo.sh) script. Run `./demo.sh -h` to get a list and
description of the available commands.

You are also encouraged to review the shell script to understand what each
command does.

You can also `source ./demo.sh` in your local bash and run shell functions
directly. The intended usage is however to run one of the described commands.

### Local Setup

First connect to VPN to [access the OPI lab][accesslab].

[accesslab]: https://github.com/opiproject/lab?tab=readme-ov-file#access-the-lab

You will need a user name and password for that. You will also need to obtain
the password for the root user to access tgen1 host. Ask the OPI lab team.

To access the demo setup, start with the commands

```bash
./demo.sh ssh_copy_id "$SSHKEY"

./demo.sh etc_hosts update

./demo.sh kubeconfigs /tmp

./demo.sh info
```

These install your SSH public key on the hosts, update /etc/hosts and download
kubeconfig files. These steps are optional but useful.

### Console and Cluster Access

The installed setup is a two cluster design. That means, there is a microshift
cluster running on the IPU inside dh4 node and there is an OCP cluster
"opicluster" on tgen1 and dh4 hosts.

If you setup hostnames via `./demo etc_hosts update`, you can access the OCP console in your browser:

<https://console-openshift-console.apps.opicluster.opicluster.opiproject-lab.org>

Login with user "kubeadmin". Get the password via `./demo.sh kubeadmin`. This allows you to look
at the OCP cluster. See [the documentation][console-doc] of the Openshift Web Console.

[console-doc]: https://docs.redhat.com/en/documentation/openshift_container_platform/4.19/html-single/web_console/index

You can also download the kubeconfig files via `./demo.sh kubeconfigs /tmp`. Once you have them, you
can run for example

- `oc --kubeconfig=/tmp/kubeconfig.ocpcluster get node`

- `oc --kubeconfig=/tmp/kubeconfig.microshift-ipu get node`

Alternatively, if you wish, you can also run `./demos.sh oc_ocp ARGS...` and
`./demo.sh oc_ipu ARGS...`. That is the same as calling the `oc` commands with
the right kubeconfig. The `oc` command line tool allows you to investigate the
cluster.

### Access the Demo Hosts

With the [lab plan][opilab] you know the IP addresses and names of the relevant
clusters. There are for example the `tgen1` and `dh4` hosts.

When connected to the VPN, you can directly ssh into those machines. Here it is
especially useful to first register your SSH public key via `./demo.sh
ssh_copy_id "$SSHKEY" all`.

Alternatively, you can also run commands via `./demo.sh exec HOST CMD...` where
`HOST` is one of `tgen1`, `dh4`, `dh4-imc`, `dh4-acc`, `master-1`, `master-2`,
`master-3`.  For example,

- Run an interactive shell session on the IPU with `./demo.sh exec dh4-acc bash`

- Run a command on the remote host like `./demo.sh exec tgen1 bash -c 'hostname; whoami; pwd; uptime'`

### Run the Demo "predict" Example

The core of the demo is a network function that runs on the IPU. This
network function is an NGINX loadbalancer that accepts HTTPS requests and
forwards them to AI example pods running on the host.

Run `./demo.sh predict`. This command calls the example [client script](./predict_images/run.py)
which starts 20 parallel threads. Each thread
will take a random picture of an animal and ask the NGINX loadbalancer to
classify the picture. NGINX on the IPU will accept the HTTPS traffic and
forward an unencrypted request to one of the demo pods on the host. Then the
response is returned. Finally, the script prints the classification result
of the images.

See the NGINX configuration with the upstream pods via
`./demo.sh oc_ipu -n openshift-dpu-operator exec -ti pod/nginx -- cat /etc/nginx/conf.d/91-upstream.conf`.

### Investigate the Demo Setup

Now we can look at the system. Try some of the example commands.

1. Look at the DPU operator's namespace on the OCP and microshift sides

    ```bash
    ./demo.sh oc_ocp -n openshift-dpu-operator get all
    ./demo.sh oc_ipu -n openshift-dpu-operator get all
    ```

1. Inspect the AI Inferencing Pods on the host

    ```bash
    ./demo.sh oc_ocp -n default get pod
    ./demo.sh oc_ocp -n default get pod -o yaml
    ./demo.sh oc_ocp -n default exec -ti pod/resnet50-model-server-1 -- bash
    ```

1. Inspect the NGINX loadbalancer Pod on the IPU

    ```bash
    ./demo.sh oc_ipu -n openshift-dpu-operator get pod/nginx
    ./demo.sh oc_ipu -n openshift-dpu-operator exec -ti pod/nginx -- bash
    ```

1. Check IP configuration and traffic

    Inside the pods, via `oc exec ...`, you can run `ip link`, `ip addr`,
    `ping` and `tcpdump` to investigate what is happening.

    Either start first an interactive shell via
    `./demo.sh oc_ipu -n openshift-dpu-operator exec -ti pod/nginx -- bash`,
    or run as separate commands:

    ```bash
    ./demo.sh oc_ipu -n openshift-dpu-operator exec -ti pod/nginx -- ps aux

    ./demo.sh oc_ipu -n openshift-dpu-operator exec -ti pod/nginx -- ip addr

    POD_IP=$(./demo.sh pod_detect_net1_ip resnet50-model-server-2)
    ./demo.sh oc_ipu -n openshift-dpu-operator exec -ti pod/nginx -- ping -c 3 "$POD_IP"

    ./demo.sh oc_ipu -n openshift-dpu-operator exec -ti pod/nginx -- tcpdump -ni net2
    ```

### Run the Demo "inspect" Example

Run `./demo.sh inspect`, which starts a tmx session with various panes. It will
then execute a series of commands to step through the setup. This is an alternative
to performing those steps yourself, as described [previously](#investigate-the-demo-setup).

At the end of this "inspect" run, you end up with several tmux panes. On the right
hand side you see tcpdump monitoring the traffic inside the AI pods and on the
upper right pane you see the NGINX access log.

On the upper left pane you can run `./demo.sh predict`. This triggers the client
request that you can see in the NGINX access log and the network traffic forwarded
to the pods.

### Redeploy the NGINX Network Function and AI Pods

The way to use DPU Operator in Openshift/Kubernetes is to create the CR (Custom
Resource) following the fields and APIs from the CRD (Custom Resource
Definition). The DPU Operator's manifests describe CRD and its fields and
APIs.

For example, to create the demo pods, we use a YAML configuration that
describes the desired behavior and apply it via `oc create -f -` (or via the
Web Console). You can see how the [demo script](demo.sh) does that by looking
at the `pod_create()` function.

Likewise, to create the NGINX network function on the IPU side, we create a
`ServiceFunctionChain` CR. You can see how that is done in the `sfc_create()`
function. You can also see this service function chain via `./demo.sh oc_ipu
get sfc -A`. Based on that custom resource will the DPU operator create the
NGINX pod.

For our demo setup, we use the public NGINX pod from the Docker Hub Registry.
This does not yet have the right configuration for our load balancing
application.  Due to that, after creating the NGINX network function, our setup
script will access the pod and perform some ad-hoc configuration (see
`nginx_setup()` function).  In a production environment you would use a
specialized container image that already has the right configuration or detects
the configuration automatically. You also probably would not create pods
directly, but rather an Deployment. But for our demo, we configure this more
straight forward via the `./demo.sh redeploy` script.

You can call the `./demo.sh redeploy` script yourself. This will first delete
the `ServiceFunctionChain`, the NGINX pod and the AI demo pods on the host. It
then proceeds to recreate that setup. You can freely call this, especially if
something breaks. You can also look at the script's `do_redeploy()` function
and see all the steps that this involves.

This does not reinstall the OCP cluster or the DPU operator itself.

### Reboot

If something gets broken, and a [redeploy](#redeploy-the-nginx-network-function-and-ai-pods)
cannot fix it, try a reboot via `./demo.sh reboot`. This reboots the involved
machines. This takes a while but is supposed to repair a potentially broken
setup.

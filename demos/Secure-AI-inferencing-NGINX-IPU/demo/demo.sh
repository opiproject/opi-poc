#!/usr/bin/env bash

# https://issues.redhat.com/browse/IIC-727
# https://access.redhat.com/articles/7120276
# https://github.com/opiproject/opi-poc/tree/0de38746cc8ed8eccb3af689a7d96028d0c835c3/demos/Secure-AI-inferencing-NGINX-IPU
# /root/summit on wsfd-advnetlab239.anl.eng.bos2.dc.redhat.com
#
# Needs workaround https://github.com/openshift/dpu-operator/pull/510, hence
# we must use dpu-operator images from quay.io/thaller/dpu-operator-tmp:dpu-daemon-4.19
# See https://docs.google.com/document/d/140XMmFQKQorSDLL0IEmtWUDpMNYoF1RECNC8izBKru0/edit?tab=t.0 aboud
# how the cluster and operator was installed.


NUM_PODS="${NUM_PODS:-3}"

HOST_TGEN1_IP="172.22.1.100"
HOST_DH4_IP="172.22.1.4"

KC_OCP="${KC_OCP:-}"
if [ -z "$KC_OCP" ] ; then
    KC_OCP="$PWD/kubeconfig.ocpcluster"
    if [ ! -f "$KC_OCP" ] && [ -f "/tmp/kubeconfig.ocpcluster"  ] ; then KC_OCP="/tmp/kubeconfig.ocpcluster"  ; fi
    if [ ! -f "$KC_OCP" ] && [ -f "/root/kubeconfig.ocpcluster" ] ; then KC_OCP="/root/kubeconfig.ocpcluster" ; fi
fi
KC_MSH="${KC_MSH:-}"
if [ -z "$KC_MSH" ] ; then
    KC_MSH="$PWD/kubeconfig.microshift"
    if [ ! -f "$KC_MSH" ] && [ -f "/tmp/kubeconfig.microshift"  ] ; then KC_MSH="/tmp/kubeconfig.microshift"  ; fi
    if [ ! -f "$KC_MSH" ] && [ -f "/root/kubeconfig.microshift" ] ; then KC_MSH="/root/kubeconfig.microshift" ; fi
fi

POD_NAME_BASE="resnet50-model-server"

INFERENCE_VENV="/tmp/opilab-demo-inference-venv"

_ECHO_P_INDENT="${_ECHO_P_INDENT:-}"

###############################################################################

_OPI_BASH_X=()
if [[ $- == *x* ]] ; then
    _OPI_BASH_X=( "-x" )
fi

C_RED=''
C_GREEN=''
C_BLUE=''
C_YELLOW=''
C_RESET=''
if [ -t 1 ] ; then
    if [ "${#_OPI_BASH_X[@]}" -eq 0 ] || [ "$_FORCE_COLOR" = 1 ] ; then
        C_RED=$'\033[1;31m'
        C_GREEN=$'\033[1;32m'
        C_BLUE=$'\033[1;34m'
        C_YELLOW=$'\033[1;33m'
        C_RESET=$'\033[0m'
    fi
fi

_OPI_SSH_T=()
if [ -t 1 ] ; then
    _OPI_SSH_T=( "-t" "-o" "LogLevel=ERROR" )
fi

_echo_n() {
    printf '%s' "$*"
}

_echo() {
    printf '%s\n' "$*"
}

_echo_p() {
    printf '%s%s%s\n' "${C_YELLOW}opilab[$(date +"%H:%M:%S.%3N")]>${C_RESET} " "$_ECHO_P_INDENT" "$*"
}

_indent() {
    _ECHO_P_INDENT="$_ECHO_P_INDENT  " \
    "$@"
}

_indent2() {
    _ECHO_P_INDENT="$_ECHO_P_INDENT    " \
    "$@"
}

die() {
    _echo "$C_RED$*$C_RESET"
    exit 1
}

_validate_ip() {
    local ip="$1"

    [[ "$ip" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]
}

_now() {
    local s
    local _
    IFS=$' \t\n' read -r s _ < /proc/uptime
    printf '%d\n' "${s%.*}"
}

_retry_with_timeout() {
    local duration="$1"
    shift || die "missing timeout"
    local deadline="$(( "$(_now)" + duration ))"

    [ "$#" -gt 0 ] || die "missing command"

    while : ; do
        local rc=0
        "$@" || rc="$?"

        if [ "$rc" -eq 0 ] ; then
            return 0
        fi

        if (( "$(_now)" >= deadline )) ; then
            return 124
        fi

        sleep 1
    done
}

usage() {
    local commands

    mapfile -t commands < <(sed -n '/^    case "''$''cmd" in/,/    esac$/ s/^        \([a-zA-Z_0-9]\+\).*/\1/p' "$SCRIPTNAME")

    _echo "Usage $C_GREEN$SCRIPTNAME$C_RESET COMMAND..."
    _echo
    _echo "Inspect the script or run with \`bash -x\` to see what the shell script does."
    _echo "You may source the script and call the shell functions directly."
    _echo
    _echo "Uses:"
    _echo "  export KC_OCP=$C_BLUE$(printf '%q' "$KC_OCP")$C_RESET"
    _echo "  export KC_MSH=$C_BLUE$(printf '%q' "$KC_MSH")$C_RESET"
    _echo

    # Parse the shell script itself, to get out the available commands and
    # doc string.
    python -c 'if 1:
        import os, sys, re

        _, scriptname, c_yellow, c_green, c_reset, *commands = sys.argv

        with open(scriptname, "r") as f:
            script_reversed = f.readlines()
        script_reversed.reverse()
        script_reversed = [s.rstrip("\n") for s in script_reversed]

        def _find_cmd_idx(cmd):
            return next(i for i, line in enumerate(script_reversed) if re.match(f"^(do_)?{cmd}\\(", line))

        def _find_cmd_doc(cmd):
            lst = []
            idx = _find_cmd_idx(cmd) + 1
            while script_reversed[idx].startswith("#"):
                lst.append(re.sub("^# ?", "", script_reversed[idx]))
                idx += 1
            lst.reverse()
            return lst

        def _parse(cmd):
            mode, synopsis, *doc = _find_cmd_doc(cmd)
            mode, prio = mode.rsplit(":", 1)
            mode_prio, mode = mode.split(":", 1)
            return {
                "cmd": cmd,
                "mode": mode,
                "synopsis": synopsis,
                "mode_prio": int(mode_prio),
                "prio": int(prio),
                "doc": doc,
            }

        docs = [_parse(c) for c in commands]
        docs.sort(key=lambda v: (v["mode_prio"], v["prio"]))

        last_mode = ""
        for doc in docs:
            cmd, mode, synopsis, doc = doc["cmd"], doc["mode"], doc["synopsis"], doc["doc"]
            if last_mode != mode:
                last_mode = mode
                print(f"{c_yellow}{last_mode}{c_reset} operations:")
            if synopsis:
                synopsis = f" {synopsis}"
            print(f"  {c_green}{cmd}{c_reset}{synopsis}")
            print("".join(f"    {s}\n" for s in doc), end="")
    ' \
        "$SCRIPTNAME" \
        "$C_YELLOW" \
        "$C_GREEN" \
        "$C_RESET" \
        "${commands[@]}"

    _echo
    _echo "Setup localhost:"
    _echo " - # connect to VPN (\`openconnect --protocol=f5 vpn.opiproject-lab.org\`)"
    _echo " - $(printf '%q' "$SCRIPTNAME0") ssh_copy_id \$SSHKEY"
    _echo " - $(printf '%q' "$SCRIPTNAME0") etc_hosts update"
    _echo " - $(printf '%q' "$SCRIPTNAME0") kubeconfigs /tmp"
    _echo " - $(printf '%q' "$SCRIPTNAME0") info"
}

is_tgen1() {
    [ "$(hostname)" = "tgen1" ] || return 1
    return 0
}

get_podnames() {
    oc_ocp -n default get pods \
        -l app="$POD_NAME_BASE-service" \
        -o custom-columns=:metadata.name --no-headers \
        --sort-by=.metadata.name
}

get_podname() {
    local pod_name="$1"

    case "$pod_name" in
        [0-9] | \
        [0-9][0-9] | \
        pod-[0-9] | \
        pod-[0-9][0-9] )
            get_podnames | sed -n "${pod_name##*-} p"
            ;;
        *)
            _echo "$pod_name"
            ;;
    esac
}

# 2:check cluster:5
# OC_ARGS...
# Call `oc` command for OCP cluster "opicluster".
oc_ocp() {
    oc --kubeconfig="$KC_OCP" "$@"
}

# 2:check cluster:6
# OC_ARGS...
# Call `oc` command for microshift on dh4-acc.
oc_msh() {
    oc --kubeconfig="$KC_MSH" "$@"
}

oc_nginx_exec() {
    oc_msh -n openshift-dpu-operator exec -i pod/nginx -- "$@"
}

_oc_node_pattern_filter() {
    local node_pattern="$1"
    local state_pattern="$2"

    awk \
        -v node_pattern="$1" \
        -v state_pattern="$2" \
        '$1 ~ node_pattern && $2 ~ state_pattern' \
        | sort
}

oc_node_in_state() {
    local node_pattern="$1"
    local mode="$2"
    local state_pattern="$3"
    local out
    local out_all
    local out_match

    # Fetch the node state (in True/False) of nodes with "$node_pattern".
    # Then check whether they are all/any ("$mode") in state "$state_pattern").

    [ "$#" -eq 3 ] || die "invalid arguments"

    out="$(
        oc get nodes \
           -o 'custom-columns=NAME:.metadata.name,STATE:.status.conditions[?(@.type=="Ready")].status' \
           --no-headers \
           2>/dev/null
    )" || return 1

    out_all="$(_oc_node_pattern_filter "$node_pattern" "." <<< "$out")"

    # We must match at least one node, otherwise the result is going to be
    # false.
    [ -n "$out_all" ] || return 1

    out_match="$(_oc_node_pattern_filter "$node_pattern" "$state_pattern" <<< "$out")"

    case "$mode" in
        all)
            [ "$out_match" = "$out_all" ]
            ;;
        any)
            [ -n "$out_match" ]
            ;;
        *)
            die "Invalid mode $mode"
            ;;
    esac
}

oc_node_has_dpu_resources() {
    local node="$1"
    local out
    local num

    [ "$#" -eq 1 ] || die "invalid arguments"

    out="$(oc describe "node/$node")" \
        || return 1

    num="$(sed -n '/^Capacity/,/^[-A-Za-z_]/ s/.*\<openshift\.io\/dpu: *\([0-9]\+\)/\1/p' <<< "$out")" ||
        return 1

    [[ "$num" =~ ^[0-9]+$ ]] || return 1
    [ "$num" -gt 0 ] || return 1
}

wait_ping() {
    local exec_host="$1"
    local timeout_val="$2"
    local host="$3"
    shift 3 || die "missing arguments"

    [ "$#" -eq 0 ] || die "invalid arguments"

    _echo_p "On $exec_host wait for $host to answer ping (with timeout $timeout_val)"

    _EXEC_NOTTY=1 \
    _EXEC_SILENT=1 \
    _exec "$exec_host" \
        timeout "$timeout_val" \
            bash -c \
                "while ! ping -c 1 -W 1 \"$host\" &> /dev/null ; do sleep 1 ; done" \
        && {
            _indent _echo_p "Host $host reachable"
            return 0
        }

    _indent _echo_p "${C_YELLOW}WARNING$C_RESET: Host $exec_host can still not ping $host after $timeout_val"
    return 1
}

wait_ssh() {
    local exec_host="$1"
    local timeout_val="$2"
    local host="$3"
    shift 3 || die "missing arguments"

    [ "$#" -eq 0 ] || die "invalid arguments"

    _echo_p "On $exec_host wait for $host to be reachable via SSH (with timeout $timeout_val)"

    local rc=0

    _EXEC_SILENT=1 \
    _EXEC_CMDSILENT=1 \
    _retry_with_timeout "$timeout_val" \
        _exec "$exec_host" \
            ssh -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=QUIET "$host" true \
                || rc="$?"

    if [ "$rc" -eq 0 ] ; then
        return 0
    fi

    _indent _echo_p "${C_YELLOW}WARNING$C_RESET: Host $exec_host can not ssh to $host after $timeout_val"
    return 1
}

sfc_create() {
    _echo_p "Create SFC \"sfc-test\" in openshift-dpu-operator"
    cat <<'    EOF' | sed 's/^        //' | oc_msh create -f -
        apiVersion: config.openshift.io/v1
        kind: ServiceFunctionChain
        metadata:
          name: sfc-test
          namespace: openshift-dpu-operator
        spec:
          networkFunctions:
          - name: nginx
            image: nginx
            imagePullPolicy: Always
    EOF
}

sfc_wait() {
    _echo_p "Wait for SFC pod with nginx to be ready"
    _retry_with_timeout "60" \
        oc_msh -n openshift-dpu-operator get pod/nginx &>/dev/null \
        || return 1
    oc_msh -n openshift-dpu-operator wait --for=condition=Ready pod/nginx --timeout=15m \
        || return 1
    return 0
}

sfc_delete() {
    oc_msh -n openshift-dpu-operator delete ServiceFunctionChain/sfc-test || true
}

sfc_wait_deleted() {
    oc_msh -n openshift-dpu-operator wait --for="delete" "pod/nginx" --timeout=15m
}

pod_deployment_create() {
    _echo_p "Create resnet deployment \"$POD_NAME_BASE\" in default namespace"

    cat <<EOF | sed 's/^        //' | oc_ocp create -f -
        apiVersion: apps/v1
        kind: Deployment
        metadata:
          name: $POD_NAME_BASE
          namespace: default
        spec:
          replicas: $NUM_PODS
          selector:
            matchLabels:
              app: $POD_NAME_BASE-service
          template:
            metadata:
              annotations:
                k8s.v1.cni.cncf.io/networks: default-sriov-net
              labels:
                app: $POD_NAME_BASE-service
            spec:
              securityContext:
                runAsUser: 0
              nodeSelector:
                dpu: "true"
              volumes:
                - name: model-volume
                  emptyDir: {}
              initContainers:
                - name: model-downloader
                  image: ubuntu:latest
                  securityContext:
                    runAsUser: 0
                  command:
                    - bash
                    - -c
                    - |
                      apt-get update && \
                      apt-get install -y wget ca-certificates && \
                      mkdir -p /models/1 && \
                      wget --no-check-certificate https://storage.openvinotoolkit.org/repositories/open_model_zoo/2022.1/models_bin/2/resnet50-binary-0001/FP32-INT1/resnet50-binary-0001.xml -O /models/1/model.xml && \
                      wget --no-check-certificate https://storage.openvinotoolkit.org/repositories/open_model_zoo/2022.1/models_bin/2/resnet50-binary-0001/FP32-INT1/resnet50-binary-0001.bin -O /models/1/model.bin
                  volumeMounts:
                    - name: model-volume
                      mountPath: /models
              containers:
                - name: ovms
                  image: openvino/model_server:latest
                  args:
                    - "--model_path=/models"
                    - "--model_name=resnet50"
                    - "--port=9000"
                    - "--rest_port=8000"
                  ports:
                    - containerPort: 8000
                    - containerPort: 9000
                  volumeMounts:
                    - name: model-volume
                      mountPath: /models
                  securityContext:
                      privileged: true
                  resources:
                    requests:
                      openshift.io/dpu: '1'
                    limits:
                      openshift.io/dpu: '1'
EOF
}

_pod_detect_net1_ip() {
    local pod_name="$1"
    local rc=0
    local ip

    ip="$(
        oc_ocp -n default get pod "$pod_name" \
            -o jsonpath='{.metadata.annotations.k8s\.v1\.cni\.cncf\.io/network-status}' \
            | jq -r '
                   .[]
                   | select(.interface == "net1")
                   | .ips[0]
            ' \
        )" \
        || rc=1

    if [ "$rc" -eq 0 ] && [ -n "$ip" ] ; then
        printf '%s' "$ip"
        return 0
    fi

    return 1
}

pod_detect_net1_ip() {
    local timeout="$1"
    local pod_name="$2"
    local rc

    pod_name="$(get_podname "$pod_name")"
    rc="0"
    _retry_with_timeout "$timeout" \
        _pod_detect_net1_ip "$pod_name" || rc="$?"
    return "$rc"
}

pods_wait() {
    _echo_p "Wait for resnet pods to be ready"
    oc_ocp -n default wait --for="condition=Ready" -l app="$POD_NAME_BASE-service" pod  --timeout=15m
}

pods_setup() {
    local pod_name
    local pod_names

    if [ -n "$_OPI_DEMO_PODS_SETUP_POD" ] ; then
        mapfile -t pod_names <<< "$_OPI_DEMO_PODS_SETUP_POD"
    else
        mapfile -t pod_names < <(get_podnames)
    fi

    for pod_name in "${pod_names[@]}" ; do
        pod_name="$(get_podname "$pod_name")"
        _echo_p "Install tcpdump and tools in resnet pod $pod_name"
        oc_ocp -n default exec -i "pod/$pod_name" -- /bin/bash -c '
            set -x &&
            apt-get update &&
            apt-get install -y procps iproute2 iputils-ping net-tools tcpdump &&
            true
        '
    done
}

pod_deployment_delete() {
    _echo_p "Delete resnet Deployment $POD_NAME_BASE"
    oc_ocp -n default delete --wait=false "deployment/$POD_NAME_BASE" || :
}

pods_wait_deleted() {
    _echo_p "Wait for resnet pods to be deleted"
    oc_ocp -n default wait --for="delete" -l app="$POD_NAME_BASE-service" pod --timeout=15m
}

nginx_setup_base() {
    _echo_p "Configure pod/nginx with base configuration"

    # Install TLS key and certificate.
    oc_nginx_exec /bin/bash -c 'cat > /etc/nginx/server.crt' < ./server.crt
    oc_nginx_exec /bin/bash -c 'cat > /etc/nginx/server.key' < ./server.key

    # Move conflicting default configuration.
    oc_nginx_exec /bin/bash -c '[ ! -f /etc/nginx/conf.d/default.conf ] || mv /etc/nginx/conf.d/default.conf "/etc/nginx/conf.d/default.conf~"'

    # Install a "91-upstream.conf" configuration file. See function
    # nginx_setup_upstream() which modifies this file.
    cat <<EOF | sed 's/^        //' | oc_nginx_exec /bin/bash -c 'cat > /etc/nginx/conf.d/91-upstream.conf'
        upstream model_servers {
            # Add "server $IP:9000;" entries to pods.
        }
EOF

    # Install a "90-base.conf" file for proxying requests to the model servers.
    cat <<'    EOF' | sed 's/^        //' | oc_nginx_exec /bin/bash -c 'cat > /etc/nginx/conf.d/90-base.conf'
        log_format main2 '$remote_addr:$remote_port > '
                         '[$time_iso8601.$msec] "$request" '
                         '$status $body_bytes_sent "$http_referer" '
                         '"$http_user_agent" "$http_x_forwarded_for" '
                         'rt=$request_time';

        server {
            access_log /var/log/nginx/access.log main2;

            listen *:443 ssl http2;

            server_name n1.nginx.ipu.opicluster.opiproject-lab.org;
            ssl_certificate /etc/nginx/server.crt;
            ssl_certificate_key /etc/nginx/server.key;

            # proxy gRPC → your upstream
            location / {
                # these must match your gRPC host header
                grpc_set_header   Host              $http_host;
                grpc_set_header   X-Real-IP         $remote_addr;
                grpc_set_header   X-Forwarded-For   $proxy_add_x_forwarded_for;
                grpc_pass         grpc://model_servers;
            }
            location /nginx_status {
              stub_status on;
            }
        }
    EOF
}

nginx_setup_ipaddr() {
    # Install a few useful packages and configure static IP addresses
    # inside the pod.
    # - 10.56.217.3/24: IP address to talk with AI demo pods (see also
    #   nad_ip_range() function).
    # - 172.16.3.200/24: this IP address is reachable from external. HTTP
    #   requests on this address will be load balanced by the nginx proxy.
    _echo_p "Configure pod/nginx with IP addresses"
    oc_nginx_exec bash -c "
        set -x && \\
        apt-get update && \\
        apt-get install -y procps iproute2 iputils-ping net-tools tcpdump && \\
        ip addr flush dev net1 && \\
        ip addr flush dev net2 && \\
        ip addr add 10.56.217.3/24 dev net2 && \\
        ip addr add 172.16.3.200/24 dev net2 && \\
        ip addr
    "
}

nginx_setup_upstream_ip() {
    local ip="$1"

    [ "$#" -eq 1 ] || die "invalid arguments"
    _validate_ip "$ip" || die "Invalid IP address format: $ip"

    oc_nginx_exec sed -i \
        -e "s/^}/    server $ip:9000;\\n}/" \
        /etc/nginx/conf.d/91-upstream.conf
}

nginx_setup_upstream() {
    local pod_name
    local ip

    # Detect and configure the IP addresses of the AI pods as upstream for the
    # nginx load balancer.
    #
    # See also nad_ip_range(). On the host side the network attachment
    # definition is configured to hand out a certain 10.56.217.0/24 range to
    # the pods.

    while IFS= read -r pod_name ; do
        ip="$(pod_detect_net1_ip 180 "$pod_name")" \
            || { _echo_p "${C_RED}ERROR${C_RESET}: Failure to detect IP address in $pod_name"; return 1; }

        if ! _validate_ip "$ip" ; then
            _echo_p "${C_RED}ERROR${C_RESET}: Invalid IP address detected for pod $pod_name: $ip"
            return 1
        fi

        _echo_p "Configure pod/nginx with upstream IP address $ip for $pod_name"
        nginx_setup_upstream_ip "$ip"
    done < <(get_podnames)
}

nginx_setup_reload() {
    _echo_p "Reload pod/nginx pod"
    oc_nginx_exec /bin/bash -c "nginx -t && nginx -s reload"
}

nginx_setup() {
    # We run the unmodified nginx pod from Docker Container Registry.
    #
    # For the demo setup, we must configure it in a suitable way. Primarily the
    # nginx configuration.
    _echo_p "Setup nginx pod on dh4-acc"
    (
        # shellcheck disable=SC2030
        _ECHO_P_INDENT="$_ECHO_P_INDENT  "

        nginx_setup_base
        nginx_setup_ipaddr
        nginx_setup_upstream
        nginx_setup_reload
    )
}

nginx_wait() {
    wait_ping tgen1 30 172.16.3.200
}

# 2:check cluster:7
#
# Reconfigure the nginx pod on dh4-acc. This is a subset of redeploy.
do_nginx_setup() {
    nginx_setup
    nginx_wait
}

# 2:check cluster:8
# [POD]
# Reconfigure the resnet pods in the OCP cluster. This is a subset of redeploy.
# If no pod is given, all of them are reconfigured.
do_pods_setup() {
    [ "$#" -le 1 ] || die "Invalid arguments"
    _OPI_DEMO_PODS_SETUP_POD="$1" pods_setup
}

# 2:check cluster:9
# POD
# Lookup the IP address on the secondary network net1 of the AI pod on the OCP
# side. This is the upstream IP address for the nginx load balancer.
do_pod_detect_net1_ip() {
    [ "$#" -le 1 ] || die "Invalid arguments"
    [ "$#" -eq 1 ] || die "Missing pod name"
    pod_detect_net1_ip 0 "$1"
}

_cmd_nginx_macs() {
    oc_msh -n openshift-dpu-operator exec pod/nginx -- bash -c '(ip link show net1; ip link show net2) | sed -n "s/^.*link\/ether \+\([^ ]\+\) \+.*$/\1/p"'
}

# Check interface rx/tx statistics on dh4-ipu. Call via `watch`
# shellcheck disable=SC2120
_cmd_ipu_statistics() {
    local watch="$1"
    shift || :
    local macs=( "$@" )
    local macs_re
    local cmd

    if [ "${#macs[@]}" -eq 0 ] ; then
        mapfile -t macs <<< "$(_cmd_nginx_macs)"
    fi

    macs_re="$(printf '%s\\|' "${macs[@]}")"
    macs_re="${macs_re%\\|}"

    cmd="$(cat <<EOF | sed 's/^ */ /' | tr -d '\n'
            while read -r mac vsi_id ; do
                printf "%s\\n" "Interface \$mac (\$vsi_id)" ;
                cli_client -q -s -v "\$vsi_id" | sed -n "s/.*\\(ingress\\|egress\\) packet:.*/  \\0/p" ;
            done < <(cli_client -c -q | sed -n "/$macs_re/ s/.*vsi_id: \\([0-9xa-fA-F]\\+\\) .*mac addr: \\([^ ]*\\).*/\\2 \\1/p") ;
EOF
    )"

    local args=( bash -c "$cmd" )
    local _notty=1
    if [ "$watch" = 1 ] ; then
        args=( watch -x -n 0.3 "${args[@]}" )
        _notty=0
    fi

    _EXEC_NOTTY="$_notty" \
    _EXEC_SILENT=1 \
    _exec dh4-imc "${args[@]}"
}

# 2:check cluster:1
# HOST [-s] [-T] [CMD...]
# Run command on host. First argument is the host followed by the command.
# If no command is given, it defaults to bash.
#
# Supported hosts are "localhost", "tgen1", "opicluster-master-[123]", "dh4",
# "dh4-acc", dh4-imc", "mgmt". In those case, the command uses SSH via tgen1.
#
# Also supported is "dh4-acc-nginx" and "dh4-pod-[123]", in which we use
# oc-exec to run a command inside the pods.
#
# The flags "-s" and "-T" make the command silent and without TTY, respectively.
# Also, OPI_DEMO_EXEC_{SILENT,NOTTY}=1 environment is honored.
do_exec() {
    _EXEC_SILENT="$OPI_DEMO_EXEC_SILENT" \
    _EXEC_NOTTY="$OPI_DEMO_EXEC_NOTTY" \
    _exec "$@"
}

_exec() {
    local host
    local args=()
    local args_cmd
    local name

    while : ; do
        case "$1" in
            -s)
                local _EXEC_SILENT=1
                shift
                ;;
            -T)
                local _EXEC_NOTTY=1
                shift
                ;;
            *)
                break
                ;;
        esac
    done

    host="$1"
    [ -n "$host" ] || die "missing name of target host for exec"

    shift

    if [ "$#" -gt 0 ] ; then
        args=( "$@" )
    else
        args=( bash )
    fi

    case "$host" in
        master-[123]|opicluster-master-[123])
            host="opicluster-master-${host##*-}"
            args=( sudo -- bash -c "cd &&$(printf ' %q' "${args[@]}")" )
            ;;
        dh4)
            args=( sudo -- bash -c "cd &&$(printf ' %q' "${args[@]}")" )
            ;;
    esac

    args_cmd="$(printf '%q ' "${args[@]}")"
    args_cmd="${args_cmd% }"

    if [ "$_EXEC_SILENT" != 1 ] ; then
        _echo_p "Run on $host: $args_cmd"
    fi

    local _ssh_t=()
    local _oc_t=()
    if [ "$_EXEC_NOTTY" != 1 ] ; then
        _ssh_t=( "${_OPI_SSH_T[@]}" )
        if [ "${#_OPI_SSH_T[@]}" -gt 0 ] ; then
            _oc_t=( '-t' )
        fi
    fi

    local ssh_cmd_tgen1=()
    if ! is_tgen1 ; then
        # If we are not on tgen1, we progably will first remote the call via
        # ssh. that has two reasons:
        # - some hosts (dh4-acc) are not directly accessible via the SSH. We
        #   Need to go through tgen1
        # - by going through tgen1, we only need authentication via that host.
        #   tgen1 can from there password-less login to the other hosts.
        ssh_cmd_tgen1=( ssh "${_ssh_t[@]}" "root@$HOST_TGEN1_IP" )
    fi

    case "$host" in
        localhost)
            ssh_cmd_tgen1=()
            ;;
        tgen1)
            if [ "${#ssh_cmd_tgen1[@]}" -gt 0 ] ; then
                args=( "${ssh_cmd_tgen1[@]}" "$args_cmd" )
            fi
            ssh_cmd_tgen1=()
            ;;
        opicluster-master-[123])
            args=( ssh "${_ssh_t[@]}" "core@192.168.122.$(( "${host##*-}" + 1 ))" "$args_cmd" )
            ;;
        dh4-acc)
            args=( ssh "${_ssh_t[@]}" "root@172.16.3.16" "$args_cmd" )
            ;;
        dh4-imc)
            args=( ssh "${_ssh_t[@]}" -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null -o LogLevel=ERROR "root@172.22.4.4" "$args_cmd" )
            ;;
        dh4)
            args=( ssh "${_ssh_t[@]}" "core@$HOST_DH4_IP" "$args_cmd" )
            ;;
        mgmt)
            args=( ssh "${_ssh_t[@]}" "root@172.22.0.1" "$args_cmd" )
            ;;
        dh4-pod-[1-9])
            name="$(get_podname "${host##*-}")"
            args=( oc --kubeconfig="/root/kubeconfig.opicluster" exec -n default "${_oc_t[@]}" -i "pod/$name" -- "${args[@]}" )
            ;;
        dh4-acc-nginx)
            args=( oc --kubeconfig="/root/kubeconfig.dh4-acc" exec -n openshift-dpu-operator "${_oc_t[@]}" -i pod/nginx -- "${args[@]}" )
            ;;
        *)
            die "Invalid host $host for exec"
            ;;
    esac

    if [ "${#ssh_cmd_tgen1[@]}" -gt 0 ] ; then
        local cmd2

        cmd2="$(printf '%q ' "${args[@]}")"
        cmd2="${cmd2% }"
        args=( "${ssh_cmd_tgen1[@]}" "$cmd2" )
    fi

    (
        cd "$ORIGINAL_PWD"

        if [ "$_EXEC_CMDSILENT" = 1 ] ; then
            "${args[@]}" &>/dev/null
        else
            "${args[@]}"
        fi
    )
}

# 2:check cluster:2
# DEMO_CMD...
# Run demo command on tgen1 host. This first copies the demo over to tgen1
# using rsync. For example, use `remote redeploy` to call the redeploy command
# of this script on the tgen1 host.
do_remote() {
    local tdir="/tmp/opilab-demo"

    _echo_p "Copy demo to tgen1 host ($tdir)"
    rsync -a --info=NAME \
        --delete \
        --exclude="*.swp" \
        --exclude=".mypy_cache/" \
        --exclude="__pycache__/" \
        --exclude="kubeconfig*" \
        . \
        "root@$HOST_TGEN1_IP:$tdir/" \
        ;

    local envs=()
    if [ -n "$OPI_DEMO_SLEEP_TIME" ] ; then
        envs+=( "OPI_DEMO_SLEEP_TIME=$OPI_DEMO_SLEEP_TIME" )
    fi

    # shellcheck disable=SC2031
    _exec tgen1 \
        /usr/bin/env \
        _ECHO_P_INDENT="$_ECHO_P_INDENT  " \
        "${envs[@]}" \
        bash "${_OPI_BASH_X[@]}" "$tdir/$SCRIPTNAME" "$@"
}

exec_podman_k8stft_on_host() {
    local host="$1"

    shift

    _exec "$host" podman run --rm -i --privileged --network=host -v '/:/host' ghcr.io/ovn-kubernetes/kubernetes-traffic-flow-tests:latest "$@"
}

# 2:check cluster:3
# [TCPDUMP_ARGS...]
# Call tcpdump on dh4 host. Additional arguments passed to tcpdump.
do_tcpdump_dh4() {
    _echo_p "Run tcpdump for dh4 host with arguments: $*"
    exec_podman_k8stft_on_host dh4 \
        tcpdump "$@"
}

# 2:check cluster:4
# [TCPDUMP_ARGS...]
# Call tcpdump on resnet pod's net1 interface. First argument is the pod name
# (or number). Further arguments are passed to tcpdump
do_tcpdump_pod() {
    local pod_name="$1"
    local args=()
    local cmd

    [ -n "$pod_name" ] || die "Requires podname as first parameter"

    shift

    pod_name="$(get_podname "$pod_name")"

    args=( -i net1 -n "$@" )

    local cmd
    cmd="$(printf '%q ' "${args[@]}")"
    cmd="${cmd% }"

    _echo_p "Run tcpdump inside pod $pod_name with arguments: tcpdump $cmd"

    oc_ocp -n default exec -i "pod/$pod_name" -- /bin/bash -c '
        set +x &&
        ip -4 addr &&
        tcpdump '"$cmd"'
    '
}

show_info() {
    _echo_p "Information about OCP cluster for $KC_OCP"
    oc_ocp get node
    oc_ocp -n openshift-dpu-operator get all
    oc_ocp -n default get all

    _echo_p "Information about microshift on dh4 for $KC_MSH"
    oc_msh get node
    oc_msh -n openshift-dpu-operator get all
    oc_msh -n default get all

    pod="$(oc_msh get pods -A -o name | grep 'vsp-p4' | head -n1)"
    _echo_p "Show ovs-vsctl inside VSP pod $pod on DPU side"
    oc_msh -n openshift-dpu-operator exec -i "$pod" -- /opt/p4/p4-cp-nws/bin/ovs-vsctl show || true
    oc_msh -n openshift-dpu-operator exec -i "$pod" -- /opt/p4/p4-cp-nws/bin/p4rt-ctl dump-entries br0 | head -n10 || true
    _cmd_ipu_statistics
}

cleanup_all() {
    _echo_p "Delete SFC"
    sfc_delete
    pod_deployment_delete

    _echo_p "Wait for SFC and Pods to be gone"
    sfc_wait_deleted \
        || _echo_p "  Failed"
    pods_wait_deleted \
        || _echo_p "  Failed"
}

py_pip_install() {
    pip install --upgrade pip
    pip install \
        black \
        flake8 \
        mypy \
        numpy \
        opencv-python \
        ovmsclient \
        ;
}

py_activate_venv() {
    local setup=0

    if declare -F deactivate 2>/dev/null ; then
        deactivate
    fi

    if [ ! -f "$INFERENCE_VENV/bin/activate" ] ; then
        _echo_p "Create Python venv at $INFERENCE_VENV"
        "${PYTHON:-python}" -mvenv "$INFERENCE_VENV"
        setup=1
    fi

    _echo_p "Activate Python venv at $INFERENCE_VENV"
    # shellcheck disable=1091
    source "$INFERENCE_VENV/bin/activate"

    if [ "$setup" = 1 ] ; then
        _echo_p "Install Python packages in venv"
        py_pip_install
    fi
}

###############################################################################

nad_ip_range() {
    local value
    local new_value

    # We want to manually configure the 10.56.217.3 address in the nginx NF.
    # To avoid clashes, restrict the range that gets automatically assigned to
    # the resnet pods. We do that by patching the network attachment definition.
    #
    # See also https://github.com/openshift/dpu-operator/pull/514

    _echo_p "Patch net-attach-def/default-sriov-net to assign restricted IP range"

    value="$(oc_ocp -n default get net-attach-def default-sriov-net -o jsonpath='{.spec.config}')"
    new_value="$(
        printf '%s' "$value" \
            | jq '.ipam.subnet = "10.56.217.0/24"' \
            | jq '.ipam.rangeStart = "10.56.217.11"' \
            | jq '.ipam.rangeEnd = "10.56.217.239"' \
            | jq -c \
            | jq -Rs .
    )"
    oc_ocp -n default patch net-attach-def default-sriov-net --type=merge \
          -p "{\"spec\":{\"config\":$new_value}}"
}

_wait_node_ready_ocp_masters() {
    _echo_p "Wait for master nodes in opicluster to be ready"
    KUBECONFIG="$KC_OCP" \
    _retry_with_timeout 180 \
        oc_node_in_state '^opicluster-master-[0-9]+$' all True \
        || { _echo_p "  not ready"; return 1; }
}

_wait_node_ready_ocp_dh4() {
    _echo_p "Wait for worker node dh4 in opicluster to be ready"
    KUBECONFIG="$KC_OCP" \
    _retry_with_timeout 180 \
        oc_node_in_state '^dh4$' all True \
        || { _echo_p "  not ready"; return 1; }
}

_wait_node_ready_microshift() {
    _echo_p "Wait for nodes in microshift on IPU to be ready"
    KUBECONFIG="$KC_MSH" \
    _retry_with_timeout 180 \
        oc_node_in_state '^dh4-acc$' all True
}

_wait_is_up_tgen1() {
    wait_ssh localhost 300 "root@$HOST_TGEN1_IP"
    _wait_node_ready_ocp_masters
}

_wait_is_up_dh4() {
    wait_ssh tgen1 300 "core@$HOST_DH4_IP"
    _wait_node_ready_ocp_dh4
}

_wait_is_up_dpu() {
    wait_ssh tgen1 300 "root@172.22.4.4" || return 1
    wait_ssh tgen1 180 "root@172.16.3.16" || return 1
    wait_ssh localhost 60 "root@172.22.1.104" || return 1
    _wait_node_ready_microshift || return 1
}

_opidemo_rebootcmd='( nohup bash -c "sleep 3; reboot" & sleep 1 ; exit ) &>/dev/null'
_opidemo_poweroffcmd='( nohup bash -c "sleep 3; poweroff" & sleep 1 ; exit ) &>/dev/null'

_reboot_tgen1_tgen1() {
    _echo_p "reboot tgen1"

    _EXEC_NOTTY=1 \
    _indent _exec tgen1 bash -c "$_opidemo_rebootcmd"

    sleep 60

    _indent wait_ssh localhost 300 "root@$HOST_TGEN1_IP"

    _EXEC_NOTTY=1 \
    _indent _exec tgen1 /etc/nftables.sh

    _indent _wait_is_up_tgen1
}

_reboot_dh4_acc_one() {
    _echo_p "Reboot ACC inside dh4"

    _indent _echo_p "power off ACC inside dh4"
    _EXEC_NOTTY=1 \
    _EXEC_CMDSILENT=1 \
    _indent2 _exec dh4-acc \
        bash -c "$_opidemo_poweroffcmd" \
        || true

    sleep 20

    _indent _echo_p "reboot IMC and ACC inside dh4"

    _EXEC_NOTTY=1 \
    _EXEC_CMDSILENT=1 \
    _indent2 _exec dh4-imc \
        bash -c "$_opidemo_poweroffcmd" \
        || { _echo_p "  ${C_YELLOW}power off command for dh4-imc failed$C_RESET"; return 1 ; }

    sleep 60

    _indent _wait_is_up_dpu

    _indent _echo_p "Give additional time for ACC to settle"
    sleep 240
}

_reboot_dh4_acc() {
    local try_count=0

    while true ; do
        try_count="$((try_count + 1))"
        if [ "$try_count" -gt 10 ] ; then
            die "Cannot bring up IPU inside dh4"
        fi

        if ! _reboot_dh4_acc_one ; then
            _indent _echo_p "${C_YELLOW}Failed to bring up IPU.$C_RESET Retry #$try_count..."
            continue
        fi

        return 0
    done
}

_reboot_dh4_one() {
    _echo_p "reboot worker node dh4"

    _EXEC_NOTTY=1 \
    _indent _exec dh4 bash -c "$_opidemo_rebootcmd"

    sleep 60

    _indent _wait_is_up_dh4

    sleep 15
}

_reboot_dh4() {
    local try_count=0

    while true ; do
        try_count="$((try_count + 1))"
        if [ "$try_count" -gt 10 ] ; then
            die "Cannot bring up dh4 host"
        fi

        if ! _reboot_dh4_one ; then
            _indent _echo_p "${C_YELLOW}Failed to bring up dh4.$C_RESET Retry #$try_count..."
            continue
        fi

        return 0
    done
}

_reboot_dh4_reload_idpf_reload_idpf() {
    local pod_name
    local pod_names

    _echo_p "reload idpf driver on dh4"
    _indent _exec dh4 bash -c "rmmod idpf && sleep 10 && modprobe idpf" \
        || { _indent _echo_p "${C_YELLOW}failed$C_RESET"; return 1; }

    _indent _echo_p "restart dpu-daemon"
    pod_names="$(
        oc_ocp -n openshift-dpu-operator get pod -o 'custom-columns=NAME:.metadata.name' --no-headers \
            | sed -n 's/^\(dpu-daemon-\|vsp-\)/&/p'
    )"
    while IFS= read -r pod_name ; do
        oc_ocp -n openshift-dpu-operator delete "pod/$pod_name"
    done <<< "$pod_names"

    sleep 10
}


# 1:main:1
# [full]
# Reboot machines and redeploy.
#
# This should be always safe to do, especially in an attempt to make the system
# healthy again.
#
# With "full", also reboot provisioning host "tgen1" and the master VM nodes.
#
# Warning: this will take a while.
do_reboot() {
    local full=0

    [ "$#" -le 1 ] || die "invalid arguments"
    if [ "$#" -ge 1 ] ; then
        case "$1" in
            full)
                full=1
                ;;
            *)
                die "invalid argument"
                ;;
        esac
    fi

    cleanup_all

    if [ "$full" -eq 1 ] ; then
        _echo_p "reboot (tgen1, dh4, dh4-acc)..."
        if is_tgen1 ; then
            die "Cannot run this command on tgen1, because it will reboot itself."
        fi
        _indent _reboot_tgen1_tgen1
    else
        _echo_p "reboot (dh4, dh4-acc)..."
    fi

    local try_count=0

    while true ; do
        try_count="$((try_count + 1))"

        _indent _reboot_dh4
        _indent _reboot_dh4_acc
        _indent _reboot_dh4_reload_idpf_reload_idpf

        _indent _echo_p "Wait for worker node dh4 to have \"openshift.io/dpu\" resources"

        local rc=0
        KUBECONFIG="$KC_OCP" \
        _retry_with_timeout 120 oc_node_has_dpu_resources 'dh4' \
         || rc=1

        if [ "$rc" -eq 0 ] ; then
            break
        fi

        if [ "$try_count" -gt 16 ] ; then
            die "worker node dh4 does not get into a working state"
        fi

        _indent2 _echo_p "${C_YELLOW}Worker node dh4 does not have \"openshift.io/dpu\" resources. Retry$C_RESET"
    done

    _EXEC_NOTTY=1 \
    _indent _exec tgen1 /etc/nftables.sh

    _indent do_redeploy

    _indent _echo_p "Show info"
    _indent do_info
}

# 1:main:2
#
# Delete the pods and SFC (if any) and redeploy them from scratch.  As we
# manually adjust the running pods for the demo, we need to redeploy after a
# pod gets deleted or restarted.
#
# Consider running `remote redeploy`, if you didn't setup your system
# according to "Setup localhost".
do_redeploy() {
    _echo_p "Redeploy nginx and resnet"

    (
        # shellcheck disable=SC2031
        _ECHO_P_INDENT="$_ECHO_P_INDENT  "

        nad_ip_range

        cleanup_all

        sfc_create
        pod_deployment_create

        sfc_wait
        pods_wait

        pods_setup
        nginx_setup

        nginx_wait
    )
}

# 1:main:3
#
# Test the nginx load balancer by asking the resnet pods to classify images
# of animals. This is the demo usage.
do_predict() {
    if ! is_tgen1 ; then
        # The nginx IP is only accessible from tgen1. Remote the call.
        do_remote predict
        return 0
    fi

    py_activate_venv

    _echo_p "${C_GREEN}PREDICT IMAGES...${C_RESET}"
    ./predict_images/run.py ./predict_images/images/

    deactivate
}

# 1:main:4
#
# Print various information.
do_info() {
    for h in localhost tgen1 dh4 dh4-acc dh4-imc ; do
        _EXEC_NOTTY=1 \
        _EXEC_SILENT=1 \
        _exec "$h" \
            sh -c 'printf "HOST %-12s : '"$C_GREEN"'%-12s'"$C_RESET"' - %s\n" "'"$h"'" "'\$'(hostname)" "'\$'(uptime)"'
    done

    _echo "Console: https://console-openshift-console.apps.opicluster.opicluster.opiproject-lab.org"
    _echo "     User \"kubeadmin\". Get secret via \`$SCRIPTNAME0 kubeadmin\`."
    show_info
}

# 1:main:5
# [local]
# Run inspect.sh script for a demo using tmux windows.
#
# By default runs via `remote inspect` on tgen1. Specify "local"
# argument to force running locally.
do_inspect() {
    local _local=0

    [ "$#" -le 1 ] || die "invalid arguments"
    if [ "$#" -ge 1 ] ; then
        case "$1" in
            local)
                _local=1
                ;;
            *)
                die "invalid argument"
                ;;
        esac
    fi

    if [ "$_local" != 1 ] && ! is_tgen1 ; then
        do_remote inspect
        return 0
    fi

    /usr/bin/env bash "./inspect.sh"
}

# 3:extra helper:1
#
# Regenerate TLS certificate for the nginx pod. This overwrites
# "server.{key,crt}" in the script directory and is used during redeploy.
do_tls_generate_keys() {
    _echo_p "Generate TLS keys $PWD/server.key and $PWD/server.crt"

    openssl genpkey -quiet -algorithm RSA -out server.key -pkeyopt rsa_keygen_bits:2048

    cat <<'    EOL' | sed 's/^        //' > /tmp/opilab-demo-san.cnf
        [req]
        distinguished_name = req_distinguished_name
        req_extensions = v3_req
        prompt = no

        [req_distinguished_name]
        CN = myserver.local

        [v3_req]
        keyUsage = digitalSignature, keyEncipherment, dataEncipherment
        extendedKeyUsage = serverAuth
        subjectAltName = @alt_names

        [alt_names]
        DNS.1 = n1.nginx.ipu.opicluster.opiproject-lab.org
        IP.1 = 172.16.3.200
    EOL

    openssl req -x509 \
        -days $((3*365)) \
        -key server.key \
        -out server.crt \
        -config /tmp/opilab-demo-san.cnf \
        -extensions v3_req

    rm -rf /tmp/opilab-demo-san.cnf
}

# 3:extra helper:2
# [update]
# Show entries for your /etc/hosts to be able to locally access console. Pass
# "update" to update the file using sudo.
do_etc_hosts() {
    set +x
    local content
    local update=0

    [ "$#" -le 1 ] || die "invalid arguments"
    if [ "$#" -ge 1 ] ; then
        case "$1" in
            u|update)
                update=1
                ;;
            *)
                die "invalid argument"
                ;;
        esac
    fi

    content="$(cat <<EOF
# Created for opicluster demo. With these we can access
# https://console-openshift-console.apps.opicluster.opicluster.opiproject-lab.org
172.22.1.99     api.opicluster.opicluster.opiproject-lab.org
172.22.1.101    oauth-openshift.apps.opicluster.opicluster.opiproject-lab.org
172.22.1.101    console-openshift-console.apps.opicluster.opicluster.opiproject-lab.org
172.22.1.101    grafana-openshift-monitoring.apps.opicluster.opicluster.opiproject-lab.org
172.22.1.101    thanos-querier-openshift-monitoring.apps.opicluster.opicluster.opiproject-lab.org
172.22.1.101    prometheus-k8s-openshift-monitoring.apps.opicluster.opicluster.opiproject-lab.org
172.22.1.101    alertmanager-main-openshift-monitoring.apps.opicluster.opicluster.opiproject-lab.org
EOF
    )"

    if [ "$update" != 1 ] ; then
        cat <<EOF2
# Run the following command to update /etc/hosts
cat <<EOF | sudo tee -a /etc/hosts

$content
EOF2
        return 0
    fi

    if grep -q console-openshift-console.apps.opicluster.opicluster.opiproject-lab.org /etc/hosts ; then
        _echo "/etc/hosts seems already configured. SKIP"
        _echo
        _echo '...'
        grep -C10 console-openshift-console.apps.opicluster.opicluster.opiproject-lab.org /etc/hosts
    else
        cat <<EOF | sudo tee -a /etc/hosts

$content
EOF
        _echo "/etc/hosts updated"
        _echo
        _echo '...'
        tail -n 15 /etc/hosts
    fi
}

# 3:extra helper:3
# SSHKEY [WHERE]
# Call ssh-copy-id to install SSH key in the demo cluster. Pass the SSH key.
#
# By default, the SSHKEY is only installed on "tgen1" host. This usually
# suffices. If you want, set the WHERE parameter to "all" to install it on all
# related hosts. Alternatively, WHERE can be one of the host names
# tgen1, opicluster-master-[123], dh4, dh4-acc.
do_ssh_copy_id() {
    local sshkey="$1"
    local rc=0
    local i

    [ -f "$sshkey" ] || die "Command requires file to SSH key to install"

    shift

    local where='tgen1'
    case "$1" in
        all|a|-a)
            where='all'
            shift
            ;;
        tgen1|opicluster-master-[123]|dh4|dh4-acc)
            where="$1"
            shift
            ;;
        master-[123])
            where="opicluster-master-${1##*-}"
            shift
            ;;
    esac

    [ "$#" -eq 0 ] || die "Invalid arguments"

    if [ "$where" = 'all' ] || [ "$where" = 'tgen1' ] ; then
        _echo_p "Install SSH key $sshkey at tgen1 (root@$HOST_TGEN1_IP)"
        ssh-copy-id -i "$sshkey" "root@$HOST_TGEN1_IP"
    fi

    for i in {1..3} ; do
        if [ "$where" = 'all' ] || [ "$where" = "opicluster-master-$i" ] ; then
            _echo_p "Install SSH key $sshkey at opicluster-master-$i (core@192.168.122.$(( i + 1 )))"
            ssh-copy-id -i "$sshkey" -o "ProxyCommand ssh root@$HOST_TGEN1_IP nc %h %p" "core@192.168.122.$(( i + 1 ))" || rc=1
        fi
    done

    if [ "$where" = 'all' ] || [ "$where" = 'dh4' ] ; then
        _echo_p "Install SSH key $sshkey at dh4 (core@$HOST_DH4_IP)"
        ssh-copy-id -i "$sshkey" -o "ProxyCommand ssh root@$HOST_TGEN1_IP nc %h %p" "core@$HOST_DH4_IP" || rc=1
    fi

    if [ "$where" = 'all' ] || [ "$where" = 'dh4-acc' ] ; then
        _echo_p "Install SSH key $sshkey at dh4-acc (root@172.16.3.16)"
        ssh-copy-id -i "$sshkey" -o "ProxyCommand ssh root@$HOST_TGEN1_IP nc %h %p" "root@172.16.3.16" || rc=1
    fi

    return "$rc"
}

# 3:extra helper:4
#
# Show kubeadmin password for https://console-openshift-console.apps.opicluster.opicluster.opiproject-lab.org
do_kubeadmin() {
    _EXEC_NOTTY=1 \
    _EXEC_SILENT=1 \
    _exec tgen1 cat /root/opicluster/kubeadmin-passwd.txt
}

# 3:extra helper:5
# [DIR]
# Download kubeconfigs to current directory (or specify directory as argument).
do_kubeconfigs() {
    local dir="${1:-$ORIGINAL_PWD}"
    local pwd0="$PWD"

    cd "$ORIGINAL_PWD"

    dir="$(realpath -s "$dir")"

    [ -d "$dir" ] || die "Directory $dir does not exist"

    shift || true

    [ "$#" -eq 0 ] || die "Invalid arguments"

    do_kubeconfig_ocp > "$dir/kubeconfig.ocpcluster"
    do_kubeconfig_microshift > "$dir/kubeconfig.microshift"

    _echo "export KC_OCP=$(printf '%q' "$dir/kubeconfig.ocpcluster")"
    _echo "export KC_MSH=$(printf '%q' "$dir/kubeconfig.microshift")"

    cd "$pwd0"
}

# 3:extra helper:6
#
# Show the kubeconfig of the ocpcluster.
do_kubeconfig_ocp() {
    _EXEC_NOTTY=1 \
    _EXEC_SILENT=1 \
    _exec tgen1 cat /root/opicluster/kubeconfig
}

# 3:extra helper:7
#
# Show the kubeconfig for the IPU on dh4.
do_kubeconfig_microshift() {
    local hh

    if is_tgen1 ; then
        hh='172.16.3.16'
    else
        hh='172.22.1.104'
    fi

    _EXEC_NOTTY=1 \
    _EXEC_SILENT=1 \
    _exec dh4-acc cat /var/lib/microshift/resources/kubeadmin/kubeconfig \
        | sed -e 's/^\( *\)\(certificate-authority-data: .*\)$/\1#\2\n\1insecure-skip-tls-verify: true/' \
              -e 's#server: https://localhost:6443#server: https://'"$hh"':6443#'
}

# 3:extra helper:8
#
# Run shellcheck, black and mypy on the sources of the demo script.
do_check() {
    _echo_p "Validate script via shellcheck"
    shellcheck "$SCRIPTNAME"
    shellcheck "./inspect.sh"

    py_activate_venv

    _echo_p "Format python sources with black"
    if ! black --check . ; then
        black .
        die "Python sources were not correctly formatted with black."
    fi

    _echo_p "Check python sources with mypy"
    mypy

    deactivate

    _echo_p "Check usage command"
    usage 1>/dev/null
}

_main() {
    local cmd="$1"
    shift || true

    # For convenience, accept the following aliases:
    case "$cmd" in
        'predict_images'|predict_images/)
            cmd='predict'
            ;;
        'inspect.sh')
            cmd='inspect'
            ;;
    esac

    case "$cmd" in
        check | \
        info | \
        kubeadmin | \
        kubeconfig_microshift | \
        kubeconfig_ocp | \
        nginx_setup | \
        predict | \
        redeploy | \
        tls_generate_keys \
        )
            [ "$#" -eq 0 ] || die "Invalid arguments"
            "do_$cmd"
            ;;
        etc_hosts | \
        exec | \
        inspect | \
        kubeconfigs | \
        pod_detect_net1_ip | \
        pods_setup | \
        reboot | \
        remote | \
        ssh_copy_id | \
        tcpdump_dh4 | \
        tcpdump_pod \
        )
            "do_$cmd" "$@"
            ;;
        oc_msh | \
        oc_ocp \
        )
            "$cmd" "$@"
            ;;
        "-h"|"--help")
            usage
            ;;
        "")
            usage
            die "Missing command"
            ;;
        *)
            usage
            die "Unknown command '$cmd'"
            ;;
    esac
}

###############################################################################

ORIGINAL_PWD="$PWD"

if [[ "${BASH_SOURCE[0]}" != "$0" ]]; then
    # Sourced. We are done
    SCRIPTNAME0="${BASH_SOURCE[0]}"
    SCRIPTNAME="${BASH_SOURCE[0]}"
    die() { _echo "$C_RED$*$C_RESET" ; }
    return 0
fi

SCRIPTNAME0="$0"
SCRIPTNAME="$0"

cd "$(dirname "$0")"
SCRIPTNAME="$(basename "$0")"

set -eo pipefail
shopt -s inherit_errexit

_main "$@"

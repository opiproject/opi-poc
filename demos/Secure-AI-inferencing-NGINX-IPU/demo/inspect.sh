#!/usr/bin/env bash
#
# This file contains commands to run the solution end to end on diff tmux
# sessions, while also displaying comments relevant to each session.

cd "$(dirname "$0")" || exit 1

sleep_time="${OPI_DEMO_SLEEP_TIME:-5}"

###############################################################################
# === Configuration ===

session_name="ai_inference_server_offload_demo"

pane0_commands=(
  "./demo.sh exec dh4-acc cat /etc/os-release"
  "./demo.sh exec dh4-acc ifconfig enp0s1f0d1"
  "./demo.sh exec dh4-imc cat /etc/issue"
  "./demo.sh oc_ipu get nodes -o wide ; sleep 3; ./demo.sh oc_ipu get pods -A"
  "./demo.sh oc_ocp get nodes -o wide ; sleep 3; ./demo.sh oc_ocp get pods -n default"
  "export pod=\$(./demo.sh oc_ipu get pods -A -o name | grep 'vsp-p4' | head -n1); echo \"podname: \$pod\""
  "./demo.sh oc_ipu exec -it \$pod -n openshift-dpu-operator -- /opt/p4/p4-cp-nws/bin/ovs-vsctl show"
  "./demo.sh oc_ipu exec -it \$pod -n openshift-dpu-operator -- /opt/p4/p4-cp-nws/bin/p4rt-ctl dump-entries br0 | tail -n 25"
)
pane0_comments=(
  "Redhat OS running on Intel IPU Compuete Complex (ACC)"
  "Intel IPU ACC Primary Network IP"
  "IPU SDK Version on Intel IPU Management Compute Complex (IMC)"
  "Redhat Microshift Cluster Information - Running on Intel IPU Compute Complex"
  "Redhat Openshift Cluster Information - Running on Intel X86 Host Servers"
  "Get name of vsp-p4 pod"
  "Intel IPU's VSP-OVS Bridge Configuration"
  "Intel IPU's Packet Processing Engine Flows - P4 Programming rules"
)

pane1_commands=(
  #"./demo.sh nginx_setup"
  "./demo.sh oc_ipu exec -n openshift-dpu-operator -it nginx -- ifconfig"
  "./demo.sh oc_ipu exec -n openshift-dpu-operator -it nginx -- bash -c 'grep -R ^ /etc/nginx/*.conf /etc/nginx/conf.d/*.conf'"
  #"nowait:./demo.sh oc_ipu exec -n openshift-dpu-operator -it nginx -- tcpdump -ni net2"
  "nowait:./demo.sh oc_ipu logs -n openshift-dpu-operator --follow nginx"
)
pane1_comments=(
  #"Bring In NGINX Server Configuration Running at Microshift Cluster on Intel IPU"
  "Service Function Chaining(SFC): NGINX Attached With 2 Virtual Function Interfaces for SFC Packet Forwarding"
  "NGINX Configured as a Proxy Server to the Backend AI Inference Workload Servers (OpenVino) Running on Intel X86 Processor Node"
  #"Run TCPDUMP to Visualize the Packet Flows at NGINX Pod "
  "Show NGINX Access Log"
)

pane2_commands=(
  #"./demo.sh pods_setup 1"
  "./demo.sh oc_ocp exec -n default -it resnet50-model-server-1 -- ifconfig"
  "nowait:./demo.sh oc_ocp exec -n default -it resnet50-model-server-1 -- tcpdump -ni net1"
)
pane2_comments=(
  #"Bring In AI Workload Server resnet50-model-server-1 Configuration Running at Openshift Cluster on Intel X86 Host"
  "AI workload Pod Configured with 1 SRIOV-VF Interface to Receive Data Traffic from IPU"
  "Enable TCPDUMP to Visualize the Packet Flows at AI Workload Pod (OpenVino)"
)

pane3_commands=(
  #"./demo.sh pods_setup 2"
  "./demo.sh oc_ocp exec -n default -it resnet50-model-server-2 -- ifconfig"
  "nowait:./demo.sh oc_ocp exec -n default -it resnet50-model-server-2 -- tcpdump -ni net1"
)
pane3_comments=(
  #"Bring In AI Workload Server resnet50-model-server-2 Configuration Running at Openshift Cluster on Intel X86 Host"
  "AI workload Pod Configured with 1 SRIOV-VF Interface to Receive Data Traffic from IPU"
  "Enable TCPDUMP to Visualize the Packet Flows at AI Workload Pod (OpenVino)"
)

pane4_commands=(
  #"./demo.sh pods_setup 3"
  "./demo.sh oc_ocp exec -n default -it resnet50-model-server-3 -- ifconfig"
  "nowait:./demo.sh oc_ocp exec -n default -it resnet50-model-server-3 -- tcpdump -ni net1"
)
pane4_comments=(
  #"Bring In AI Workload Server resnet50-model-server-3 Configuration Running at Openshift Cluster on Intel X86 Host"
  "AI workload Pod Configured with 1 SRIOV-VF Interface to Receive Data Traffic from IPU"
  "Enable TCPDUMP to Visualize the Packet Flows at AI Workload Pod (OpenVino)"
)

pane0_client_commands=(
  "nowait:./demo.sh predict"
)
pane0_client_comments=(
  "AI Inference (Image Prediction) Demonstration with Proxy Server Offload At Intel IPU (SmartNIC)"
)

pane_nginx_status_commands=(
   #"nowait:watch -n1 ./demo.sh remote exec tgen1 curl -s --cacert /tmp/opilab-demo/server.crt https://172.16.3.200/nginx_status"
   "nowait:bash -c 'source ./demo.sh ; _cmd_ipu_statistics 1'"
)
pane_nginx_status_comments=(
  #"Offloaded NGINX Proxy Server Realtime Stats - Running On Intel IPU (SmartNIC)"
  "Interface network statistics for VFs inside NGINX pod"
)

export OPI_DEMO_EXEC_SILENT=1

_tmux_command_idx=0

_tmux_command() {
    local window="$1"
    local comment="$2"
    local cmd="$3"

    tmux send-keys -t "$session_name:$window" "echo -e \"\e[30;43m###### $comment ######\e[0m\"" C-m

    if [[ "$cmd" != 'nowait:'* ]] ; then
        _tmux_command_idx="$((_tmux_command_idx + 1))"
        tmux send-keys -t "$session_name:$window" "$cmd ; tmux wait-for -S \"$window:$_tmux_command_idx:done\"" C-m
        tmux wait-for "$window:$_tmux_command_idx:done"
    else
        cmd="${cmd#nowait:}"
        tmux send-keys -t "$session_name:$window" "$cmd" C-m
        sleep 1
    fi
    sleep "$sleep_time"
}

_tmux_caption() {
    local session="$1"
    shift
    local name="$*"

    tmux set -t "$session" pane-border-style fg=cyan,bg=black
    tmux set -t "$session" pane-border-status top
    tmux set -t "$session" pane-border-format "#{pane_index} #{pane_title}"
    tmux send-keys -t "$session" "printf '\\033]2;$name\\007'" C-m
}

_title_screen() {
    local msg1="Demonstration of an AI inference use case on Red Hat OpenShift and MicroShift cluster environments running on Intel x86 and IPU (Smart NIC) platforms."
    local msg2="Spawning TMUX Sessions To Configure and Execute This Demonstration."

    clear
    sleep 3

    local cols
    local lines

    # grab your terminal size
    cols="$(tput cols)"
    lines="$(tput lines)"

    # compute vertical padding (half the remaining lines)
    local vpad="$(( (lines - 1) / 2 ))"

    # compute horizontal padding (half the remaining columns)
    local hpad="$(( (cols - "${#msg1}") / 2 ))"
    (( hpad < 0 )) && hpad=0

    # print blank lines for vertical centering
    for ((i=0; i<vpad; i++)); do
      echo
    done

    # print spaces + colored message
    # \e[30;43m = black text on yellow background; \e[0m = reset
    printf '%*s\e[30;43m%s\e[0m\n' "$hpad" '' "$msg1"
    printf '%*s\e[30;43m%s\e[0m\n' "$hpad" '' "$msg2"
    sleep 5
}

_setup_tmux() {
    # --- Setup tmux session with one window split into three panes ---
    if tmux has-session -t "$session_name" 2>/dev/null; then
      tmux kill-session -t "$session_name"
    fi

    # Create session with single window
    tmux new-session -d -s "$session_name" -n main
    # Split main window: first vertical, then horizontal on the right pane
    tmux split-window -h -t "$session_name:0"
    tmux split-window -v -t "$session_name:0.1"
    tmux split-window -v -t "$session_name:0.2"
    tmux split-window -v -t "$session_name:0.3"
    tmux split-window -v -t "$session_name:0.4"

    _tmux_caption "$session_name:0.0" main
    _tmux_caption "$session_name:0.1" nginx pod
    _tmux_caption "$session_name:0.2" AI pod 1
    _tmux_caption "$session_name:0.3" AI pod 2
    _tmux_caption "$session_name:0.4" AI pod 3
}

_run_in_background() {
    # Pane 0 (left) – zoom before running
    tmux select-pane -t "$session_name:0.0"
    tmux resize-pane -Z -t "$session_name:0.0"

    for i in "${!pane0_commands[@]}" ; do
        _tmux_command "0.0" "${pane0_comments[i]}" "${pane0_commands[i]}"
    done
    sleep 3
    tmux resize-pane -Z -t "$session_name:0.0"

    # Pane 1 (top-right) – zoom before running
    tmux select-pane -t "$session_name:0.1"
    tmux resize-pane -Z -t "$session_name:0.1"
    for i in "${!pane1_commands[@]}"; do
        _tmux_command "0.1" "${pane1_comments[i]}" "${pane1_commands[i]}"
    done
    sleep 3
    tmux resize-pane -Z -t "$session_name:0.1"

    # Pane 2 (bottom-right)
    tmux select-pane -t "$session_name:0.2"
    tmux resize-pane -Z -t "$session_name:0.2"
    for i in "${!pane2_commands[@]}"; do
        _tmux_command "0.2" "${pane2_comments[i]}" "${pane2_commands[i]}"
    done
    sleep 3
    tmux resize-pane -Z -t "$session_name:0.2"

    # Pane 3 (bottom-right)
    tmux select-pane -t "$session_name:0.3"
    tmux resize-pane -Z -t "$session_name:0.3"
    for i in "${!pane3_commands[@]}"; do
        _tmux_command "0.3" "${pane3_comments[i]}" "${pane3_commands[i]}"
    done
    sleep 3
    tmux resize-pane -Z -t "$session_name:0.3"

    # Pane 4 (bottom-right)
    tmux select-pane -t "$session_name:0.4"
    tmux resize-pane -Z -t "$session_name:0.4"
    for i in "${!pane4_commands[@]}"; do
        _tmux_command "0.4" "${pane4_comments[i]}" "${pane4_commands[i]}"
    done
    sleep 3
    tmux resize-pane -Z -t "$session_name:0.4"

    tmux split-window -v -t "$session_name:0.0" -p 25
    _tmux_caption "$session_name:0.1" ipu stats
    tmux select-pane -t "$session_name:0.1"
    for i in "${!pane_nginx_status_commands[@]}"; do
        _tmux_command "main.1" "${pane_nginx_status_comments[i]}" "${pane_nginx_status_commands[i]}"
    done
    sleep 5

    tmux select-pane -t "$session_name:0.0"
    for i in "${!pane0_client_commands[@]}"; do
        _tmux_command "main.0" "${pane0_client_comments[i]}" "${pane0_client_commands[i]}"
    done
}

###############################################################################

_title_screen
_setup_tmux
_run_in_background &

# Attach session to view the layout
tmux attach-session -t "$session_name"

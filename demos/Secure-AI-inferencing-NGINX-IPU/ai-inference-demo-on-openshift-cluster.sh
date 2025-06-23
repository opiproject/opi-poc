#!/usr/bin/env bash
#This file contains commands to run the solution end to end on diff tmux sessions, 
#while also displaying comments relevant to each session.


# === Configuration ===
session_name="ai_inference_server_offload_demo"            # tmux session name
pause_seconds=5                 # seconds between commands

# Commands & comments per pane
pane0_commands=(
  "ssh-keygen -R 172.16.3.16; sshpass -p 'redhat' ssh -o StrictHostKeyChecking=no root@172.16.3.16"
  "cat /etc/os-release"
  "ifconfig enp0s1f0d1"
  "exit"
  "ssh-keygen -R 172.16.3.16; sshpass -p '' ssh -o StrictHostKeyChecking=no root@10.26.16.111"
  "cat /etc/issue"
  "exit"
  "export KUBECONFIG=/root/kubeconfig.microshift; oc get nodes -L kubernetes.io/hostname=worker-238 -o wide ; sleep 3; oc get pods -A"
  "export KUBECONFIG=/root/kubeconfig.ocpcluster; oc get nodes -L kubernetes.io/hostname=worker-238 -o wide ; sleep 3; oc get pods"
  "export KUBECONFIG=/root/kubeconfig.microshift"
  "pod=\$(oc get pods -A -o name | grep 'vsp-p4' | head -n1) ; oc exec -it \${pod} -n openshift-dpu-operator -- /opt/p4/p4-cp-nws/bin/ovs-vsctl show"
  "pod=\$(oc get pods -A -o name | grep 'vsp-p4' | head -n1) ; oc exec -it \${pod} -n openshift-dpu-operator -- /opt/p4/p4-cp-nws/bin/p4rt-ctl dump-entries br0"
)
pane0_comments=(
  "###### Connect to Intel IPU Compute Complex (ACC) ######"
  "###### Redhat OS running on Intel IPU Compuete Complex (ACC) ######"
  "###### Intel IPU ACC Primary Network IP ######"
  "###### Logout from ACC ######"
  "###### Connect to Intel IPU Management Compute Complex (IMC) ######"
  "###### IPU SDK Version ######"
  "###### Logout from IMC ######"
  "###### Redhat Microshift Cluster Information - Running on Intel IPU Compute Complex ######"
  "###### Redhat Openshift Cluster Information - Running on Intel X86 Host Servers ######"
  "###### Intel IPU Specific Redhat Microshift Cluster Key Components Required to Offload Network Workload ######"
  ""
  "Intel IPU's VSP-OVS Bridge Configuration"
  "Intel IPU's Packet Processing Engine Flows - P4 Programming rules"
)

pane1_commands=(
  "/root/summit/f5-redhat-on-intel-ipu/update_nginx_upstream_server_info.sh"
  "export KUBECONFIG=/root/kubeconfig.microshift"
  "oc exec -it nginx -n openshift-dpu-operator -- ifconfig"
  "oc exec -it nginx -n openshift-dpu-operator -- cat /etc/nginx/nginx.conf"
  "oc exec -it nginx -n openshift-dpu-operator -- tcpdump -ni net2"
)
pane1_comments=(
  "###### Bring In NGINX Server Configuration Running at Microshift Cluster on Intel IPU ######"
  ""
  "###### Service Function Chaining(SFC): NGINX Attached With 2 Virtual Function Interfaces for SFC Packet Forwarding ######"
  "###### NGINX Configured as a Proxy Server to the Backend AI Inference Workload Servers (OpenVino) Running on Intel X86 Processor Node ######"
  "###### Enable TCPDUMP to Visualize the Packet Flows at NGINX Pod  ######"
)

pane2_commands=(
  "export KUBECONFIG=/root/kubeconfig.ocpcluster"
  "oc exec -it resnet50-model-server-1 -- apt-get update"
  "oc exec -it resnet50-model-server-1 -- apt-get install -y net-tools"
  "oc exec -it resnet50-model-server-1 -- apt-get install -y tcpdump"
  "oc exec -it resnet50-model-server-1 -- ifconfig"
  "oc exec -it resnet50-model-server-1 -- tcpdump -ni net1"
)
pane2_comments=(
  "###### Bring In AI Workload Server resnet50-model-server-1 Configuration Running at Openshift Cluster on Intel X86 Host ######"
  ""
  ""
  ""
  "###### AI workload Pod Configured with 1 SRIOV-VF Interface to Receive Data Traffic from IPU ######"
  "###### Enable TCPDUMP to Visualize the Packet Flows at AI Workload Pod (OpenVino) ######"
)

pane3_commands=(
  "export KUBECONFIG=/root/kubeconfig.ocpcluster"
  "oc exec -it resnet50-model-server-2 -- apt-get update"
  "oc exec -it resnet50-model-server-2 -- apt-get install -y net-tools"
  "oc exec -it resnet50-model-server-2 -- apt-get install -y tcpdump"
  "oc exec -it resnet50-model-server-2 -- ifconfig"
  "oc exec -it resnet50-model-server-2 -- tcpdump -ni net1"
)
pane3_comments=(
  "###### Bring In AI Workload Server resnet50-model-server-2 Configuration Running at Openshift Cluster on Intel X86 Host ######"
  ""
  ""
  ""
  "###### AI workload Pod Configured with 1 SRIOV-VF Interface to Receive Data Traffic from IPU ######"
  "###### Enable TCPDUMP to Visualize the Packet Flows at AI Workload Pod (OpenVino) ######"
)

pane4_commands=(
  "export KUBECONFIG=/root/kubeconfig.ocpcluster"
  "oc exec -it resnet50-model-server-3 -- apt-get update"
  "oc exec -it resnet50-model-server-3 -- apt-get install -y net-tools"
  "oc exec -it resnet50-model-server-3 -- apt-get install -y tcpdump"
  "oc exec -it resnet50-model-server-3 -- ifconfig"
  "oc exec -it resnet50-model-server-3 -- tcpdump -ni net1"
)
pane4_comments=(
  "###### Bring In AI Workload Server resnet50-model-server-3 Configuration Running at Openshift Cluster on Intel X86 Host ######"
  ""
  ""
  ""
  "###### AI workload Pod Configured with 1 SRIOV-VF Interface to Receive Data Traffic from IPU ######"
  "###### Enable TCPDUMP to Visualize the Packet Flows at AI Workload Pod (OpenVino) ######"
)

# Commands & comments per pane
pane0_client_commands=(
  "/root/summit/f5-redhat-on-intel-ipu/run_image_prediction_random.py"
)

pane0_client_comments=(
  "###### AI Inference (Image Prediction) Demonstration with Proxy Server Offload At Intel IPU (SmartNIC) ######"
)

pane_nginx_status_commands=(
   "/root/summit/f5-redhat-on-intel-ipu/get_nginx_status.sh"
)

pane_nginx_status_comments=(
  "###### Offloaded NGINX Proxy Server Realtime Stats - Running On Intel IPU (SmartNIC) ######"
)

clear
msg="Demonstration of an AI inference use case on Red Hat OpenShift and MicroShift cluster environments running on Intel x86 and IPU (Smart NIC) platforms.\n
Spawning TMUX Sessions To Configure and Execute This Demonstration\n"
sleep 3

# grab your terminal size
cols=$(tput cols)
lines=$(tput lines)

# compute vertical padding (half the remaining lines)
vpad=$(( (lines - 1) / 2 ))

# compute horizontal padding (half the remaining columns)
hpad=$(( (cols - ${#msg}) / 2 ))
(( hpad < 0 )) && hpad=0

# print blank lines for vertical centering
for ((i=0; i<vpad; i++)); do
  echo
done

# print spaces + colored message
# \e[30;43m = black text on yellow background; \e[0m = reset
printf '%*s\e[30;43m%s\e[0m\n' "$hpad" '' "$msg"
sleep 5

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

# --- Execute commands in their designated panes ---
(
  # Pane 0 (left) – zoom before running
  tmux select-pane -t "${session_name}:0.0"
  tmux resize-pane -Z -t "${session_name}:0.0"

  for i in "${!pane0_commands[@]}"; do
    # print comment in pane 0
    tmux send-keys -t "${session_name}:0.0" "echo -e \"\e[30;43m ${pane0_comments[i]} \e[0m\"" C-m
    # run command in pane 0
    cmd="${pane0_commands[i]}"
    tmux send-keys -t "${session_name}:0.0" "$cmd" C-m
    sleep "$pause_seconds"
  done

  # Zoom out pane 0
  sleep 3
  tmux resize-pane -Z -t "${session_name}:0.0"

  # Pane 1 (top-right) – zoom before running
  tmux select-pane -t "${session_name}:0.1"
  tmux resize-pane -Z -t "${session_name}:0.1"

  for i in "${!pane1_commands[@]}"; do
    tmux send-keys -t "${session_name}:0.1" "echo -e \"\e[30;43m ${pane1_comments[i]} \e[0m\"" C-m
    cmd="${pane1_commands[i]}"
    tmux send-keys -t "${session_name}:0.1" "$cmd" C-m
    sleep "$pause_seconds"
  done
  sleep 3
  tmux resize-pane -Z -t "${session_name}:0.1"

  # Pane 2 (bottom-right)
  tmux select-pane -t "${session_name}:0.2"
  tmux resize-pane -Z -t "${session_name}:0.2"

  for i in "${!pane2_commands[@]}"; do
    tmux send-keys -t "${session_name}:0.2" "echo -e \"\e[30;43m ${pane2_comments[i]} \e[0m\"" C-m
    cmd="${pane2_commands[i]}"
    tmux send-keys -t "${session_name}:0.2" "$cmd" C-m
    sleep "$pause_seconds"
  done
  sleep 3
  tmux resize-pane -Z -t "${session_name}:0.2"

  # Pane 3 (bottom-right)
  tmux select-pane -t "${session_name}:0.3"
  tmux resize-pane -Z -t "${session_name}:0.3"

  for i in "${!pane2_commands[@]}"; do
    tmux send-keys -t "${session_name}:0.3" "echo -e \"\e[30;43m ${pane2_comments[i]} \e[0m\"" C-m
    cmd="${pane2_commands[i]}"
    tmux send-keys -t "${session_name}:0.3" "$cmd" C-m
    sleep "$pause_seconds"
  done
  sleep 3
  tmux resize-pane -Z -t "${session_name}:0.3"

  # Pane 4 (bottom-right)
  tmux select-pane -t "${session_name}:0.4"
  tmux resize-pane -Z -t "${session_name}:0.4"

  for i in "${!pane2_commands[@]}"; do
    tmux send-keys -t "${session_name}:0.4" "echo -e \"\e[30;43m ${pane2_comments[i]} \e[0m\"" C-m
    cmd="${pane2_commands[i]}"
    tmux send-keys -t "${session_name}:0.4" "$cmd" C-m
    sleep "$pause_seconds"
  done
  sleep 3
  tmux resize-pane -Z -t "${session_name}:0.4"

  for run in {0..2}; do
    for i in "${!pane0_client_commands[@]}"; do
      # print comment in pane 0
      tmux send-keys -t "${session_name}:main.${run}" "echo -e \"\e[30;43m ${pane0_client_comments[i]} \e[0m\"" C-m
      # run command in pane 0
      cmd="${pane0_client_commands[i]}"
      tmux send-keys -t "${session_name}:main.${run}" "$cmd" C-m
    done
    sleep 3
    tmux split-window -v -t "$session_name:main.${run}" 
  done

  sleep 5
  tmux select-pane -t "${session_name}:main.3"
  tmux resize-pane -Z -t "${session_name}:main.3"
  for i in "${!pane_nginx_status_commands[@]}"; do
    # print comment in pane 0
    tmux send-keys -t "${session_name}:main.3" "echo -e \"\e[30;43m ${pane_nginx_status_comments[i]} \e[0m\"" C-m
    # run command in pane 0
    cmd="${pane_nginx_status_commands[i]}"
    tmux send-keys -t "${session_name}:main.3" "$cmd" C-m
  done
  sleep 10
  tmux select-pane -t "${session_name}:main.3"
  tmux resize-pane -Z -t "${session_name}:main.3"

) &

# Attach session to view the layout
tmux attach-session -t "$session_name"

# How to Add an SSH Key to a GitHub Runner

This guide explains the manual steps required to add and configure an SSH key on a self-hosted GitHub Runner so it can securely access other servers or resources during CI/CD workflows.

## Prerequisites
* Access to the server running the GitHub Runner.
* An SSH Key pair (Public and Private keys). If you don't have one, generate it using:
  ```bash
  ssh-keygen -t ed25519 -C "your_email@example.com"

Step 1: Log into the GitHub Runner Server

1.Open your terminal or command prompt.
2.Connect to the machine where your GitHub self-hosted runner is installed:
ssh user@runner-ip-address

Step 2: Switch to the GitHub Runner User
It is highly recommended to store the SSH key under the specific user account that runs the GitHub Runner service (usually named runner or github).

1.Switch to the runner user:
sudo su - runner

Step 3: Add the Private SSH Key
1.Navigate to the .ssh directory of the runner user. If it doesn't exist, create it:
mkdir -p ~/.ssh
chmod 700 ~/.ssh

2.Create or open the id_ed25519 (or id_rsa) file to paste your private key:
nano ~/.ssh/id_ed25519

3.Paste your Private SSH Key into this file, then save and close it (Ctrl+O, Enter, Ctrl+X).

4.Set the correct permissions so the system allows it to be used:
chmod 600 ~/.ssh/id_ed25519

Step 4: Add the Public Key to the Target Server
To allow the runner to connect to a target destination server:

1.Copy the public key content (cat ~/.ssh/id_ed25519.pub).

2.Paste it into the target server's ~/.ssh/authorized_keys file.

Step 5: Test the Connection Manually
Before running a GitHub Actions workflow, test the connection once manually from the runner user to verify that it doesn't ask for a password prompt:
ssh -T target_user@target_server_ip

If it asks to trust the host, type yes and press Enter. This will add the server to known_hosts.

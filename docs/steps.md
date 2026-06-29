# Manual Steps to Bring Up Mt. Evans with Nginx

This document outlines the manual steps required to set up and bring up the Mt. Evans (Intel DPU) environment with Rocky Linux and Nginx.

## Prerequisites
* Intel Mt. Evans DPU Card installed in the server.
* Rocky Linux OS installed and running on the host/DPU.
* Root or sudo access to the system.

---

## Step 1: System and Interface Verification
Before installing Nginx, ensure that the Mt. Evans network interfaces are up and properly detected by the OS.

1. Check the available network interfaces:
   ```bash
   ip a

2.Verify that the Mt. Evans specific interfaces are initialized and have correct IP configurations matching the lab testbed topology

Step 2: Install Nginx on Rocky Linux
Since the environment uses Rocky Linux, we will use the standard package manager (dnf) to install Nginx.

1.Update the repository configuration:
sudo dnf update -y

2.Install Nginx server:
sudo dnf install nginx -y

Step 3: Configure Nginx for Mt. Evans Environment
Configure Nginx to handle traffic through the Mt. Evans DPU interfaces.

1.Open the default Nginx configuration file:
sudo nano /etc/nginx/nginx.conf

2.Adjust the listen directive to bind to the specific Mt. Evans interface IP address or keep it default (listen 80;) depending on the routing setup.

3.Save and close the file.

Step 4: Start and Enable Nginx Service
Enable the service so it starts automatically upon system reboot.

1.Start the Nginx service:
sudo systemctl start nginx

2.Enable it to run on startup:
sudo systemctl enable nginx

3.Verify the status of the service:
sudo systemctl status nginx

Step 5: Firewall Configuration (Optional)
If the firewall is running, allow HTTP traffic through the Mt. Evans interfaces.
sudo firewall-cmd --permanent --add-service=http
sudo firewall-cmd --reload

# DPU Host 3

HP RL300

## BIOS settings

- change power to always on

## Link speed

an aditional service was created in order to set the nic at boot in 1G mode
otherwise it starts in 10G/25G mode and link will be down

- run manually:

```bash
ethtool -s enP2p1s0f0np0 speed 1000
```

- or via systemctl

```bash
systemctl edit  --force --full link-speed.service
systemctl enable link-speed.service
systemctl start link-speed.service
```

where `link-speed.service` is:

```ini
[Unit]
After=network.target

[Service]
ExecStart=ethtool -s enP2p1s0f0np0 speed 1000

[Install]
WantedBy=default.target
```

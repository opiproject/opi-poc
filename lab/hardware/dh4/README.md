# DPU Host 4

Dell PowerEdge R760

## BIOS settings

- change power to always on

## IPU

Intel Dayton Peak

- to enable port forwarding for `OTEL` from `ACC` through `IMC`, run manually:

```bash
dnf install -y socat
socat tcp-l:4317,fork,reuseaddr tcp:172.22.0.1:4317
```

- or via systemctl

```bash
systemctl edit  --force --full socat-otel.service
systemctl enable socat-otel.service
systemctl start socat-otel.service
```

where `socat-otel.service` is:

```ini
[Unit]
After=network.target

[Service]
ExecStart=socat tcp-l:4317,fork,reuseaddr tcp:172.22.0.1:4317

[Install]
WantedBy=default.target
```

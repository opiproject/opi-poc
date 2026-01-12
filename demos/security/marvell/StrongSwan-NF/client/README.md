# StrongSwan Network Function - Client

This directory contains the client setup and execution instructions for testing the StrongSwan IPSec network function on Marvell DPUs.

## Client Setup and Execution

### Prerequisites

Clone the OPI godpu repository:

```bash
git clone https://github.com/opiproject/godpu.git
cd godpu
```

### Checkout Specific Commit

Checkout to the required commit:

```bash
git checkout 2bc04190a2db77b03ab2ea5f850e6ad93be4c385
```

This commit includes the fix: "change call for error messages with no format required" by Mark Sanders (committed on Mon Jan 20 10:25:57 2025 -0600).

### Apply Configuration Patch

Apply the IPSec tunnel and selector configuration patch:

```bash
git am /path/to/0001-ipsec-update-tunnel-and-selectors-to-match-the-setup.patch
```

This patch updates the default IPSec tunnel configuration:

- Local address: `20.20.20.1`
- Remote address: `20.20.20.20`
- Local traffic selector: `10.56.217.0/24`
- Remote traffic selector: `20.20.20.0/24`
- Skips cleanup of SAs for testing purposes

### Run the Client

Execute the IPSec test client:

```bash
go run main.go --addr 20.20.20.1:50151 ipsec test
```

To customize the traffic selectors for different network configurations:

```bash
go run main.go --addr 20.20.20.1:50151 ipsec test \
  --local-ts "10.56.217.0/24" \
  --remote-ts "20.20.20.0/24"
```

**Parameters:**

- `--addr 20.20.20.1:50151` - gRPC server address on the DPU
- `ipsec test` - Run the IPSec test command
- `--local-ts` (optional) - Local traffic selector subnet (default: from patch)
- `--remote-ts` (optional) - Remote traffic selector subnet (default: from patch)

### Expected Output

On successful execution, you should see:

```text
Rekeyed IKE_SA opi-test: success:"yes" matches:1
GRPC connection closed successfully
```

This indicates:

- IPSec tunnel is successfully established
- IKE_SA rekeying completed successfully
- gRPC connection to the DPU is functioning correctly

## Troubleshooting

If you encounter issues:

1. **Connection refused**: Verify the DPU is reachable at `20.20.20.1:50151`
2. **Patch fails**: Ensure you're on the correct commit before applying the patch
3. **Test fails**: Check the DPU logs and network connectivity between endpoints

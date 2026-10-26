# IPsec batch API benchmark

This measures the cost of the API, not of an IPsec implementation. It compares
creating and deleting IPsec SAs one RPC at a time with the `BatchCreateIpsecSas`
and `BatchDeleteIpsecSas` RPCs proposed in
[opiproject/opi-api#496](https://github.com/opiproject/opi-api/pull/496), which
fixes [opiproject/opi-api#452](https://github.com/opiproject/opi-api/issues/452)
(loading 1M tunnels took 45 minutes one at a time).

The server in `server.go` is an in-memory SAD behind a real gRPC server. The
client talks to it over loopback TCP. There is no datapath, so the numbers show
per-RPC overhead and batching gain only. Add `-rtt` to model a remote client.

## Run

All commands run from this folder.

Comparison of unary and batch modes (default 100000 SAs per mode):

```bash
go test -run TestUnaryVersusBatch -v -sas 100000 -results baseline.md ./...
```

Batch RPCs only, at scale:

```bash
go test -run TestBatchAtScale -v -timeout 60m \
  -scale 1000000,10000000,32000000 -batch 1000 -workers 8 -results scale.md ./...
```

Benchmarks:

```bash
go test -run '^$' -bench . -benchmem ./...
```

Flags:

| Flag       | Default | Meaning                                            |
| ---------- | ------- | -------------------------------------------------- |
| `-sas`     | 100000  | SAs created and deleted per mode                   |
| `-rtt`     | 0       | extra round-trip latency added to every RPC        |
| `-results` | empty   | write the results table (markdown) to this file    |
| `-scale`   | empty   | SA counts for `TestBatchAtScale`, comma-separated  |
| `-batch`   | 1000    | SAs per batch RPC in `TestBatchAtScale`            |
| `-workers` | 1       | batch RPCs in flight at once in `TestBatchAtScale` |

## Results

Measured on windows/amd64 with 24 CPUs and go1.26.0, gRPC over loopback TCP.

### Unary versus batch, 100000 SAs

The 1M and 32M columns are estimates extrapolated from the 100000 SA run, not
measurements.

| Mode                 |   RPCs |  Create | Create SAs/s | Speedup | Delete | Delete SAs/s | Speedup | 1M create (est.) | 32M create (est.) |
| -------------------- | -----: | ------: | -----------: | ------: | -----: | -----------: | ------: | ---------------: | ----------------: |
| unary, sequential    | 100000 | 15.013s |         6661 |    1.0x | 11.22s |         8912 |    1.0x |            2m30s |           1h20m4s |
| unary, 64 concurrent | 100000 |   583ms |       171514 |   25.8x |  514ms |       194532 |   21.8x |               6s |              3m7s |
| batch of 100         |   1000 |   528ms |       189290 |   28.4x |  122ms |       819972 |   92.0x |               5s |             2m49s |
| batch of 1000        |    100 |   461ms |       217105 |   32.6x |   25ms |      3951304 |  443.3x |               5s |             2m27s |
| batch of 10000       |     10 |   425ms |       235192 |   35.3x |   14ms |      7315503 |  820.8x |               4s |             2m16s |

### Batch at scale, 1000 SAs per batch, 1 RPC in flight

|      SAs |  RPCs |    Create | Create SAs/s |  Delete | Delete SAs/s | Heap with all SAs loaded |
| -------: | ----: | --------: | -----------: | ------: | -----------: | -----------------------: |
|  1000000 |  1000 |    5.194s |       192541 |   280ms |      3569728 |                 0.21 GiB |
| 10000000 | 10000 |   54.735s |       182697 |  3.065s |      3262282 |                 1.89 GiB |
| 32000000 | 32000 | 2m53.713s |       184212 | 10.837s |      2952735 |                 6.59 GiB |

### Batch at scale, 1000 SAs per batch, 8 RPCs in flight

|      SAs |  RPCs |  Create | Create SAs/s | Delete | Delete SAs/s | Heap with all SAs loaded |
| -------: | ----: | ------: | -----------: | -----: | -----------: | -----------------------: |
|  1000000 |  1000 |   869ms |      1151136 |   81ms |     12271294 |                 0.23 GiB |
| 10000000 | 10000 |  9.655s |      1035760 | 1.044s |      9581393 |                 1.90 GiB |
| 32000000 | 32000 | 33.444s |       956824 | 4.237s |      7552075 |                 6.60 GiB |

A second run of the same configuration:

|      SAs |  RPCs |  Create | Create SAs/s | Delete | Delete SAs/s | Heap with all SAs loaded |
| -------: | ----: | ------: | -----------: | -----: | -----------: | -----------------------: |
|  1000000 |  1000 |   916ms |      1091836 |   91ms |     11037260 |                 0.21 GiB |
| 10000000 | 10000 |  9.521s |      1050362 | 1.034s |      9671113 |                 1.90 GiB |
| 32000000 | 32000 | 33.521s |       954630 |  4.39s |      7288999 |                 6.60 GiB |

### Batch at scale, 10000 SAs per batch, 8 RPCs in flight

|      SAs | RPCs |  Create | Create SAs/s | Delete | Delete SAs/s | Heap with all SAs loaded |
| -------: | ---: | ------: | -----------: | -----: | -----------: | -----------------------: |
| 32000000 | 3200 | 31.786s |      1006720 | 6.989s |      4578937 |                 6.60 GiB |

## What this shows

- One RPC per SA runs at about 6.7k SAs/s. Extrapolated, 32M SAs take about
  1h20m.
- With batch 1000 and 8 RPCs in flight, 32M SAs are created in 33 s, about 1M
  SAs/s, and deleted in about 4 s. The two runs agree within 1%.
- Batch size matters less than concurrency at this scale: 10000 per batch
  gives 31.8 s against 33.4 s for 1000 per batch.
- With one batch in flight, throughput is about 185k SAs/s regardless of
  scale, so the gain from batching alone is roughly 28x on create.
- Server heap grows linearly, about 6.6 GiB for 32M SAs.

These are API-only numbers on loopback. A real implementation will be slower
per SA, and a remote client adds round trips, which favors batching further.

## Note on the replace lines

`go.mod` has `replace` lines that point the opi-api security and opinetcommon
modules at the `mgheorghe/opi-api` fork, because the batch RPCs are not in
upstream until opiproject/opi-api#496 merges. Drop the `replace` block and
bump the opi-api requirements once it is merged.

// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Keysight Technologies

package ipsecbatch

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	nc "github.com/opiproject/opi-api/network/opinetcommon/v1alpha1/gen/go"
	pb "github.com/opiproject/opi-api/v1/gen/go/security"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/proto"
)

var (
	numSas  = flag.Int("sas", 100000, "number of IPsec SAs to create and delete per mode")
	rtt     = flag.Duration("rtt", 0, "extra round-trip latency added to every RPC, to model a remote client")
	results = flag.String("results", "", "write the results table (markdown) to this file")

	scale     = flag.String("scale", "", "comma-separated SA counts for TestBatchAtScale, e.g. 1000000,32000000")
	batchSize = flag.Int("batch", 1000, "SAs per batch RPC in TestBatchAtScale")
	workers   = flag.Int("workers", 1, "batch RPCs in flight at once in TestBatchAtScale")
)

const concurrency = 64

// startServer serves an in-memory IpsecSaService and returns a connected client.
func startServer(tb testing.TB) (*Server, pb.IpsecSaServiceClient) {
	tb.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		tb.Fatal(err)
	}
	srv := NewServer()
	gs := grpc.NewServer()
	pb.RegisterIpsecSaServiceServer(gs, srv)
	go func() { _ = gs.Serve(lis) }()
	tb.Cleanup(gs.Stop)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if *rtt > 0 {
		opts = append(opts, grpc.WithUnaryInterceptor(
			func(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
				time.Sleep(*rtt)
				return invoker(ctx, method, req, reply, cc, opts...)
			}))
	}
	conn, err := grpc.NewClient(lis.Addr().String(), opts...)
	if err != nil {
		tb.Fatal(err)
	}
	tb.Cleanup(func() { _ = conn.Close() })
	return srv, pb.NewIpsecSaServiceClient(conn)
}

func ipv4(a uint32) *nc.IPAddress {
	return &nc.IPAddress{Af: nc.IpAf_IP_AF_INET, V4OrV6: &nc.IPAddress_V4Addr{V4Addr: a}}
}

// saRequest builds a tunnel-mode ESP SA with unique endpoints.
func saRequest(i int) *pb.CreateIpsecSaRequest {
	local := 0x0a000000 + uint32(i)  // 10.x.x.x
	remote := 0x14000000 + uint32(i) // 20.x.x.x
	return &pb.CreateIpsecSaRequest{
		IpsecSaId: fmt.Sprintf("sa-%d", i),
		IpsecSa: &pb.IpsecSa{
			Reqid: int64(i + 1),
			Config: &pb.IpsecSaConfig{
				Spi:      uint32(0x1000 + i),
				Protocol: pb.IpsecProtocol_IPSEC_PROTOCOL_ESP,
				Mode:     pb.IpsecMode_IPSEC_MODE_TUNNEL_MODE,
				EspAlgorithms: &pb.EspAlgorithms{
					Encryption: []pb.EncAlgorithm{pb.EncAlgorithm_ENC_ALGORITHM_AES256GCM128},
				},
				Tunnel: &pb.Tunnel{Local: ipv4(0x01000000 + uint32(i)), Remote: ipv4(0x02000000 + uint32(i))},
				TrafficSelector: &pb.TrafficSelector{
					LocalPrefix:  &nc.IPPrefix{Addr: ipv4(local), Len: 32},
					RemotePrefix: &nc.IPPrefix{Addr: ipv4(remote), Len: 32},
				},
			},
		},
	}
}

func saRequests(n int) []*pb.CreateIpsecSaRequest {
	reqs := make([]*pb.CreateIpsecSaRequest, n)
	for i := range reqs {
		reqs[i] = saRequest(i)
	}
	return reqs
}

// mode is one way of loading and unloading the SAs.
type mode struct {
	name   string
	rpcs   func(n int) int
	create func(ctx context.Context, c pb.IpsecSaServiceClient, reqs []*pb.CreateIpsecSaRequest) error
	delete func(ctx context.Context, c pb.IpsecSaServiceClient, names []string) error
}

func unaryMode() mode {
	return mode{
		name: "unary, sequential",
		rpcs: func(n int) int { return n },
		create: func(ctx context.Context, c pb.IpsecSaServiceClient, reqs []*pb.CreateIpsecSaRequest) error {
			for _, r := range reqs {
				if _, err := c.CreateIpsecSa(ctx, r); err != nil {
					return err
				}
			}
			return nil
		},
		delete: func(ctx context.Context, c pb.IpsecSaServiceClient, names []string) error {
			for _, n := range names {
				if _, err := c.DeleteIpsecSa(ctx, &pb.DeleteIpsecSaRequest{Name: n}); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

// parallel runs fn(i) for i in [0, n) on concurrency workers.
func parallel(n int, fn func(i int) error) error {
	return parallelN(concurrency, n, fn)
}

// parallelN runs fn(i) for i in [0, n) on w workers.
func parallelN(w, n int, fn func(i int) error) error {
	var next atomic.Int64
	var wg sync.WaitGroup
	errs := make(chan error, w)
	for ; w > 0; w-- {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := int(next.Add(1) - 1)
				if i >= n {
					return
				}
				if err := fn(i); err != nil {
					errs <- err
					return
				}
			}
		}()
	}
	wg.Wait()
	close(errs)
	return <-errs
}

func concurrentUnaryMode() mode {
	return mode{
		name: fmt.Sprintf("unary, %d concurrent", concurrency),
		rpcs: func(n int) int { return n },
		create: func(ctx context.Context, c pb.IpsecSaServiceClient, reqs []*pb.CreateIpsecSaRequest) error {
			return parallel(len(reqs), func(i int) error {
				_, err := c.CreateIpsecSa(ctx, reqs[i])
				return err
			})
		},
		delete: func(ctx context.Context, c pb.IpsecSaServiceClient, names []string) error {
			return parallel(len(names), func(i int) error {
				_, err := c.DeleteIpsecSa(ctx, &pb.DeleteIpsecSaRequest{Name: names[i]})
				return err
			})
		},
	}
}

func batchMode(size int) mode {
	return mode{
		name: fmt.Sprintf("batch of %d", size),
		rpcs: func(n int) int { return (n + size - 1) / size },
		create: func(ctx context.Context, c pb.IpsecSaServiceClient, reqs []*pb.CreateIpsecSaRequest) error {
			for start := 0; start < len(reqs); start += size {
				end := min(start+size, len(reqs))
				resp, err := c.BatchCreateIpsecSas(ctx, &pb.BatchCreateIpsecSasRequest{Requests: reqs[start:end]})
				if err != nil {
					return err
				}
				if len(resp.GetIpsecSas()) != end-start {
					return fmt.Errorf("batch returned %d SAs, want %d", len(resp.GetIpsecSas()), end-start)
				}
			}
			return nil
		},
		delete: func(ctx context.Context, c pb.IpsecSaServiceClient, names []string) error {
			for start := 0; start < len(names); start += size {
				end := min(start+size, len(names))
				if _, err := c.BatchDeleteIpsecSas(ctx, &pb.BatchDeleteIpsecSasRequest{Names: names[start:end]}); err != nil {
					return err
				}
			}
			return nil
		},
	}
}

func rate(n int, d time.Duration) string {
	return fmt.Sprintf("%.0f", float64(n)/d.Seconds())
}

// estimate scales a measured duration for n SAs to total SAs.
func estimate(d time.Duration, n, total int) string {
	return (time.Duration(float64(d) * float64(total) / float64(n))).Round(time.Second).String()
}

// TestUnaryVersusBatch times every load mode and checks the SAD content.
func TestUnaryVersusBatch(t *testing.T) {
	n := *numSas
	reqs := saRequests(n)
	names := make([]string, n)
	for i, r := range reqs {
		names[i] = saPrefix + r.GetIpsecSaId()
	}
	t.Logf("%d SAs, %d bytes each on the wire; a batch of 10000 is %d bytes (gRPC default limit 4 MiB)",
		n, proto.Size(reqs[0]), proto.Size(&pb.BatchCreateIpsecSasRequest{Requests: reqs[:min(n, 10000)]}))

	modes := []mode{unaryMode(), concurrentUnaryMode(), batchMode(100), batchMode(1000), batchMode(10000)}
	ctx := context.Background()

	var table strings.Builder
	fmt.Fprintf(&table, "%d IPsec SAs, gRPC over loopback TCP, extra RTT %v, %s/%s, %d CPUs, %s\n\n",
		n, *rtt, runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.Version())
	fmt.Fprintln(&table, "| Mode | RPCs | Create | Create SAs/s | Speedup | Delete | Delete SAs/s | Speedup | 1M create (est.) | 32M create (est.) |")
	fmt.Fprintln(&table, "| ---- | ---: | -----: | -----------: | ------: | -----: | -----------: | ------: | ---------------: | ----------------: |")

	var baseCreate, baseDelete time.Duration
	var batch1000 time.Duration
	for _, m := range modes {
		srv, client := startServer(t)

		start := time.Now()
		if err := m.create(ctx, client, reqs); err != nil {
			t.Fatalf("%s: create: %v", m.name, err)
		}
		create := time.Since(start)
		if got := srv.Len(); got != n {
			t.Fatalf("%s: %d SAs after create, want %d", m.name, got, n)
		}

		start = time.Now()
		if err := m.delete(ctx, client, names); err != nil {
			t.Fatalf("%s: delete: %v", m.name, err)
		}
		del := time.Since(start)
		if got := srv.Len(); got != 0 {
			t.Fatalf("%s: %d SAs after delete, want 0", m.name, got)
		}

		if baseCreate == 0 {
			baseCreate, baseDelete = create, del
		}
		if m.name == batchMode(1000).name {
			batch1000 = create
		}
		fmt.Fprintf(&table, "| %s | %d | %v | %s | %.1fx | %v | %s | %.1fx | %s | %s |\n",
			m.name, m.rpcs(n),
			create.Round(time.Millisecond), rate(n, create), float64(baseCreate)/float64(create),
			del.Round(time.Millisecond), rate(n, del), float64(baseDelete)/float64(del),
			estimate(create, n, 1_000_000), estimate(create, n, 32_000_000))
	}

	t.Log("\n" + table.String())
	if *results != "" {
		if err := os.WriteFile(*results, []byte(table.String()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if batch1000 >= baseCreate {
		t.Errorf("batch of 1000 took %v, not faster than unary %v", batch1000, baseCreate)
	}
}

// TestBatchAtScale loads and unloads millions of SAs with batch RPCs only.
func TestBatchAtScale(t *testing.T) {
	if *scale == "" {
		t.Skip("set -scale to run, e.g. -scale=1000000,32000000")
	}
	size, w := *batchSize, *workers
	ctx := context.Background()

	var table strings.Builder
	fmt.Fprintf(&table, "Batch RPCs only, %d SAs per batch, %d in flight, gRPC over loopback TCP, extra RTT %v, %s/%s, %d CPUs, %s\n\n",
		size, w, *rtt, runtime.GOOS, runtime.GOARCH, runtime.NumCPU(), runtime.Version())
	fmt.Fprintln(&table, "| SAs | RPCs | Create | Create SAs/s | Delete | Delete SAs/s | Heap with all SAs loaded |")
	fmt.Fprintln(&table, "| --: | ---: | -----: | -----------: | -----: | -----------: | -----------------------: |")

	for _, field := range strings.Split(*scale, ",") {
		var n int
		if _, err := fmt.Sscan(field, &n); err != nil || n <= 0 {
			t.Fatalf("bad -scale entry %q", field)
		}
		batches := (n + size - 1) / size
		srv, client := startServer(t)

		start := time.Now()
		err := parallelN(w, batches, func(b int) error {
			first, end := b*size, min((b+1)*size, n)
			reqs := make([]*pb.CreateIpsecSaRequest, end-first)
			for j := range reqs {
				reqs[j] = saRequest(first + j)
			}
			resp, err := client.BatchCreateIpsecSas(ctx, &pb.BatchCreateIpsecSasRequest{Requests: reqs})
			if err == nil && len(resp.GetIpsecSas()) != len(reqs) {
				err = fmt.Errorf("batch %d returned %d SAs, want %d", b, len(resp.GetIpsecSas()), len(reqs))
			}
			return err
		})
		create := time.Since(start)
		if err != nil {
			t.Fatalf("%d SAs: create: %v", n, err)
		}
		if got := srv.Len(); got != n {
			t.Fatalf("%d SAs: %d after create", n, got)
		}
		runtime.GC()
		var ms runtime.MemStats
		runtime.ReadMemStats(&ms)
		t.Logf("%d SAs created in %v", n, create.Round(time.Millisecond))

		start = time.Now()
		err = parallelN(w, batches, func(b int) error {
			first, end := b*size, min((b+1)*size, n)
			names := make([]string, end-first)
			for j := range names {
				names[j] = fmt.Sprintf("%ssa-%d", saPrefix, first+j)
			}
			_, err := client.BatchDeleteIpsecSas(ctx, &pb.BatchDeleteIpsecSasRequest{Names: names})
			return err
		})
		del := time.Since(start)
		if err != nil {
			t.Fatalf("%d SAs: delete: %v", n, err)
		}
		if got := srv.Len(); got != 0 {
			t.Fatalf("%d SAs: %d left after delete", n, got)
		}
		t.Logf("%d SAs deleted in %v", n, del.Round(time.Millisecond))

		fmt.Fprintf(&table, "| %d | %d | %v | %s | %v | %s | %.2f GiB |\n",
			n, batches, create.Round(time.Millisecond), rate(n, create),
			del.Round(time.Millisecond), rate(n, del), float64(ms.HeapAlloc)/(1<<30))
		srv.Reset()
		runtime.GC()
	}

	t.Log("\n" + table.String())
	if *results != "" {
		if err := os.WriteFile(*results, []byte(table.String()), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func BenchmarkCreateIpsecSa(b *testing.B) {
	srv, client := startServer(b)
	ctx := context.Background()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := client.CreateIpsecSa(ctx, saRequest(i)); err != nil {
			b.Fatal(err)
		}
	}
	b.StopTimer()
	b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N), "ns/SA")
	srv.Reset()
}

func BenchmarkBatchCreateIpsecSas(b *testing.B) {
	for _, size := range []int{100, 1000, 10000} {
		b.Run(fmt.Sprintf("batch=%d", size), func(b *testing.B) {
			srv, client := startServer(b)
			ctx := context.Background()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				if srv.Len() >= 1_000_000 {
					b.StopTimer()
					srv.Reset()
					b.StartTimer()
				}
				// Built inside the timed loop, like the unary benchmark.
				reqs := make([]*pb.CreateIpsecSaRequest, size)
				for j := range reqs {
					reqs[j] = saRequest(i*size + j)
				}
				if _, err := client.BatchCreateIpsecSas(ctx, &pb.BatchCreateIpsecSasRequest{Requests: reqs}); err != nil {
					b.Fatal(err)
				}
			}
			b.StopTimer()
			b.ReportMetric(float64(b.Elapsed().Nanoseconds())/float64(b.N*size), "ns/SA")
			srv.Reset()
		})
	}
}

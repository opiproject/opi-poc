// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"context"
	"log"

	pb "github.com/opiproject/opi-api/security/proto"
	"google.golang.org/grpc"
)

func do_ipsec(conn grpc.ClientConnInterface, ctx context.Context) {
	// IPsec
	c1 := pb.NewIPsecClient(conn)

	// Load IPsec connection
	local_ipsec := pb.IPsecLoadConnReq {
		Connection: &pb.Connection {
			Name: "opi-test",
			Version: "2",
			LocalAddrs: []*pb.Addrs {
				{
					Addr: "192.168.200.200",
				},
			},
			RemoteAddrs: []*pb.Addrs {
				{
					Addr: "192.168.200.210",
				},
			},
			Proposals: &pb.Proposals {
					CryptoAlg: []pb.CryptoAlgorithm { pb.CryptoAlgorithm_AES256GCM128 },
					IntegAlg: []pb.IntegAlgorithm { pb.IntegAlgorithm_SHA256_96 },
					Dhgroups: []pb.DiffieHellmanGroups { pb.DiffieHellmanGroups_CURVE25519 },
			},
			Children: []*pb.Child {
				{
					Name: "opi-child",
					EspProposals: &pb.Proposals {
							CryptoAlg: []pb.CryptoAlgorithm { pb.CryptoAlgorithm_AES256GMAC },
							IntegAlg: []pb.IntegAlgorithm { pb.IntegAlgorithm_SHA512 },
							Dhgroups: []pb.DiffieHellmanGroups { pb.DiffieHellmanGroups_CURVE25519 },
					},
				},
			},
		},
	}

	rs1, err := c1.IPsecLoadConn(ctx, &local_ipsec)
	if err != nil {
		log.Fatalf("could not load IPsec tunnel: %v", err)
	}
	log.Printf("Loaded: %v", rs1)

	// Unload

	unload_ipsec := pb.IPsecUnloadConnReq {
		Name: "opi-test",
	}

	rs2, err := c1.IPsecUnloadConn(ctx, &unload_ipsec)
	if err != nil {
		log.Fatalf("could not unload IPsec tunnel: %v", err)
	}
	log.Printf("Unloaded: %v", rs2)
}

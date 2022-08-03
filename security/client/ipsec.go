// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"context"
	"log"
	"math/rand"

	pb "github.com/opiproject/opi-api/security/proto"
	"google.golang.org/grpc"
)

func do_ipsec(conn grpc.ClientConnInterface, ctx context.Context) {
	// IPsec
	c1 := pb.NewIPsecClient(conn)

	// Create IPsec Connection
	local_ipsec := pb.IPsecCreateRequest{
		Tunnel: &pb.TunnelInterfaces{
			Tunnels: []*pb.TunnelInterfaces_Tunnel{
				{
					Name: "Test Tunnel",
					LocalIp: "192.168.200.200",
					RemoteIp: "192.168.200.210",
					LocalSpi: rand.Uint32(),
					CryptoAlg: pb.CryptoAlgorithm_AES_GCM_256,
					IntegAlg: pb.IntegAlgorithm_SHA_512_256,
					Mode: pb.IPsecMode_TUNNEL_MODE,
				},
			},
		},
	}

	rs1, err := c1.IPsecCreate(ctx, &local_ipsec)
	if err != nil {
		log.Fatalf("could not create IPsec tunnel: %v", err)
	}
	log.Printf("Added: %v", rs1)

	// Get the IPsec Connection
	rs2, err := c1.IPsecGet(ctx, &pb.IPsecGetRequest { Id: rs1.GetId() })
	if err != nil {
		log.Fatalf("could not get IPsec tunnel: %v", err)
	}
	log.Printf("Get: %v", rs2)

	// Update the IPsec Connection
	update_ipsec := pb.IPsecUpdateRequest {
		Id: rs1.GetId(),
		Tunnel: &pb.TunnelInterfaces {
			Tunnels: []*pb.TunnelInterfaces_Tunnel {
				{
						Name: "Renamed Test Tunnel",
						LocalIp: "192.168.200.220",
						RemoteIp: "192.168.200.250",
						LocalSpi: rand.Uint32(),
						CryptoAlg: pb.CryptoAlgorithm_AES_GMAC_256,
						IntegAlg: pb.IntegAlgorithm_SHA1_96,
						Mode: pb.IPsecMode_TRANSPORT_MODE,
				},
			},
		},
	}

	rs3, err := c1.IPsecUpdate(ctx, &update_ipsec)
	if err != nil {
		log.Fatalf("could not update IPsec tunnel: %v", err)
	}
	log.Printf("Updated: %v", rs3)

	// Get the IPsec Connection
	rs4, err := c1.IPsecGet(ctx, &pb.IPsecGetRequest { Id: rs1.GetId() })
	if err != nil {
		log.Fatalf("could not get IPsec tunnel: %v", err)
	}
	log.Printf("Get: %v", rs4)

	// Delete IPsec Connection
	rs5, err := c1.IPsecDelete(ctx, &pb.IPsecDeleteRequest { Id: rs1.GetId() })
	if err != nil {
		log.Fatalf("could not delete IPsec tunnel: %v", err)
	}
	log.Printf("Deleted: %v", rs5)
}

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
			Vips: &pb.Vips { Vip: []string { "0.0.0.0", }, },
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
			LocalAuth: &pb.LocalAuth { Auth: pb.AuthType_PSK, Id: "hacker@strongswan.org" },
			RemoteAuth: &pb.RemoteAuth { Auth: pb.AuthType_PSK, Id: "server.strongswan.org" },
			Children: []*pb.Child {
				{
					Name: "opi-child",
					EspProposals: &pb.Proposals {
						CryptoAlg: []pb.CryptoAlgorithm { pb.CryptoAlgorithm_AES256GCM128 },
						IntegAlg: []pb.IntegAlgorithm { pb.IntegAlgorithm_SHA512 },
						Dhgroups: []pb.DiffieHellmanGroups { pb.DiffieHellmanGroups_CURVE25519 },
					},
					RemoteTs: &pb.TrafficSelectors {
						Ts: []*pb.TrafficSelectors_TrafficSelector {
							{
								Cidr: "10.1.0.0/16",
							},
						},
					},
				},
			},
		},
	}

	ver_req := pb.IPsecVersionReq {}

	vresp, err := c1.IPsecVersion(ctx, &ver_req)
	if err != nil {
		log.Fatalf("could not get IPsec version")
	}
	log.Printf("Daemon  [%v]", vresp.GetDaemon())
	log.Printf("Version [%v]", vresp.GetVersion())
	log.Printf("Sysname [%v]", vresp.GetSysname())
	log.Printf("Release [%v]", vresp.GetRelease())
	log.Printf("Machine [%v]", vresp.GetMachine())

	rs1, err := c1.IPsecLoadConn(ctx, &local_ipsec)
	if err != nil {
		log.Fatalf("could not load IPsec tunnel: %v", err)
	}
	log.Printf("Loaded: %v", rs1)

	// Bring the connection up
	init_conn := pb.IPsecInitiateReq {
		Ike: "opi-test",
		Child: "opi-child",
	}

	init_ret, err := c1.IPsecInitiate(ctx, &init_conn)
	if err != nil {
		log.Fatalf("could not initiate IPsec tunnel: %v", err)
	}
	log.Printf("Initiated: %v", init_ret)

	// Terminate the connection
	term_conn := pb.IPsecTerminateReq{
		Ike: "opi-test",
		Child: "opi-child",
	}

	term_ret, err := c1.IPsecTerminate(ctx, &term_conn)
	if err != nil {
		log.Fatalf("could not terminate IPsec tunnel: %v", err)
	}
	log.Printf("Initiated: %v", term_ret)


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

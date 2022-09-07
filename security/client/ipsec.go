// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"context"
	"log"

	pb "github.com/opiproject/opi-api/security/proto"
	"google.golang.org/grpc"
	"github.com/go-ping/ping"
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

	stats_req := pb.IPsecStatsReq {}

	stats_resp, err := c1.IPsecStats(ctx, &stats_req)
	if err != nil {
		log.Fatalf("could not get IPsec stats")
	}
	log.Printf("IPsec stats\n%s", stats_resp.GetStatus())

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

	// List the IKE_SAs
	ike_sas := pb.IPsecListSasReq {
		Ike: "opi-test",
	}

	list_sas_ret, err := c1.IPsecListSas(ctx, &ike_sas)
	if err != nil {
		log.Fatalf("could not list IKE_SAs: %v", err)
	}
	log.Printf("Returned IKE_SAs: %v", list_sas_ret)

	// List the connections
	list_conn := pb.IPsecListConnsReq {
		Ike: "opi-test",
	}

	list_conns_ret, err := c1.IPsecListConns(ctx, &list_conn)
	if err != nil {
		log.Fatalf("could not list connections: %v", err)
	}
	log.Printf("Returned connections: %v", list_conns_ret)

	// List the certificats
	list_certs := pb.IPsecListCertsReq {
		Type: "any",
	}

	list_certs_ret, err := c1.IPsecListCerts(ctx, &list_certs)
	if err != nil {
		log.Fatalf("could not list certificates: %v", err)
	}
	log.Printf("Returned connections: %v", list_certs_ret)

	// Ping across the tunnel.
	// .NOTE: The container this test runs in is linked to the appropriate
	//        strongSwan container.
	pinger, err := ping.NewPinger(*pingaddr)
	if err != nil {
		log.Fatalf("Cannot create Pinger")
	}
	pinger.Count = 5
	// .NOTE: This blocks until it finishes
	err = pinger.Run()
	if err != nil {
		log.Fatalf("Ping command to host 10.3.0.1 failed")
	}
	stats := pinger.Statistics() // get send/receive/duplicate/rtt stats

	log.Printf("Ping stats: %v", stats)

	// Rekey the IKE_SA
	rekey_conn := pb.IPsecRekeyReq {
		Ike: "opi-test",
	}

	rekey_ret, err := c1.IPsecRekey(ctx, &rekey_conn)
	if err != nil {
		log.Fatalf("could not rekey IPsec tunnel: %v", err)
	}
	log.Printf("Rekeyed IKE_SA %s: %v", "opi-test", rekey_ret)

	// Terminate the connection
	term_conn := pb.IPsecTerminateReq{
		Ike: "opi-test",
	}

	term_ret, err := c1.IPsecTerminate(ctx, &term_conn)
	if err != nil {
		log.Fatalf("could not terminate IPsec tunnel: %v", err)
	}
	log.Printf("Terminate: %v", term_ret)


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

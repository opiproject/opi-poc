// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/opiproject/opi-api/storage/proto"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
)

func main() {
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// NVMeSubsystem
	c1 := pb.NewNVMeSubsystemServiceClient(conn)
	rs1, err := c1.NVMeSubsystemCreate(ctx, &pb.NVMeSubsystemCreateRequest{Subsystem: &pb.NVMeSubsystem{NQN: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", rs1)
	rs2, err := c1.NVMeSubsystemDelete(ctx, &pb.NVMeSubsystemDeleteRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", rs2)
	rs3, err := c1.NVMeSubsystemUpdate(ctx, &pb.NVMeSubsystemUpdateRequest{Subsystem: &pb.NVMeSubsystem{NQN: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not update NVMe subsystem: %v", err)
	}
	log.Printf("Updated: %v", rs3)
	rs4, err := c1.NVMeSubsystemList(ctx, &pb.NVMeSubsystemListRequest{})
	if err != nil {
		log.Fatalf("could not list NVMe subsystem: %v", err)
	}
	log.Printf("Listed: %v", rs4)
	rs5, err := c1.NVMeSubsystemGet(ctx, &pb.NVMeSubsystemGetRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", rs5.Subsystem.NQN)
	rs6, err := c1.NVMeSubsystemStats(ctx, &pb.NVMeSubsystemStatsRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not stats NVMe subsystem: %v", err)
	}
	log.Printf("Stats: %s", rs6.Stats)

	// NVMeController
	c2 := pb.NewNVMeControllerServiceClient(conn)
	rc1, err := c2.NVMeControllerCreate(ctx, &pb.NVMeControllerCreateRequest{Controller: &pb.NVMeController{Name: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", rc1)
	rc2, err := c2.NVMeControllerDelete(ctx, &pb.NVMeControllerDeleteRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", rc2)
	rc3, err := c2.NVMeControllerUpdate(ctx, &pb.NVMeControllerUpdateRequest{Controller: &pb.NVMeController{Name: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not update NVMe subsystem: %v", err)
	}
	log.Printf("Updated: %v", rc3)
	rc4, err := c2.NVMeControllerList(ctx, &pb.NVMeControllerListRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not list NVMe subsystem: %v", err)
	}
	log.Printf("Listed: %s", rc4)
	rc5, err := c2.NVMeControllerGet(ctx, &pb.NVMeControllerGetRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", rc5.Controller.Name)
	rc6, err := c2.NVMeControllerStats(ctx, &pb.NVMeControllerStatsRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not stats NVMe subsystem: %v", err)
	}
	log.Printf("Stats: %s", rc6.Stats)

	// NVMeNamespace
	c3 := pb.NewNVMeNamespaceServiceClient(conn)
	rn1, err := c3.NVMeNamespaceCreate(ctx, &pb.NVMeNamespaceCreateRequest{Namespace: &pb.NVMeNamespace{Name: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", rn1)
	rn2, err := c3.NVMeNamespaceDelete(ctx, &pb.NVMeNamespaceDeleteRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", rn2)
	rn3, err := c3.NVMeNamespaceUpdate(ctx, &pb.NVMeNamespaceUpdateRequest{Namespace: &pb.NVMeNamespace{Name: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not update NVMe subsystem: %v", err)
	}
	log.Printf("Updated: %v", rn3)
	rn4, err := c3.NVMeNamespaceList(ctx, &pb.NVMeNamespaceListRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not list NVMe subsystem: %v", err)
	}
	log.Printf("Listed: %v", rn4)
	rn5, err := c3.NVMeNamespaceGet(ctx, &pb.NVMeNamespaceGetRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %v", rn5.Namespace.Name)
	rn6, err := c3.NVMeNamespaceStats(ctx, &pb.NVMeNamespaceStatsRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not stats NVMe subsystem: %v", err)
	}
	log.Printf("Stats: %v", rn6.Stats)

	// NVMfRemoteController
	c4 := pb.NewNVMfRemoteControllerServiceClient(conn)
	rr0, err := c4.NVMfRemoteControllerConnect(ctx, &pb.NVMfRemoteControllerConnectRequest{Controller: &pb.NVMfRemoteController{Id: 8}})
	if err != nil {
		log.Fatalf("could not connect to Remote NVMf controller: %v", err)
	}
	log.Printf("Connected: %v", rr0)
	rr1, err := c4.NVMfRemoteControllerDisconnect(ctx, &pb.NVMfRemoteControllerDisconnectRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not disconnect from Remote NVMf controller: %v", err)
	}
	log.Printf("Disconnected: %v", rr1)
	rr2, err := c4.NVMfRemoteControllerReset(ctx, &pb.NVMfRemoteControllerResetRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not reset Remote NVMf controller: %v", err)
	}
	log.Printf("Reset: %v", rr2)
	rr3, err := c4.NVMfRemoteControllerList(ctx, &pb.NVMfRemoteControllerListRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not list Remote NVMf controllerd: %v", err)
	}
	log.Printf("List: %v", rr3)
	rr4, err := c4.NVMfRemoteControllerGet(ctx, &pb.NVMfRemoteControllerGetRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not get Remote NVMf controller: %v", err)
	}
	log.Printf("Got: %v", rr4)
	rr5, err := c4.NVMfRemoteControllerStats(ctx, &pb.NVMfRemoteControllerStatsRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not stats from Remote NVMf controller: %v", err)
	}
	log.Printf("Stats: %v", rr5)

}
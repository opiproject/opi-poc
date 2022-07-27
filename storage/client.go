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
	r1, err := c1.NVMeSubsystemCreate(ctx, &pb.NVMeSubsystemCreateRequest{Subsystem: &pb.NVMeSubsystem{NQN: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", r1)

	r2, err := c1.NVMeSubsystemDelete(ctx, &pb.NVMeSubsystemDeleteRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", r2)

	r3, err := c1.NVMeSubsystemGet(ctx, &pb.NVMeSubsystemGetRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", r3.Subsystem.NQN)

	// NVMeController
	c2 := pb.NewNVMeControllerServiceClient(conn)
	r4, err := c2.NVMeControllerCreate(ctx, &pb.NVMeControllerCreateRequest{Controller: &pb.NVMeController{Name: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", r4)

	r5, err := c2.NVMeControllerDelete(ctx, &pb.NVMeControllerDeleteRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", r5)

	r6, err := c2.NVMeControllerGet(ctx, &pb.NVMeControllerGetRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", r6.Controller.Name)

	// NVMeNamespace
	c3 := pb.NewNVMeNamespaceServiceClient(conn)
	r7, err := c3.NVMeNamespaceCreate(ctx, &pb.NVMeNamespaceCreateRequest{Namespace: &pb.NVMeNamespace{Name: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", r7)

	r8, err := c3.NVMeNamespaceDelete(ctx, &pb.NVMeNamespaceDeleteRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", r8)

	r9, err := c3.NVMeNamespaceGet(ctx, &pb.NVMeNamespaceGetRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %v", r9.Namespace.Name)

	// NVMfRemoteController
	c4 := pb.NewNVMfRemoteControllerServiceClient(conn)
	r0, err := c4.NVMfRemoteControllerDisconnect(ctx, &pb.NVMfRemoteControllerDisconnectRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not disconnect from Remote NVMf controller: %v", err)
	}
	log.Printf("Disconnected: %v", r0)

}
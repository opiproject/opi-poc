package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "opi.storage.v1/proto"
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
	c1 := pb.NewNVMeSubsystemClient(conn)
	r1, err := c1.NVMeSubsystemCreate(ctx, &pb.NVMeSubsystemCreateRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %s", r1.Name)

	r2, err := c1.NVMeSubsystemDelete(ctx, &pb.NVMeSubsystemDeleteRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %s", r2.Name)

	r3, err := c1.NVMeSubsystemGet(ctx, &pb.NVMeSubsystemGetRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", r3.Name)

	// NVMeController
	c2 := pb.NewNVMeControllerClient(conn)
	r4, err := c2.NVMeControllerCreate(ctx, &pb.NVMeControllerCreateRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %s", r4.Name)

	r5, err := c2.NVMeControllerDelete(ctx, &pb.NVMeControllerDeleteRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %s", r5.Name)

	r6, err := c2.NVMeControllerGet(ctx, &pb.NVMeControllerGetRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", r6.Name)

	// NVMeNamespace
	c3 := pb.NewNVMeNamespaceClient(conn)
	r7, err := c3.NVMeNamespaceCreate(ctx, &pb.NVMeNamespaceCreateRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %s", r7.Name)

	r8, err := c3.NVMeNamespaceDelete(ctx, &pb.NVMeNamespaceDeleteRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %s", r8.Name)

	r9, err := c3.NVMeNamespaceGet(ctx, &pb.NVMeNamespaceGetRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", r9.Name)

	// NVMfRemoteController
	c4 := pb.NewNVMfRemoteControllerClient(conn)
	r0, err := c4.NVMfRemoteControllerDisconnect(ctx, &pb.NVMfRemoteControllerDisconnectRequest{Name: "OPI-Nvme"})
	if err != nil {
		log.Fatalf("could not disconnect from Remote NVMf controller: %v", err)
	}
	log.Printf("Disconnected: %s", r0.Name)

}
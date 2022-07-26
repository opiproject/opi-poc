package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "opi.storage.v1"
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
	c := pb.nVMeSubsystemClient(conn)
	r, err := c.NVMeSubsystemCreate(ctx, &pb.NVMeSubsystemCreateRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %s", r.GetMessage())

	r, err := c.NVMeSubsystemDelete(ctx, &pb.NVMeSubsystemDeleteRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %s", r.GetMessage())

	r, err := c.NVMeSubsystemGet(ctx, &pb.NVMeSubsystemGetRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", r.GetMessage())

	// NVMeController
	c := pb.NVMeControllerClient(conn)
	r, err := c.NVMeControllerCreate(ctx, &pb.NVMeControllerCreateRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %s", r.GetMessage())

	r, err := c.NVMeControllerDelete(ctx, &pb.NVMeControllerDeleteRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %s", r.GetMessage())

	r, err := c.NVMeControllerGet(ctx, &pb.NVMeControllerGetRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", r.GetMessage())

	// NVMeNamespace
	c := pb.NVMeNamespaceClient(conn)
	r, err := c.NVMeNamespaceCreate(ctx, &pb.NVMeNamespaceCreateRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %s", r.GetMessage())

	r, err := c.NVMeNamespaceDelete(ctx, &pb.NVMeNamespaceDeleteRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %s", r.GetMessage())

	r, err := c.NVMeNamespaceGet(ctx, &pb.NVMeNamespaceGetRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", r.GetMessage())

	// NVMfRemoteController
	c := pb.NVMfRemoteControllerClient(conn)
	r, err := c.NVMfRemoteControllerDisconnect(ctx, &pb.NVMfRemoteControllerDisconnectRequest{Name: *name})
	if err != nil {
		log.Fatalf("could not disconnect from Remote NVMf controller: %v", err)
	}
	log.Printf("Disconnected: %s", r.GetMessage())

}
package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	pb "github.com/opiproject/opi-api/storage/proto"
)

func do_backend(conn grpc.ClientConnInterface, ctx context.Context) {
	// NVMfRemoteController
	c4 := pb.NewNVMfRemoteControllerServiceClient(conn)
	rr0, err := c4.NVMfRemoteControllerConnect(ctx, &pb.NVMfRemoteControllerConnectRequest{Controller: &pb.NVMfRemoteController{Id: 1}})
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
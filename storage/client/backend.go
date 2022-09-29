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
	rr0, err := c4.NVMfRemoteControllerConnect(ctx, &pb.NVMfRemoteControllerConnectRequest{Ctrl: &pb.NVMfRemoteController{Id: 8, Traddr:"127.0.0.1", Trsvcid:4444, Subnqn:"nqn.2016-06.io.spdk:cnode1"}})
	if err != nil {
		log.Fatalf("could not connect to Remote NVMf controller: %v", err)
	}
	log.Printf("Connected: %v", rr0)
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
	rr1, err := c4.NVMfRemoteControllerDisconnect(ctx, &pb.NVMfRemoteControllerDisconnectRequest{Id: 8})
	if err != nil {
		log.Fatalf("could not disconnect from Remote NVMf controller: %v", err)
	}
	log.Printf("Disconnected: %v", rr1)

	// NullDebug
	c1 := pb.NewNullDebugServiceClient(conn)
	log.Printf("Testing NewNullDebugServiceClient")
	rs1, err := c1.NullDebugCreate(ctx, &pb.NullDebugCreateRequest{Device: &pb.NullDebug{Name: "OpiNull9"}})
	if err != nil {
		log.Fatalf("could not create NULL device: %v", err)
	}
	log.Printf("Added: %v", rs1)
	rs3, err := c1.NullDebugUpdate(ctx, &pb.NullDebugUpdateRequest{Device: &pb.NullDebug{Name: "OpiNull9"}})
	if err != nil {
		log.Fatalf("could not update NULL device: %v", err)
	}
	log.Printf("Updated: %v", rs3)
	rs4, err := c1.NullDebugList(ctx, &pb.NullDebugListRequest{})
	if err != nil {
		log.Fatalf("could not list NULL device: %v", err)
	}
	log.Printf("Listed: %v", rs4)
	rs5, err := c1.NullDebugGet(ctx, &pb.NullDebugGetRequest{Id: 9})
	if err != nil {
		log.Fatalf("could not get NULL device: %v", err)
	}
	log.Printf("Got: %s", rs5.Device.Name)
	rs6, err := c1.NullDebugStats(ctx, &pb.NullDebugStatsRequest{Id: 9})
	if err != nil {
		log.Fatalf("could not stats NULL device: %v", err)
	}
	log.Printf("Stats: %s", rs6.Stats)
	rs2, err := c1.NullDebugDelete(ctx, &pb.NullDebugDeleteRequest{Id: 9})
	if err != nil {
		log.Fatalf("could not delete NULL device: %v", err)
	}
	log.Printf("Deleted: %v", rs2)

}
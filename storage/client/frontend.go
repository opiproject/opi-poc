package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	pb "github.com/opiproject/opi-api/storage/proto"
)

func do_frontend(conn grpc.ClientConnInterface, ctx context.Context) {

	// NVMeSubsystem
	c1 := pb.NewNVMeSubsystemServiceClient(conn)
	rs1, err := c1.NVMeSubsystemCreate(ctx, &pb.NVMeSubsystemCreateRequest{Subsystem: &pb.NVMeSubsystem{Nqn: "Malloc7"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", rs1)
	rs2, err := c1.NVMeSubsystemDelete(ctx, &pb.NVMeSubsystemDeleteRequest{Id: 7})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", rs2)
	rs3, err := c1.NVMeSubsystemUpdate(ctx, &pb.NVMeSubsystemUpdateRequest{Subsystem: &pb.NVMeSubsystem{Nqn: "Malloc1"}})
	if err != nil {
		log.Fatalf("could not update NVMe subsystem: %v", err)
	}
	log.Printf("Updated: %v", rs3)
	rs4, err := c1.NVMeSubsystemList(ctx, &pb.NVMeSubsystemListRequest{})
	if err != nil {
		log.Fatalf("could not list NVMe subsystem: %v", err)
	}
	log.Printf("Listed: %v", rs4)
	rs5, err := c1.NVMeSubsystemGet(ctx, &pb.NVMeSubsystemGetRequest{Id: 1})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", rs5.Subsystem.Nqn)
	rs6, err := c1.NVMeSubsystemStats(ctx, &pb.NVMeSubsystemStatsRequest{Id: 1})
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

	// VirtioBlk
	c4 := pb.NewVirtioBlkServiceClient(conn)
	rv1, err := c4.VirtioBlkCreate(ctx, &pb.VirtioBlkCreateRequest{Controller: &pb.VirtioBlk{Name: "VirtioBlk8", Bdev:"Malloc1"}})
	if err != nil {
		log.Fatalf("could not create VirtioBlk Controller: %v", err)
	}
	log.Printf("Added: %v", rv1)
	rv3, err := c4.VirtioBlkUpdate(ctx, &pb.VirtioBlkUpdateRequest{Controller: &pb.VirtioBlk{Name: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not update VirtioBlk Controller: %v", err)
	}
	log.Printf("Updated: %v", rv3)
	rv4, err := c4.VirtioBlkList(ctx, &pb.VirtioBlkListRequest{SubsystemId: 8})
	if err != nil {
		log.Fatalf("could not list VirtioBlk Controller: %v", err)
	}
	log.Printf("Listed: %v", rv4)
	rv5, err := c4.VirtioBlkGet(ctx, &pb.VirtioBlkGetRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not get VirtioBlk Controller: %v", err)
	}
	log.Printf("Got: %v", rv5.Controller.Name)
	rv6, err := c4.VirtioBlkStats(ctx, &pb.VirtioBlkStatsRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not stats VirtioBlk Controller: %v", err)
	}
	log.Printf("Stats: %v", rv6.Stats)
	rv2, err := c4.VirtioBlkDelete(ctx, &pb.VirtioBlkDeleteRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not delete VirtioBlk Controller: %v", err)
	}
	log.Printf("Deleted: %v", rv2)
}
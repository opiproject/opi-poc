package main

import (
	"context"
	"log"

	pb "github.com/opiproject/opi-api/storage/v1/gen/go"
	"google.golang.org/grpc"
)

func doFrontend(ctx context.Context, conn grpc.ClientConnInterface) {
	err := executeNVMeSubsystem(ctx, conn)
	if err != nil {
		log.Fatalf("Error executeNVMeSubsystem: %v", err)
		return
	}

	err = executeNVMeController(ctx, conn)
	if err != nil {
		log.Fatalf("Error executeNVMeController: %v", err)
		return
	}

	err = executeNVMeNamespace(ctx, conn)
	if err != nil {
		log.Fatalf("Error executeNVMeNamespace: %v", err)
		return
	}

	err = executeVirtioBlk(ctx, conn)
	if err != nil {
		log.Fatalf("Error executeVirtioBlk: %v", err)
		return
	}

	c5, err := executeVirtioScsiController(ctx, conn)
	if err != nil {
		log.Fatalf("Error executeVirtioScsiController: %v", err)
		return
	}

	err = executeVirtioScsiLun(ctx, conn, c5)
	if err != nil {
		log.Fatalf("Error executeVirtioScsiLun: %v", err)
		return
	}
}

func executeVirtioScsiLun(ctx context.Context, conn grpc.ClientConnInterface, c5 pb.VirtioScsiControllerServiceClient) error {
	// VirtioScsiLun
	c6 := pb.NewVirtioScsiLunServiceClient(conn)
	log.Printf("Testing NewVirtioScsiLunServiceClient")
	rl1, err := c6.VirtioScsiLunCreate(ctx, &pb.VirtioScsiLunCreateRequest{Lun: &pb.VirtioScsiLun{ControllerId: 8, Bdev: "Malloc1"}})
	if err != nil {
		log.Fatalf("could not create VirtioScsi subsystem: %v", err)
	}
	log.Printf("Added: %v", rl1)
	rl3, err := c6.VirtioScsiLunUpdate(ctx, &pb.VirtioScsiLunUpdateRequest{Lun: &pb.VirtioScsiLun{ControllerId: 8, Bdev: "Malloc1"}})
	if err != nil {
		log.Fatalf("could not update VirtioScsi subsystem: %v", err)
	}
	log.Printf("Updated: %v", rl3)
	rl4, err := c6.VirtioScsiLunList(ctx, &pb.VirtioScsiLunListRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not list VirtioScsi subsystem: %v", err)
	}
	log.Printf("Listed: %v", rl4)
	rl5, err := c6.VirtioScsiLunGet(ctx, &pb.VirtioScsiLunGetRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not get VirtioScsi subsystem: %v", err)
	}
	log.Printf("Got: %v", rl5.Lun.Bdev)
	rl6, err := c6.VirtioScsiLunStats(ctx, &pb.VirtioScsiLunStatsRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not stats VirtioScsi subsystem: %v", err)
	}
	log.Printf("Stats: %v", rl6.Stats)
	rl2, err := c6.VirtioScsiLunDelete(ctx, &pb.VirtioScsiLunDeleteRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not delete VirtioScsi subsystem: %v", err)
	}
	log.Printf("Deleted: %v", rl2)

	rss2, err := c5.VirtioScsiControllerDelete(ctx, &pb.VirtioScsiControllerDeleteRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not delete VirtioScsi subsystem: %v", err)
	}
	log.Printf("Deleted: %v", rss2)
	return err
}

func executeVirtioScsiController(ctx context.Context, conn grpc.ClientConnInterface) (pb.VirtioScsiControllerServiceClient, error) {
	// VirtioScsiController
	c5 := pb.NewVirtioScsiControllerServiceClient(conn)
	log.Printf("Testing NewVirtioScsiControllerServiceClient")
	rss1, err := c5.VirtioScsiControllerCreate(ctx, &pb.VirtioScsiControllerCreateRequest{Controller: &pb.VirtioScsiController{Name: "OPI-VirtioScsi8"}})
	if err != nil {
		log.Fatalf("could not create VirtioScsi subsystem: %v", err)
	}
	log.Printf("Added: %v", rss1)
	rss3, err := c5.VirtioScsiControllerUpdate(ctx, &pb.VirtioScsiControllerUpdateRequest{Controller: &pb.VirtioScsiController{Name: "OPI-VirtioScsi8"}})
	if err != nil {
		log.Fatalf("could not update VirtioScsi subsystem: %v", err)
	}
	log.Printf("Updated: %v", rss3)
	rss4, err := c5.VirtioScsiControllerList(ctx, &pb.VirtioScsiControllerListRequest{})
	if err != nil {
		log.Fatalf("could not list VirtioScsi subsystem: %v", err)
	}
	log.Printf("Listed: %s", rss4)
	rss5, err := c5.VirtioScsiControllerGet(ctx, &pb.VirtioScsiControllerGetRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not get VirtioScsi subsystem: %v", err)
	}
	log.Printf("Got: %s", rss5.Controller.Name)
	rss6, err := c5.VirtioScsiControllerStats(ctx, &pb.VirtioScsiControllerStatsRequest{ControllerId: 8})
	if err != nil {
		log.Fatalf("could not stats VirtioScsi subsystem: %v", err)
	}
	log.Printf("Stats: %s", rss6.Stats)
	return c5, err
}

func executeVirtioBlk(ctx context.Context, conn grpc.ClientConnInterface) error {
	// VirtioBlk
	c4 := pb.NewVirtioBlkServiceClient(conn)
	log.Printf("Testing NewVirtioBlkServiceClient")
	rv1, err := c4.VirtioBlkCreate(ctx, &pb.VirtioBlkCreateRequest{Controller: &pb.VirtioBlk{Name: "VirtioBlk8", Bdev: "Malloc1"}})
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

	return err
}

func executeNVMeNamespace(ctx context.Context, conn grpc.ClientConnInterface) error {
	// NVMeNamespace
	c3 := pb.NewNVMeNamespaceServiceClient(conn)
	log.Printf("Testing NewNVMeNamespaceServiceClient")
	rn1, err := c3.NVMeNamespaceCreate(ctx, &pb.NVMeNamespaceCreateRequest{Namespace: &pb.NVMeNamespace{SubsystemId: "nqn.2016-06.io.spdk:cnode1", Bdev: "Malloc1"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", rn1)
	rn2, err := c3.NVMeNamespaceDelete(ctx, &pb.NVMeNamespaceDeleteRequest{SubsystemId: 1, NamespaceId: 1})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", rn2)
	rn3, err := c3.NVMeNamespaceUpdate(ctx, &pb.NVMeNamespaceUpdateRequest{Namespace: &pb.NVMeNamespace{Name: "OPI-Nvme"}})
	if err != nil {
		log.Fatalf("could not update NVMe subsystem: %v", err)
	}
	log.Printf("Updated: %v", rn3)
	rn4, err := c3.NVMeNamespaceList(ctx, &pb.NVMeNamespaceListRequest{SubsystemId: 1})
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

	return err
}

func executeNVMeController(ctx context.Context, conn grpc.ClientConnInterface) error {
	// NVMeController
	c2 := pb.NewNVMeControllerServiceClient(conn)
	log.Printf("Testing NewNVMeControllerServiceClient")
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

	return err
}

func executeNVMeSubsystem(ctx context.Context, conn grpc.ClientConnInterface) error {
	// NVMeSubsystem
	c1 := pb.NewNVMeSubsystemServiceClient(conn)
	log.Printf("Testing NewNVMeSubsystemServiceClient")
	rs1, err := c1.NVMeSubsystemCreate(ctx, &pb.NVMeSubsystemCreateRequest{Subsystem: &pb.NVMeSubsystem{Nqn: "nqn.2022-09.io.spdk:opi7"}})
	if err != nil {
		log.Fatalf("could not create NVMe subsystem: %v", err)
	}
	log.Printf("Added: %v", rs1)
	rs3, err := c1.NVMeSubsystemUpdate(ctx, &pb.NVMeSubsystemUpdateRequest{Subsystem: &pb.NVMeSubsystem{Nqn: "nqn.2022-09.io.spdk:opi7"}})
	if err != nil {
		log.Fatalf("could not update NVMe subsystem: %v", err)
	}
	log.Printf("Updated: %v", rs3)
	rs4, err := c1.NVMeSubsystemList(ctx, &pb.NVMeSubsystemListRequest{})
	if err != nil {
		log.Fatalf("could not list NVMe subsystem: %v", err)
	}
	log.Printf("Listed: %v", rs4)
	rs5, err := c1.NVMeSubsystemGet(ctx, &pb.NVMeSubsystemGetRequest{Nqn: "7"})
	if err != nil {
		log.Fatalf("could not get NVMe subsystem: %v", err)
	}
	log.Printf("Got: %s", rs5.Subsystem.Nqn)
	rs6, err := c1.NVMeSubsystemStats(ctx, &pb.NVMeSubsystemStatsRequest{Nqn: "7"})
	if err != nil {
		log.Fatalf("could not stats NVMe subsystem: %v", err)
	}
	log.Printf("Stats: %s", rs6.Stats)
	rs2, err := c1.NVMeSubsystemDelete(ctx, &pb.NVMeSubsystemDeleteRequest{Nqn: "7"})
	if err != nil {
		log.Fatalf("could not delete NVMe subsystem: %v", err)
	}
	log.Printf("Deleted: %v", rs2)
	return err
}

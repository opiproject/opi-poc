// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/opiproject/opi-api/storage/proto"
	"github.com/google/uuid"
)

//////////////////////////////////////////////////////////

func (s *server) NVMeSubsystemCreate(ctx context.Context, in *pb.NVMeSubsystemCreateRequest) (*pb.NVMeSubsystemCreateResponse, error) {
	log.Printf("NVMeSubsystemCreate: Received from client: %v", in)
	params := struct {
		Name        string `json:"name"`
		BlockSize   int64  `json:"block_size"`
		NumBlocks   int64  `json:"num_blocks"`
		Uuid        string `json:"uuid"`
	}{
		Name:       in.GetSubsystem().GetNqn(),
		BlockSize:  512,
		NumBlocks:  64,
		Uuid:      uuid.New().String(),
	}
	var result string
	err := call("bdev_malloc_create", &params, &result)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	log.Printf("Received from SPDK: %v", result)
	return &pb.NVMeSubsystemCreateResponse{}, nil
}

func (s *server) NVMeSubsystemDelete(ctx context.Context, in *pb.NVMeSubsystemDeleteRequest) (*pb.NVMeSubsystemDeleteResponse, error) {
	log.Printf("NVMeSubsystemDelete: Received from client: %v", in)
	params := struct {
		Name        string `json:"name"`
	}{
		Name:       fmt.Sprint("Malloc", in.GetId()),
	}
	var result bool
	err := call("bdev_malloc_delete", &params, &result)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	log.Printf("Received from SPDK: %v", result)
	if (!result) {
		log.Printf("Could not delete: %v", in)
	}
	return &pb.NVMeSubsystemDeleteResponse{}, nil
}

func (s *server) NVMeSubsystemUpdate(ctx context.Context, in *pb.NVMeSubsystemUpdateRequest) (*pb.NVMeSubsystemUpdateResponse, error) {
	log.Printf("NVMeSubsystemUpdate: Received from client: %v", in)
	params1 := struct {
		Name        string `json:"name"`
	}{
		Name:       in.GetSubsystem().GetNqn(),
	}
	var result1 bool
	err1 := call("bdev_malloc_delete", &params1, &result1)
	if err1 != nil {
		log.Printf("error: %v\n", err1)
	}
	log.Printf("Received from SPDK: %v", result1)
	if (!result1) {
		log.Printf("Could not delete: %v", in)
	}
	params2 := struct {
		Name        string `json:"name"`
		BlockSize   int64  `json:"block_size"`
		NumBlocks   int64  `json:"num_blocks"`
		Uuid        string `json:"uuid"`
	}{
		Name:       in.GetSubsystem().GetNqn(),
		BlockSize:  512,
		NumBlocks:  64,
		Uuid:      uuid.New().String(),
	}
	var result2 string
	err2 := call("bdev_malloc_create", &params2, &result2)
	if err2 != nil {
		log.Printf("error: %v\n", err2)
	}
	log.Printf("Received from SPDK: %v", result2)
	return &pb.NVMeSubsystemUpdateResponse{}, nil
}

func (s *server) NVMeSubsystemList(ctx context.Context, in *pb.NVMeSubsystemListRequest) (*pb.NVMeSubsystemListResponse, error) {
	log.Printf("NVMeSubsystemList: Received from client: %v", in)
	var result []struct {
		Name        string `json:"name"`
		BlockSize   int64  `json:"block_size"`
		NumBlocks   int64  `json:"num_blocks"`
		Uuid        string `json:"uuid"`
	}
	err := call("bdev_get_bdevs", nil, &result)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	log.Printf("Received from SPDK: %v", result)
	Blobarray := make([]*pb.NVMeSubsystem, len(result))
	for i := range result {
		r := &result[i]
		Blobarray[i] = &pb.NVMeSubsystem{Nqn: r.Name}
	}
	return &pb.NVMeSubsystemListResponse{Subsystem: Blobarray}, nil
}

func (s *server) NVMeSubsystemGet(ctx context.Context, in *pb.NVMeSubsystemGetRequest) (*pb.NVMeSubsystemGetResponse, error) {
	log.Printf("NVMeSubsystemGet: Received from client: %v", in)
	params := struct {
		Name string `json:"name"`
	}{
		Name:       fmt.Sprint("Malloc", in.GetId()),
	}
	var result []struct {
		Name        string `json:"name"`
		BlockSize   int64  `json:"block_size"`
		NumBlocks   int64  `json:"num_blocks"`
		Uuid        string `json:"uuid"`
	}
	err := call("bdev_get_bdevs", &params, &result)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	log.Printf("Received from SPDK: %v", result)
	if (len(result) != 1) {
		log.Printf("expecting exactly 1 result")
	}
	return &pb.NVMeSubsystemGetResponse{Subsystem: &pb.NVMeSubsystem{Nqn: result[0].Name}}, nil
}

func (s *server) NVMeSubsystemStats(ctx context.Context, in *pb.NVMeSubsystemStatsRequest) (*pb.NVMeSubsystemStatsResponse, error) {
	log.Printf("NVMeSubsystemStats: Received from client: %v", in)
	params := struct {
		Name string `json:"name"`
	}{
		Name:     fmt.Sprint("Malloc", in.GetId()),
	}
	// See https://mholt.github.io/json-to-go/
	var result struct {
		TickRate int64 `json:"tick_rate"`
		Ticks    int64 `json:"ticks"`
		Bdevs    []struct {
			Name              string `json:"name"`
			BytesRead         int    `json:"bytes_read"`
			NumReadOps        int    `json:"num_read_ops"`
			BytesWritten      int    `json:"bytes_written"`
			NumWriteOps       int    `json:"num_write_ops"`
			BytesUnmapped     int    `json:"bytes_unmapped"`
			NumUnmapOps       int    `json:"num_unmap_ops"`
			ReadLatencyTicks  int    `json:"read_latency_ticks"`
			WriteLatencyTicks int    `json:"write_latency_ticks"`
			UnmapLatencyTicks int    `json:"unmap_latency_ticks"`
		} `json:"bdevs"`
	}
	err := call("bdev_get_iostat", &params, &result)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	log.Printf("Received from SPDK: %v", result)
	if (len(result.Bdevs) != 1) {
		log.Printf("expecting exactly 1 result")
	}
	return &pb.NVMeSubsystemStatsResponse{Id: 1, Stats: fmt.Sprint(result.Bdevs[0])}, nil
}

//////////////////////////////////////////////////////////

func (s *server) NVMeControllerCreate(ctx context.Context, in *pb.NVMeControllerCreateRequest) (*pb.NVMeControllerCreateResponse, error) {
	log.Printf("Received from client: %v", in.GetController())
	return &pb.NVMeControllerCreateResponse{}, nil
}

func (s *server) NVMeControllerDelete(ctx context.Context, in *pb.NVMeControllerDeleteRequest) (*pb.NVMeControllerDeleteResponse, error) {
	log.Printf("Received from client: %v", in.GetControllerId())
	return &pb.NVMeControllerDeleteResponse{}, nil
}

func (s *server) NVMeControllerUpdate(ctx context.Context, in *pb.NVMeControllerUpdateRequest) (*pb.NVMeControllerUpdateResponse, error) {
	log.Printf("Received from client: %v", in.GetController())
	return &pb.NVMeControllerUpdateResponse{}, nil
}

func (s *server) NVMeControllerList(ctx context.Context, in *pb.NVMeControllerListRequest) (*pb.NVMeControllerListResponse, error) {
	log.Printf("Received from client: %v", in.GetSubsystemId())
	Blobarray := make([]*pb.NVMeController, 3)
	return &pb.NVMeControllerListResponse{Controller: Blobarray}, nil
}

func (s *server) NVMeControllerGet(ctx context.Context, in *pb.NVMeControllerGetRequest) (*pb.NVMeControllerGetResponse, error) {
	log.Printf("Received from client: %v", in.GetControllerId())
	return &pb.NVMeControllerGetResponse{Controller: &pb.NVMeController{Name: "Hello " + fmt.Sprint(in.GetControllerId())}}, nil
}

func (s *server) NVMeControllerStats(ctx context.Context, in *pb.NVMeControllerStatsRequest) (*pb.NVMeControllerStatsResponse, error) {
	log.Printf("Received from client: %v", in.GetControllerId())
	return &pb.NVMeControllerStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

func (s *server) NVMeNamespaceCreate(ctx context.Context, in *pb.NVMeNamespaceCreateRequest) (*pb.NVMeNamespaceCreateResponse, error) {
	log.Printf("Received from client: %v", in.GetNamespace())
	return &pb.NVMeNamespaceCreateResponse{}, nil
}

func (s *server) NVMeNamespaceDelete(ctx context.Context, in *pb.NVMeNamespaceDeleteRequest) (*pb.NVMeNamespaceDeleteResponse, error) {
	log.Printf("Received from client: %v", in.GetNamespaceId())
	return &pb.NVMeNamespaceDeleteResponse{}, nil
}

func (s *server) NVMeNamespaceUpdate(ctx context.Context, in *pb.NVMeNamespaceUpdateRequest) (*pb.NVMeNamespaceUpdateResponse, error) {
	log.Printf("Received from client: %v", in.GetNamespace())
	return &pb.NVMeNamespaceUpdateResponse{}, nil
}

func (s *server) NVMeNamespaceList(ctx context.Context, in *pb.NVMeNamespaceListRequest) (*pb.NVMeNamespaceListResponse, error) {
	log.Printf("Received from client: %v", in.GetControllerId())
	Blobarray := make([]*pb.NVMeNamespace, 3)
	return &pb.NVMeNamespaceListResponse{Namespace: Blobarray}, nil
}

func (s *server) NVMeNamespaceGet(ctx context.Context, in *pb.NVMeNamespaceGetRequest) (*pb.NVMeNamespaceGetResponse, error) {
	log.Printf("Received from client: %v", in.GetNamespaceId())
	return &pb.NVMeNamespaceGetResponse{Namespace: &pb.NVMeNamespace{Name: "Hello " + fmt.Sprint(in.GetNamespaceId())}}, nil
}

func (s *server) NVMeNamespaceStats(ctx context.Context, in *pb.NVMeNamespaceStatsRequest) (*pb.NVMeNamespaceStatsResponse, error) {
	log.Printf("Received from client: %v", in.GetNamespaceId())
	return &pb.NVMeNamespaceStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

func (s *server) VirtioBlkCreate(ctx context.Context, in *pb.VirtioBlkCreateRequest) (*pb.VirtioBlkCreateResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioBlkCreateResponse{}, nil
}

func (s *server) VirtioBlkDelete(ctx context.Context, in *pb.VirtioBlkDeleteRequest) (*pb.VirtioBlkDeleteResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioBlkDeleteResponse{}, nil
}

func (s *server) VirtioBlkUpdate(ctx context.Context, in *pb.VirtioBlkUpdateRequest) (*pb.VirtioBlkUpdateResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioBlkUpdateResponse{}, nil
}

func (s *server) VirtioBlkList(ctx context.Context, in *pb.VirtioBlkListRequest) (*pb.VirtioBlkListResponse, error) {
	log.Printf("Received from client: %v", in)
	Blobarray := make([]*pb.VirtioBlk, 3)
	return &pb.VirtioBlkListResponse{Controller: Blobarray}, nil
}

func (s *server) VirtioBlkGet(ctx context.Context, in *pb.VirtioBlkGetRequest) (*pb.VirtioBlkGetResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioBlkGetResponse{Controller: &pb.VirtioBlk{Name: "Hello " + fmt.Sprint(in.GetControllerId())}}, nil
}

func (s *server) VirtioBlkStats(ctx context.Context, in *pb.VirtioBlkStatsRequest) (*pb.VirtioBlkStatsResponse, error) {
	log.Printf("Received from client: %v", in)
	return &pb.VirtioBlkStatsResponse{}, nil
}

// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/opiproject/opi-api/storage/proto"
)

//////////////////////////////////////////////////////////

func (s *server) NVMfRemoteControllerConnect(ctx context.Context, in *pb.NVMfRemoteControllerConnectRequest) (*pb.NVMfRemoteControllerConnectResponse, error) {
	log.Printf("NVMfRemoteControllerConnect: Received from client: %v", in)
	params := struct {
		Name string `json:"name"`
	}{
		Name: 		fmt.Sprint("Malloc", in.GetCtrl().GetId()),
	}
	var result []struct {
		Name        string `json:"name"`
		BlockSize   int64  `json:"block_size"`
		NumBlocks   int64  `json:"num_blocks"`
		Uuid        string `json:"uuid"`
	}
	// TODO: bdev_nvme_attach_controller -b Nvme0 -t RDMA -a 192.168.100.1 -f IPv4 -s 4420 -n nqn.2016-06.io.spdk:cnode1
	err := call("bdev_get_bdevs", &params, &result)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	log.Printf("Received from SPDK: %v", result)
	if (len(result) != 1) {
		log.Printf("expecting exactly 1 result")
	}
	return &pb.NVMfRemoteControllerConnectResponse{}, nil
}

func (s *server) NVMfRemoteControllerDisconnect(ctx context.Context, in *pb.NVMfRemoteControllerDisconnectRequest) (*pb.NVMfRemoteControllerDisconnectResponse, error) {
	log.Printf("Received: %v", in.GetId())
	// TODO: rpc.py bdev_nvme_detach_controller Nvme0
	return &pb.NVMfRemoteControllerDisconnectResponse{}, nil
}

func (s *server) NVMfRemoteControllerReset(ctx context.Context, in *pb.NVMfRemoteControllerResetRequest) (*pb.NVMfRemoteControllerResetResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMfRemoteControllerResetResponse{}, nil
}

func (s *server) NVMfRemoteControllerList(ctx context.Context, in *pb.NVMfRemoteControllerListRequest) (*pb.NVMfRemoteControllerListResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMfRemoteControllerListResponse{}, nil
}

func (s *server) NVMfRemoteControllerGet(ctx context.Context, in *pb.NVMfRemoteControllerGetRequest) (*pb.NVMfRemoteControllerGetResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMfRemoteControllerGetResponse{}, nil
}

func (s *server) NVMfRemoteControllerStats(ctx context.Context, in *pb.NVMfRemoteControllerStatsRequest) (*pb.NVMfRemoteControllerStatsResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMfRemoteControllerStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

// TODO: add NULL
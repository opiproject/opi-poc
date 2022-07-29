// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"log"

	pb "github.com/opiproject/opi-api/storage/proto"
)

//////////////////////////////////////////////////////////

func (s *server) NVMfRemoteControllerConnect(ctx context.Context, in *pb.NVMfRemoteControllerConnectRequest) (*pb.NVMfRemoteControllerConnectResponse, error) {
	log.Printf("Received: %v", in.GetController())
	return &pb.NVMfRemoteControllerConnectResponse{}, nil
}

func (s *server) NVMfRemoteControllerDisconnect(ctx context.Context, in *pb.NVMfRemoteControllerDisconnectRequest) (*pb.NVMfRemoteControllerDisconnectResponse, error) {
	log.Printf("Received: %v", in.GetId())
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
// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"log"

	pb "github.com/opiproject/opi-api/storage/proto"
)

func (s *server) NVMfRemoteControllerDisconnect(ctx context.Context, in *pb.NVMfRemoteControllerDisconnectRequest) (*pb.NVMfRemoteControllerDisconnectResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMfRemoteControllerDisconnectResponse{}, nil
}

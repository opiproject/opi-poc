// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"context"
	"log"

	pb "github.com/opiproject/opi-api/security/proto"
	"github.com/google/uuid"
)

func (s *server) IPsecCreate(ctx context.Context, in *pb.IPsecCreateRequest) (*pb.IPsecCreateResponse, error) {
	log.Printf("IPsecCreate: Received: %v", in.GetTunnel())
	ipsec_id := uuid.New()

	ip_ret := pb.IPsecCreateResponse {
		Id: &pb.Uuid {
			Value: ipsec_id.String(),
		},
	}

	return &ip_ret, nil
}

func (s *server) IPsecDelete(ctx context.Context, in *pb.IPsecDeleteRequest) (*pb.IPsecDeleteResponse, error) {
	log.Printf("IPsecDelete: Received: %v", in.GetId())

	ip_ret := pb.IPsecDeleteResponse {
		Id: &pb.Uuid {
			Value: in.GetId().GetValue(),
		},
	}
	return &ip_ret, nil
}

func (s *server) IPsecUpdate(ctx context.Context, in *pb.IPsecUpdateRequest) (*pb.IPsecUpdateResponse, error) {
	log.Printf("IPsecUpdate: Received: %v", in.GetId())
	log.Printf("IPsecUpdate: Received: %v", in.GetTunnel())

	ip_ret := pb.IPsecUpdateResponse {
		Id: &pb.Uuid {
			Value: in.GetId().GetValue(),
		},
	}
	return &ip_ret, nil
}

func (s *server) IPsecList(ctx context.Context, in *pb.IPsecListRequest) (*pb.IPsecListResponse, error) {
	log.Printf("IPsecList: Received: %v", in.GetId())
	return &pb.IPsecListResponse{}, nil
}

func (s *server) IPsecGet(ctx context.Context, in *pb.IPsecGetRequest) (*pb.IPsecGetResponse, error) {
	log.Printf("IPsecGet: Received: %v", in.GetId())

	ip_ret := pb.IPsecGetResponse {
		Id: &pb.Uuid {
			Value: in.GetId().GetValue(),
		},
	}

	return &ip_ret, nil
}

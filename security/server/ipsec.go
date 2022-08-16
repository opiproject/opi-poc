// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"context"
	"log"

	pb "github.com/opiproject/opi-api/security/proto"
)

func (s *server) IPsecVersion(ctx context.Context, in *pb.IPsecVersionReq) (*pb.IPsecVersionResp, error) {
	return nil, nil
}

func (s *server) IPsecStats(ctx context.Context, in *pb.IPsecStatsReq) (*pb.IPsecStatsResp, error) {
	return nil, nil
}

func (s *server) IPsecInitiate(ctx context.Context, in *pb.IPsecInitiateReq) (*pb.IPsecInitiateResp, error) {
	return nil, nil
}

func (s *server) IPsecTerminate(ctx context.Context, in *pb.IPsecTerminateReq) (*pb.IPsecTerminateResp, error) {
	return nil, nil
}

func (s *server) IPsecRekey(ctx context.Context, in *pb.IPsecRekeyReq) (*pb.IPsecRekeyResp, error) {
	return nil, nil
}

func (s *server) IPsecListSas(ctx context.Context, in *pb.IPsecListSasReq) (*pb.IPsecListSasResp, error) {
	return nil, nil
}

func (s *server) IPsecListConns(ctx context.Context, in *pb.IPsecListConnsReq) (*pb.IPsecListConnsResp, error) {
	return nil, nil
}

func (s *server) IPsecListCerts(ctx context.Context, in *pb.IPsecListCertsReq) (*pb.IPsecListCertsResp, error) {
	return nil, nil
}

func (s *server) IPsecLoadConn(ctx context.Context, in *pb.IPsecLoadConnReq) (*pb.IPsecLoadConnResp, error) {
	log.Printf("IPsecLoadConn: Received: %v", in.GetConnection())

	err := loadConn(in)
	if err != nil {
		log.Printf("IPsecLoadConn: Failed %v", err)
		return nil, err
	}

	ip_ret := pb.IPsecLoadConnResp {
		Success: "Yes",
	}

	return &ip_ret, nil
}

func (s *server) IPsecUnloadConn(ctx context.Context, in *pb.IPsecUnloadConnReq) (*pb.IPsecUnloadConnResp, error) {
	log.Printf("IPsecUnloadConn: Received: %v", in.GetName())

	err := unloadConn(in)
	if err != nil {
		log.Printf("IPsecUnloadConn: Failed %v", err)
		return nil, err
	}

	ip_ret := pb.IPsecUnloadConnResp {
		Success: "Yes",
	}

	return &ip_ret, nil
}

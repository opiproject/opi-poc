// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"context"
	"io"
	"log"

	pb "github.com/opiproject/opi-api/security/proto"
)

func (s *server) AddSession(stream pb.SessionTable_AddSessionServer) error {
	var total int32

	for {
		sr, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.AddSessionResponse{
				Requeststatus: pb.AddSessionStatus__SESSION_ACCEPTED,
				Errorstatus: 0,
			})
		}
		if err != nil {
			return err
		}

		total++
		log.Printf("%+v\n", sr)
	}
}

func (s *server) GetSession(ctx context.Context, in *pb.SessionId) (*pb.SessionResponse, error) {
	return &pb.SessionResponse{
			Sessionid: 0,
			Inpackets: 0,
			Outpackets: 0,
			Inbytes: 0,
			Outbytes: 0,
			Sessionstate: pb.SessionState__UNKNOWN_STATE,
		}, nil
}

func (s *server) DeleteSession(ctx context.Context, in *pb.SessionId) (*pb.SessionResponse, error) {
	return &pb.SessionResponse{
			Sessionid: 0,
			Inpackets: 0,
			Outpackets: 0,
			Inbytes: 0,
			Outbytes: 0,
			Sessionstate: pb.SessionState__CLOSED,
		}, nil
}

func (s *server) GetAllSession(ctx context.Context, in *pb.SessionRequestArgs) (*pb.SessionResponses, error) {
	return &pb.SessionResponses{
			Sessioninfo: []*pb.SessionResponse {
				{
					Sessionid: 0,
					Inpackets: 0,
					Outpackets: 0,
					Inbytes: 0,
					Outbytes: 0,
					Sessionstate: pb.SessionState__UNKNOWN_STATE,
				},
			},
		}, nil
}

// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"context"
	fw "github.com/opiproject/opi-api/security/firewall/v1/gen/go"
	"io"
	"log"
)

func (s *server) AddSession(stream fw.SessionTable_AddSessionServer) error {
	var total int32

	for {
		sr, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&fw.AddSessionResponse{
				Requeststatus: fw.AddSessionStatus__SESSION_ACCEPTED,
				Errorstatus:   0,
			})
		}
		if err != nil {
			return err
		}

		total++
		log.Printf("%+v\n", sr)
	}
}

func (s *server) GetSession(ctx context.Context, in *fw.SessionId) (*fw.SessionResponse, error) {
	return &fw.SessionResponse{
		Sessionid:    &fw.Uuid{Value: ""},
		Inpackets:    0,
		Outpackets:   0,
		Inbytes:      0,
		Outbytes:     0,
		Sessionstate: fw.SessionState__UNKNOWN_STATE,
	}, nil
}

func (s *server) DeleteSession(ctx context.Context, in *fw.SessionId) (*fw.SessionResponse, error) {
	return &fw.SessionResponse{
		Sessionid:    &fw.Uuid{Value: ""},
		Inpackets:    0,
		Outpackets:   0,
		Inbytes:      0,
		Outbytes:     0,
		Sessionstate: fw.SessionState__CLOSED,
	}, nil
}

func (s *server) GetAllSession(ctx context.Context, in *fw.SessionRequestArgs) (*fw.SessionResponses, error) {
	return &fw.SessionResponses{
		Sessioninfo: []*fw.SessionResponse{
			{
				Sessionid:    &fw.Uuid{Value: ""},
				Inpackets:    0,
				Outpackets:   0,
				Inbytes:      0,
				Outbytes:     0,
				Sessionstate: fw.SessionState__UNKNOWN_STATE,
			},
		},
	}, nil
}

/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: Apache-2.0
 */

package shared

import (
	"opi-api/proto"
	"golang.org/x/net/context"
)

// GRPCClient is an implementation of OPIAPI that talks over RPC.
type GRPCClient struct{ client proto.OPIAPIClient }

func (m *GRPCClient) CreateNetwork(key string, value []byte) error {
	_, err := m.client.CreateNetwork(context.Background(), &proto.CreateNetworkRequest{
		Key:   key,
		Value: value,
	})
	return err
}

func (m *GRPCClient) GetNetwork(key string) ([]byte, error) {
	resp, err := m.client.GetNetwork(context.Background(), &proto.GetNetworkRequest{
		Key: key,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

// Here is the gRPC server that GRPCClient talks to.
type GRPCServer struct {
	// This is the real implementation
	Impl OPIAPI
}

func (m *GRPCServer) CreateNetwork(
	ctx context.Context,
	req *proto.CreateNetworkRequest) (*proto.Empty, error) {
	return &proto.Empty{}, m.Impl.CreateNetwork(req.Key, req.Value)
}

func (m *GRPCServer) GetNetwork(
	ctx context.Context,
	req *proto.GetNetworkRequest) (*proto.GetNetworkResponse, error) {
	v, err := m.Impl.GetNetwork(req.Key)
	return &proto.GetNetworkResponse{Value: v}, err
}

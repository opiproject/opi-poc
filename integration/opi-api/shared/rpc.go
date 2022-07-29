/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: Apache-2.0
 */

package shared

import (
	"net/rpc"
)

// RPCClient is an implementation of OPIAPI that talks over RPC.
type RPCClient struct{ client *rpc.Client }

func (m *RPCClient) CreateNetwork(key string, value []byte) error {
	// We don't expect a response, so we can just use interface{}
	var resp interface{}

	// The args are just going to be a map. A struct could be better.
	return m.client.Call("Plugin.CreateNetwork", map[string]interface{}{
		"key":   key,
		"value": value,
	}, &resp)
}

func (m *RPCClient) GetNetwork(key string) ([]byte, error) {
	var resp []byte
	err := m.client.Call("Plugin.GetNetwork", key, &resp)
	return resp, err
}

// Here is the RPC server that RPCClient talks to, conforming to
// the requirements of net/rpc
type RPCServer struct {
	// This is the real implementation
	Impl OPIAPI
}

func (m *RPCServer) CreateNetwork(args map[string]interface{}, resp *interface{}) error {
	return m.Impl.CreateNetwork(args["key"].(string), args["value"].([]byte))
}

func (m *RPCServer) GetNetwork(key string, resp *[]byte) error {
	v, err := m.Impl.GetNetwork(key)
	*resp = v
	return err
}

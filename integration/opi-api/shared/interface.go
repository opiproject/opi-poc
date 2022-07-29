/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: Apache-2.0
 */

// Package shared contains shared data between the host and plugins.
package shared

import (
	"context"
	"net/rpc"

	"opi-api/proto"

	"google.golang.org/grpc"

	"github.com/hashicorp/go-plugin"
)

// Handshake is a common handshake that is shared by plugin and host.
var Handshake = plugin.HandshakeConfig{
	// This isn't required when using VersionedPlugins
	ProtocolVersion:  1,
	MagicCookieKey:   "BASIC_PLUGIN",
	MagicCookieValue: "hello",
}

// PluginMap is the map of plugins we can dispense.
var PluginMap = map[string]plugin.Plugin{
	"opi_api_grpc": &OPIAPIGRPCPlugin{},
	"opi_api":      &OPIAPIPlugin{},
}

// OPIAPI is the interface that we're exposing as a plugin.
type OPIAPI interface {
	CreateNetwork(key string, value []byte) error
	GetNetwork(key string) ([]byte, error)
}

// This is the implementation of plugin.Plugin so we can serve/consume this.
type OPIAPIPlugin struct {
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl OPIAPI
}

func (p *OPIAPIPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &RPCServer{Impl: p.Impl}, nil
}

func (*OPIAPIPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &RPCClient{client: c}, nil
}

// This is the implementation of plugin.GRPCPlugin so we can serve/consume this.
type OPIAPIGRPCPlugin struct {
	// GRPCPlugin must still implement the Plugin interface
	plugin.Plugin
	// Concrete implementation, written in Go. This is only used for plugins
	// that are written in Go.
	Impl OPIAPI
}

func (p *OPIAPIGRPCPlugin) GRPCServer(broker *plugin.GRPCBroker, s *grpc.Server) error {
	proto.RegisterOPIAPIServer(s, &GRPCServer{Impl: p.Impl})
	return nil
}

func (p *OPIAPIGRPCPlugin) GRPCClient(ctx context.Context, broker *plugin.GRPCBroker, c *grpc.ClientConn) (interface{}, error) {
	return &GRPCClient{client: proto.NewOPIAPIClient(c)}, nil
}

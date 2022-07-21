/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"io/ioutil"

	"opi-api/shared"

	"github.com/hashicorp/go-plugin"
)

/*
 * This simple implementation is based on [1]. Our implementation uses
 * examples for CreateNetwork() and GetNetwork() to show how we can use
 * plugins to implement the calls. This plugin implements those calls and
 * stores them locally in a file on disk.
 *
 * [1] https://github.com/hashicorp/go-plugin/tree/master/examples/grpc
 */
type OPIAPI struct{}

/* Save the Network name and data to a file */
func (OPIAPI) CreateNetwork(key string, value []byte) error {
	value = []byte(fmt.Sprintf("%s\n\nWritten from plugin-go-grpc", string(value)))
	return ioutil.WriteFile("opiapi_"+key, value, 0644)
}

/* Read the data matching the Network name and return it */
func (OPIAPI) GetNetwork(key string) ([]byte, error) {
	return ioutil.ReadFile("opiapi_" + key)
}

func main() {
	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: shared.Handshake,
		Plugins: map[string]plugin.Plugin{
			"opiapi": &shared.OPIAPIGRPCPlugin{Impl: &OPIAPI{}},
		},

		// A non-nil value here enables gRPC serving for this plugin...
		GRPCServer: plugin.DefaultGRPCServer,
	})
}

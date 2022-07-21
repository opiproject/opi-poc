/*
 * Copyright (C) 2022 Intel Corporation
 * SPDX-License-Identifier: Apache-2.0
 */

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"opi-api/shared"

	"github.com/hashicorp/go-plugin"
)

func main() {
	// We don't want to see the plugin logs.
	log.SetOutput(ioutil.Discard)

	// We're a host. Start by launching the plugin process.
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: shared.Handshake,
		Plugins:         shared.PluginMap,
		Cmd:             exec.Command("sh", "-c", os.Getenv("OPI_PLUGIN")),
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC, plugin.ProtocolGRPC},
	})
	defer client.Kill()

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("opi_api_grpc")
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}

	/*
	 * We should have a OPIAPI now! This feels like a normal interface
	 * implementation but is in fact over an RPC connection.
	 *
	 * NOTE: In this example, the command runs a single time. For the actual
	 *       OPI API Conglomerate, we'd want this to run as a daemon and serve
	 *       gRPC requests on the northbound side to handle the translation to
	 *       these southbound calls.
	 */
	opiapi := raw.(shared.OPIAPI)
	os.Args = os.Args[1:]
	switch os.Args[0] {
	case "GetNetwork":
		result, err := opiapi.GetNetwork(os.Args[1])
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

		fmt.Println(string(result))

	case "CreateNetwork":
		err := opiapi.CreateNetwork(os.Args[1], []byte(os.Args[2]))
		if err != nil {
			fmt.Println("Error:", err.Error())
			os.Exit(1)
		}

	default:
		fmt.Printf("Please only use 'GetNetwork' or 'CreateNetwork', given: %q", os.Args[0])
		os.Exit(1)
	}
	os.Exit(0)
}

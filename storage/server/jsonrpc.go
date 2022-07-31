// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"bufio"
	"io"
	"flag"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync/atomic"
)

var (

	rpcID      int32 // json request message ID, auto incremented
	rpc_sock = flag.String("rpc_sock", "/var/tmp/spdk.sock", "Path to SPDK JSON RPC socket")
)

// low level rpc request/response handling
func call(method string, args, result interface{}) error {
	type rpcRequest struct {
		Ver    string `json:"jsonrpc"`
		ID     int32  `json:"id"`
		Method string `json:"method"`
	}

	id := atomic.AddInt32(&rpcID, 1)
	request := rpcRequest{
		Ver:    "2.0",
		ID:     id,
		Method: method,
	}

	var data []byte
	var err error

	if args == nil {
		data, err = json.Marshal(request)
	} else {
		requestWithParams := struct {
			rpcRequest
			Params interface{} `json:"params"`
		}{
			request,
			args,
		}
		data, err = json.Marshal(requestWithParams)
	}
	if err != nil {
		return fmt.Errorf("%s: %s", method, err)
	}

	fmt.Println(string(data))

	// TODO: add also web option: resp, _ = webSocketCom(rpcClient, data)
	resp, _ := unixSocketCom(*rpc_sock, data)

	response := struct {
		ID    int32 `json:"id"`
		Error struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
		Result interface{} `json:"result"`
	}{
		Result: result,
	}

	err = json.NewDecoder(resp).Decode(&response)
	fmt.Println(response)
	if err != nil {
		return fmt.Errorf("%s: %s", method, err)
	}
	if response.ID != id {
		return fmt.Errorf("%s: json response ID mismatch", method)
	}
	if response.Error.Code != 0 {
		return fmt.Errorf("%s: json response error: %s", method, response.Error.Message)
	}

	return nil
}

func unixSocketCom(rpc_sock string, buf []byte) (io.Reader, error) {
	conn, err := net.Dial("unix", rpc_sock)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(buf)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	err = conn.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return bufio.NewReader(conn), nil
}

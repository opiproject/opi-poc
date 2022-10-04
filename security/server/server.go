// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"flag"
	"fmt"
	fw "github.com/opiproject/opi-api/security/firewall/v1/gen/go"
	"log"
	"net"

	pb "github.com/opiproject/opi-api/security/v1/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50151, "The server port")
)

type server struct {
	pb.UnimplementedIPsecServer
}

type session_server struct {
	fw.UnimplementedSessionTableServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterIPsecServer(s, &server{})
	fw.RegisterSessionTableServer(s, &session_server{})

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

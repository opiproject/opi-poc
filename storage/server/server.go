// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/opiproject/opi-api/storage/v1/gen/go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedNVMeSubsystemServiceServer
	pb.UnimplementedNVMeControllerServiceServer
	pb.UnimplementedNVMeNamespaceServiceServer
	pb.UnimplementedNVMfRemoteControllerServiceServer
	pb.UnimplementedVirtioBlkServiceServer
	pb.UnimplementedVirtioScsiControllerServiceServer
	pb.UnimplementedVirtioScsiLunServiceServer
	pb.UnimplementedNullDebugServiceServer
	pb.UnimplementedAioControllerServiceServer
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterNVMeSubsystemServiceServer(s, &server{})
	pb.RegisterNVMeControllerServiceServer(s, &server{})
	pb.RegisterNVMeNamespaceServiceServer(s, &server{})
	pb.RegisterNVMfRemoteControllerServiceServer(s, &server{})
	pb.RegisterVirtioBlkServiceServer(s, &server{})
	pb.RegisterVirtioScsiControllerServiceServer(s, &server{})
	pb.RegisterVirtioScsiLunServiceServer(s, &server{})
	pb.RegisterNullDebugServiceServer(s, &server{})
	pb.RegisterAioControllerServiceServer(s, &server{})

	reflection.Register(s)

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

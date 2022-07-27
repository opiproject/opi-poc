// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"encoding/json"
	"context"
	"fmt"
	"log"

	pb "github.com/opiproject/opi-api/storage/proto"
)

func (s *server) NVMeSubsystemCreate(ctx context.Context, in *pb.NVMeSubsystemCreateRequest) (*pb.NVMeSubsystemCreateResponse, error) {
	log.Printf("Received: %v", in.GetSubsystem())
	values := map[string]string{"jsonrpc": "2.0", "id": "1", "method": "bdev_get_bdevs"}
	jsonValue, _ := json.Marshal(values)
	jsonReply, _ := spdkCommunicate(jsonValue)
	fmt.Println(string(jsonReply))
	return &pb.NVMeSubsystemCreateResponse{}, nil
}

func (s *server) NVMeSubsystemDelete(ctx context.Context, in *pb.NVMeSubsystemDeleteRequest) (*pb.NVMeSubsystemDeleteResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMeSubsystemDeleteResponse{}, nil
}

func (s *server) NVMeSubsystemGet(ctx context.Context, in *pb.NVMeSubsystemGetRequest) (*pb.NVMeSubsystemGetResponse, error) {
	log.Printf("Received: %v", in.GetId())
	values := map[string]string{"jsonrpc": "2.0", "id": "1", "method": "bdev_get_bdevs"}
	jsonValue, _ := json.Marshal(values)
	jsonReply, _ := spdkCommunicate(jsonValue)
	return &pb.NVMeSubsystemGetResponse{Subsystem: &pb.NVMeSubsystem{NQN: "Hello " + string(in.GetId()) + " got " + string(jsonReply)}}, nil
}

func (s *server) NVMeControllerCreate(ctx context.Context, in *pb.NVMeControllerCreateRequest) (*pb.NVMeControllerCreateResponse, error) {
	log.Printf("Received: %v", in.GetController())
	return &pb.NVMeControllerCreateResponse{}, nil
}

func (s *server) NVMeControllerDelete(ctx context.Context, in *pb.NVMeControllerDeleteRequest) (*pb.NVMeControllerDeleteResponse, error) {
	log.Printf("Received: %v", in.GetControllerId())
	return &pb.NVMeControllerDeleteResponse{}, nil
}

func (s *server) NVMeControllerGet(ctx context.Context, in *pb.NVMeControllerGetRequest) (*pb.NVMeControllerGetResponse, error) {
	log.Printf("Received: %v", in.GetControllerId())
	return &pb.NVMeControllerGetResponse{Controller: &pb.NVMeController{Name: "Hello " + string(in.GetControllerId())}}, nil
}

func (s *server) NVMeNamespaceCreate(ctx context.Context, in *pb.NVMeNamespaceCreateRequest) (*pb.NVMeNamespaceCreateResponse, error) {
	log.Printf("Received: %v", in.GetNamespace())
	return &pb.NVMeNamespaceCreateResponse{}, nil
}

func (s *server) NVMeNamespaceDelete(ctx context.Context, in *pb.NVMeNamespaceDeleteRequest) (*pb.NVMeNamespaceDeleteResponse, error) {
	log.Printf("Received: %v", in.GetNamespaceId())
	return &pb.NVMeNamespaceDeleteResponse{}, nil
}

func (s *server) NVMeNamespaceGet(ctx context.Context, in *pb.NVMeNamespaceGetRequest) (*pb.NVMeNamespaceGetResponse, error) {
	log.Printf("Received: %v", in.GetNamespaceId())
	return &pb.NVMeNamespaceGetResponse{Namespace: &pb.NVMeNamespace{Name: "Hello " + string(in.GetNamespaceId())}}, nil
}


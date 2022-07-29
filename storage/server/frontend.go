// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/opiproject/opi-api/storage/proto"
)

//////////////////////////////////////////////////////////

func (s *server) NVMeSubsystemCreate(ctx context.Context, in *pb.NVMeSubsystemCreateRequest) (*pb.NVMeSubsystemCreateResponse, error) {
	log.Printf("Received: %v", in.GetSubsystem())
	params := struct {Name string `json:"name"`}{ Name: "Malloc0"}
	var Result string
	call("bdev_get_bdevs", &params, &Result)
	fmt.Println(Result)
	return &pb.NVMeSubsystemCreateResponse{}, nil
}

func (s *server) NVMeSubsystemDelete(ctx context.Context, in *pb.NVMeSubsystemDeleteRequest) (*pb.NVMeSubsystemDeleteResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMeSubsystemDeleteResponse{}, nil
}

func (s *server) NVMeSubsystemUpdate(ctx context.Context, in *pb.NVMeSubsystemUpdateRequest) (*pb.NVMeSubsystemUpdateResponse, error) {
	log.Printf("Received: %v", in.GetSubsystem())
	params := struct {Name string `json:"name"`}{ Name: "Malloc0"}
	var Result string
	call("bdev_get_bdevs", &params, &Result)
	fmt.Println(Result)
	return &pb.NVMeSubsystemUpdateResponse{}, nil
}

func (s *server) NVMeSubsystemList(ctx context.Context, in *pb.NVMeSubsystemListRequest) (*pb.NVMeSubsystemListResponse, error) {
	log.Printf("Received: %v", in)
	Blobarray := make([]*pb.NVMeSubsystem, 3)
	return &pb.NVMeSubsystemListResponse{Subsystem: Blobarray}, nil
}

func (s *server) NVMeSubsystemGet(ctx context.Context, in *pb.NVMeSubsystemGetRequest) (*pb.NVMeSubsystemGetResponse, error) {
	log.Printf("Received: %v", in.GetId())
	params := struct {Name string `json:"name"`}{ Name: "Malloc0"}
	var Result string
	call("bdev_get_bdevs", &params, &Result)
	fmt.Println(Result)
	return &pb.NVMeSubsystemGetResponse{Subsystem: &pb.NVMeSubsystem{NQN: "Hello " + string(in.GetId()) + " got " + string(Result)}}, nil
}

func (s *server) NVMeSubsystemStats(ctx context.Context, in *pb.NVMeSubsystemStatsRequest) (*pb.NVMeSubsystemStatsResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMeSubsystemStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

func (s *server) NVMeControllerCreate(ctx context.Context, in *pb.NVMeControllerCreateRequest) (*pb.NVMeControllerCreateResponse, error) {
	log.Printf("Received: %v", in.GetController())
	return &pb.NVMeControllerCreateResponse{}, nil
}

func (s *server) NVMeControllerDelete(ctx context.Context, in *pb.NVMeControllerDeleteRequest) (*pb.NVMeControllerDeleteResponse, error) {
	log.Printf("Received: %v", in.GetControllerId())
	return &pb.NVMeControllerDeleteResponse{}, nil
}

func (s *server) NVMeControllerUpdate(ctx context.Context, in *pb.NVMeControllerUpdateRequest) (*pb.NVMeControllerUpdateResponse, error) {
	log.Printf("Received: %v", in.GetController())
	return &pb.NVMeControllerUpdateResponse{}, nil
}

func (s *server) NVMeControllerList(ctx context.Context, in *pb.NVMeControllerListRequest) (*pb.NVMeControllerListResponse, error) {
	log.Printf("Received: %v", in.GetSubsystemId())
	Blobarray := make([]*pb.NVMeController, 3)
	return &pb.NVMeControllerListResponse{Controller: Blobarray}, nil
}

func (s *server) NVMeControllerGet(ctx context.Context, in *pb.NVMeControllerGetRequest) (*pb.NVMeControllerGetResponse, error) {
	log.Printf("Received: %v", in.GetControllerId())
	return &pb.NVMeControllerGetResponse{Controller: &pb.NVMeController{Name: "Hello " + string(in.GetControllerId())}}, nil
}

func (s *server) NVMeControllerStats(ctx context.Context, in *pb.NVMeControllerStatsRequest) (*pb.NVMeControllerStatsResponse, error) {
	log.Printf("Received: %v", in.GetControllerId())
	return &pb.NVMeControllerStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

func (s *server) NVMeNamespaceCreate(ctx context.Context, in *pb.NVMeNamespaceCreateRequest) (*pb.NVMeNamespaceCreateResponse, error) {
	log.Printf("Received: %v", in.GetNamespace())
	return &pb.NVMeNamespaceCreateResponse{}, nil
}

func (s *server) NVMeNamespaceDelete(ctx context.Context, in *pb.NVMeNamespaceDeleteRequest) (*pb.NVMeNamespaceDeleteResponse, error) {
	log.Printf("Received: %v", in.GetNamespaceId())
	return &pb.NVMeNamespaceDeleteResponse{}, nil
}

func (s *server) NVMeNamespaceUpdate(ctx context.Context, in *pb.NVMeNamespaceUpdateRequest) (*pb.NVMeNamespaceUpdateResponse, error) {
	log.Printf("Received: %v", in.GetNamespace())
	return &pb.NVMeNamespaceUpdateResponse{}, nil
}

func (s *server) NVMeNamespaceList(ctx context.Context, in *pb.NVMeNamespaceListRequest) (*pb.NVMeNamespaceListResponse, error) {
	log.Printf("Received: %v", in.GetControllerId())
	Blobarray := make([]*pb.NVMeNamespace, 3)
	return &pb.NVMeNamespaceListResponse{Namespace: Blobarray}, nil
}

func (s *server) NVMeNamespaceGet(ctx context.Context, in *pb.NVMeNamespaceGetRequest) (*pb.NVMeNamespaceGetResponse, error) {
	log.Printf("Received: %v", in.GetNamespaceId())
	return &pb.NVMeNamespaceGetResponse{Namespace: &pb.NVMeNamespace{Name: "Hello " + string(in.GetNamespaceId())}}, nil
}

func (s *server) NVMeNamespaceStats(ctx context.Context, in *pb.NVMeNamespaceStatsRequest) (*pb.NVMeNamespaceStatsResponse, error) {
	log.Printf("Received: %v", in.GetNamespaceId())
	return &pb.NVMeNamespaceStatsResponse{}, nil
}

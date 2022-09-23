// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Dell Inc, or its subsidiaries.

package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/opiproject/opi-api/storage/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//////////////////////////////////////////////////////////

func (s *server) NVMfRemoteControllerConnect(ctx context.Context, in *pb.NVMfRemoteControllerConnectRequest) (*pb.NVMfRemoteControllerConnectResponse, error) {
	log.Printf("NVMfRemoteControllerConnect: Received from client: %v", in)
	params := struct {
		Name        string `json:"name"`
		Type        string `json:"trtype"`
		Address     string `json:"traddr"`
		Family      string `json:"adrfam"`
		Port        string `json:"trsvcid"`
		Subsystem   string `json:"subnqn"`
	}{
		Name: 		fmt.Sprint("OpiNvme", in.GetCtrl().GetId()),
		Type: 		"TCP",
		Address: 	in.GetCtrl().GetTraddr(),
		Family:		"ipv4",
		Port: 		fmt.Sprint(in.GetCtrl().GetTrsvcid()),
		Subsystem: 	in.GetCtrl().GetSubnqn(),
	}
	var result []struct {
		Name        string `json:"name"`
	}
	err := call("bdev_nvme_attach_controller", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if (len(result) != 1) {
		log.Printf("expecting exactly 1 result")
	}
	return &pb.NVMfRemoteControllerConnectResponse{}, nil
}

func (s *server) NVMfRemoteControllerDisconnect(ctx context.Context, in *pb.NVMfRemoteControllerDisconnectRequest) (*pb.NVMfRemoteControllerDisconnectResponse, error) {
	log.Printf("NVMfRemoteControllerDisconnect: Received from client: %v", in)
	params := struct {
		Name string `json:"name"`
	}{
		Name: 		fmt.Sprint("OpiNvme", in.GetId()),
	}
	var result bool
	err := call("bdev_nvme_detach_controller", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	return &pb.NVMfRemoteControllerDisconnectResponse{}, nil
}

func (s *server) NVMfRemoteControllerReset(ctx context.Context, in *pb.NVMfRemoteControllerResetRequest) (*pb.NVMfRemoteControllerResetResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMfRemoteControllerResetResponse{}, nil
}

func (s *server) NVMfRemoteControllerList(ctx context.Context, in *pb.NVMfRemoteControllerListRequest) (*pb.NVMfRemoteControllerListResponse, error) {
	log.Printf("NVMfRemoteControllerList: Received from client: %v", in)
	var result []struct {
		Name   string `json:"name"`
		Ctrlrs []struct {
			State string `json:"state"`
			Trid  struct {
				Trtype  string `json:"trtype"`
				Adrfam  string `json:"adrfam"`
				Traddr  string `json:"traddr"`
				Trsvcid string `json:"trsvcid"`
				Subnqn  string `json:"subnqn"`
			} `json:"trid"`
			Cntlid int `json:"cntlid"`
			Host   struct {
				Nqn   string `json:"nqn"`
				Addr  string `json:"addr"`
				Svcid string `json:"svcid"`
			} `json:"host"`
		} `json:"ctrlrs"`
	}
	err := call("bdev_nvme_get_controllers", nil, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	Blobarray := make([]*pb.NVMfRemoteController, len(result))
	for i := range result {
		r := &result[i]
		Blobarray[i] = &pb.NVMfRemoteController{Subnqn: r.Name}
	}
	return &pb.NVMfRemoteControllerListResponse{Ctrl: Blobarray}, nil
}

func (s *server) NVMfRemoteControllerGet(ctx context.Context, in *pb.NVMfRemoteControllerGetRequest) (*pb.NVMfRemoteControllerGetResponse, error) {
	log.Printf("NVMfRemoteControllerGet: Received from client: %v", in)
	params := struct {
		Name string `json:"name"`
	}{
		Name:       fmt.Sprint("OpiNvme", in.GetId()),
	}
	var result []struct {
		Name        string `json:"name"`
		CtrlId      int64  `json:"cntlid"`
	}
	err := call("bdev_nvme_get_controllers", &params, &result)
	if err != nil {
		log.Printf("error: %v", err)
		return nil, err
	}
	log.Printf("Received from SPDK: %v", result)
	if (len(result) != 1) {
		msg := fmt.Sprintf("expecting exactly 1 result, got %d", len(result))
		log.Print(msg)
		return nil, status.Errorf(codes.InvalidArgument, msg)
	}
	return &pb.NVMfRemoteControllerGetResponse{Ctrl: &pb.NVMfRemoteController{Subnqn: result[0].Name}}, nil
}

func (s *server) NVMfRemoteControllerStats(ctx context.Context, in *pb.NVMfRemoteControllerStatsRequest) (*pb.NVMfRemoteControllerStatsResponse, error) {
	log.Printf("Received: %v", in.GetId())
	return &pb.NVMfRemoteControllerStatsResponse{}, nil
}

//////////////////////////////////////////////////////////

// TODO: add NULL
// TODO: add AIO
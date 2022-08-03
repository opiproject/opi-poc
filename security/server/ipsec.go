// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2022 Intel Corporation, or its subsidiaries.

package main

import (
	"context"
	"encoding/base64"
	"log"

	"github.com/go-redis/redis/v8"
	"google.golang.org/protobuf/proto"
	pb "github.com/opiproject/opi-api/security/proto"
	"github.com/google/uuid"
)

// Our client connection
var rdb *redis.Client

// Used for background Redis processing
var ctx = context.Background()

func init_backend() error {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	val, err := rdb.Info(ctx).Result()
	if err != nil {
		panic(err)
	}
	log.Printf("Redis Info:\n%v\n", val)

	return nil
}

func (s *server) IPsecCreate(ctx context.Context, in *pb.IPsecCreateRequest) (*pb.IPsecCreateResponse, error) {
	log.Printf("IPsecCreate: Received: %v", in.GetTunnel())
	ipsec_id := uuid.New()

	// Store in Redis
	b, err := proto.Marshal(in)
	if err != nil {
		log.Fatal("Cannot marshal data")
	}
	proto_str := base64.StdEncoding.EncodeToString(b)

	log.Printf("Saving entry with UUID %s", ipsec_id.String())

	err = rdb.Set(ctx, ipsec_id.String(), proto_str, 0).Err()
	if err != nil {
		panic(err)
	}

	err = load_connection("load", in)
	if err != nil {
		log.Printf("IPsecCreate: Failed with error %v", err)
		return nil, err
	}

	ip_ret := pb.IPsecCreateResponse {
		Id: &pb.Uuid {
			Value: ipsec_id.String(),
		},
	}

	return &ip_ret, nil
}

func (s *server) IPsecDelete(ctx context.Context, in *pb.IPsecDeleteRequest) (*pb.IPsecDeleteResponse, error) {
	reqId := in.GetId().GetValue()
	log.Printf("IPsecDelete: Received: %v", reqId)

	log.Printf("Retreiving entry with UUID %s", reqId)

	// Retreive from Redis
	val, err := rdb.Get(ctx, reqId).Result()
	if err != nil {
		panic(err)
	}
	p_blob, err := base64.StdEncoding.DecodeString(val)
	if err != nil {
		log.Fatal("Cannot decode data")
	}
	ipsec_req := &pb.IPsecCreateRequest{}
	err = proto.Unmarshal(p_blob, ipsec_req)
	if err != nil {
		log.Fatal("Cannot unmarshal data")
	}

	log.Printf("Dumping unmarshaled protobuf\n%v\n", ipsec_req)

	err = delete_connection("load", ipsec_req)
	if err != nil {
		log.Printf("IPsecDelete: Failed with error %v", err)
		return nil, err
	}

	// Delete from Redis
	rdb.Del(ctx, reqId)

	ip_ret := pb.IPsecDeleteResponse {
		Id: &pb.Uuid {
			Value: in.GetId().GetValue(),
		},
	}
	return &ip_ret, nil
}

func (s *server) IPsecUpdate(ctx context.Context, in *pb.IPsecUpdateRequest) (*pb.IPsecUpdateResponse, error) {
	reqId := in.GetId().GetValue()
	log.Printf("IPsecUpdate: Received: %v", reqId)
	log.Printf("IPsecUpdate: Received: %v", in.GetTunnel())

	ip_ret := pb.IPsecUpdateResponse {
		Id: &pb.Uuid {
			Value: in.GetId().GetValue(),
		},
	}
	return &ip_ret, nil
}

func (s *server) IPsecList(ctx context.Context, in *pb.IPsecListRequest) (*pb.IPsecListResponse, error) {
	reqId := in.GetId().GetValue()
	log.Printf("IPsecList: Received: %v", reqId)
	return &pb.IPsecListResponse{}, nil
}

func (s *server) IPsecGet(ctx context.Context, in *pb.IPsecGetRequest) (*pb.IPsecGetResponse, error) {
	reqId := in.GetId().GetValue()
	log.Printf("IPsecGet: Received: %v", reqId)

	ip_ret := pb.IPsecGetResponse {
		Id: &pb.Uuid {
			Value: in.GetId().GetValue(),
		},
	}

	return &ip_ret, nil
}

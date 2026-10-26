// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2026 Keysight Technologies

// Package ipsecbatch benchmarks unary versus batch IPsec SA RPCs against an in-memory SAD.
package ipsecbatch

import (
	"context"
	"sync"

	pb "github.com/opiproject/opi-api/v1/gen/go/security"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

const saPrefix = "ipsecSas/"

// Server is an in-memory IpsecSaService.
type Server struct {
	pb.UnimplementedIpsecSaServiceServer

	mu  sync.Mutex
	sas map[string][]byte
}

// NewServer returns an empty in-memory IpsecSaService.
func NewServer() *Server {
	return &Server{sas: make(map[string][]byte)}
}

// Len returns the number of SAs currently stored.
func (s *Server) Len() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.sas)
}

func validateCreate(req *pb.CreateIpsecSaRequest) error {
	if req.GetIpsecSaId() == "" {
		return status.Error(codes.InvalidArgument, "ipsec_sa_id is required")
	}
	if req.GetIpsecSa().GetConfig() == nil {
		return status.Error(codes.InvalidArgument, "ipsec_sa.config is required")
	}
	return nil
}

func newSa(req *pb.CreateIpsecSaRequest) *pb.IpsecSa {
	return &pb.IpsecSa{
		Name:   saPrefix + req.GetIpsecSaId(),
		Reqid:  req.GetIpsecSa().GetReqid(),
		Config: req.GetIpsecSa().GetConfig(),
	}
}

// Reset removes all SAs.
func (s *Server) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sas = make(map[string][]byte)
}

// CreateIpsecSa creates a single SA.
func (s *Server) CreateIpsecSa(_ context.Context, req *pb.CreateIpsecSaRequest) (*pb.IpsecSa, error) {
	if err := validateCreate(req); err != nil {
		return nil, err
	}
	sa := newSa(req)
	entry, err := proto.Marshal(sa)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sas[sa.Name]; ok {
		return nil, status.Errorf(codes.AlreadyExists, "%s already exists", sa.Name)
	}
	s.sas[sa.Name] = entry
	return sa, nil
}

// DeleteIpsecSa deletes a single SA.
func (s *Server) DeleteIpsecSa(_ context.Context, req *pb.DeleteIpsecSaRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sas[req.GetName()]; !ok {
		return nil, status.Errorf(codes.NotFound, "%s not found", req.GetName())
	}
	delete(s.sas, req.GetName())
	return &emptypb.Empty{}, nil
}

// BatchCreateIpsecSas creates all the SAs in the request, or none of them.
func (s *Server) BatchCreateIpsecSas(_ context.Context, req *pb.BatchCreateIpsecSasRequest) (*pb.BatchCreateIpsecSasResponse, error) {
	sas := make([]*pb.IpsecSa, len(req.GetRequests()))
	entries := make([][]byte, len(sas))
	for i, r := range req.GetRequests() {
		if err := validateCreate(r); err != nil {
			return nil, err
		}
		sas[i] = newSa(r)
		var err error
		if entries[i], err = proto.Marshal(sas[i]); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	seen := make(map[string]struct{}, len(sas))
	for _, sa := range sas {
		if _, ok := s.sas[sa.Name]; ok {
			return nil, status.Errorf(codes.AlreadyExists, "%s already exists", sa.Name)
		}
		if _, ok := seen[sa.Name]; ok {
			return nil, status.Errorf(codes.InvalidArgument, "%s appears twice in the batch", sa.Name)
		}
		seen[sa.Name] = struct{}{}
	}
	for i, sa := range sas {
		s.sas[sa.Name] = entries[i]
	}
	return &pb.BatchCreateIpsecSasResponse{IpsecSas: sas}, nil
}

// BatchDeleteIpsecSas deletes all the named SAs, or none of them.
func (s *Server) BatchDeleteIpsecSas(_ context.Context, req *pb.BatchDeleteIpsecSasRequest) (*emptypb.Empty, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, name := range req.GetNames() {
		if _, ok := s.sas[name]; !ok {
			return nil, status.Errorf(codes.NotFound, "%s not found", name)
		}
	}
	for _, name := range req.GetNames() {
		delete(s.sas, name)
	}
	return &emptypb.Empty{}, nil
}

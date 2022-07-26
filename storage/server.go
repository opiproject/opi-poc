package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "opi.storage.v1"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedNVMeSubsystemServer
	pb.UnimplementedNVMeControllerServer
	pb.UnimplementedNVMeNamespaceServer
	pb.UnimplementedNVMfRemoteControllerServer
}

func (s *server) NVMeSubsystemCreate(ctx context.Context, in *pb.NVMeSubsystemCreateRequest) (*pb.NVMeSubsystemCreateResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeSubsystemCreateResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeSubsystemDelete(ctx context.Context, in *pb.NVMeSubsystemDeleteRequest) (*pb.NVMeSubsystemDeleteResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeSubsystemDeleteResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeSubsystemGet(ctx context.Context, in *pb.NVMeSubsystemGetRequest) (*pb.NVMeSubsystemGetResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeSubsystemGetResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeControllerCreate(ctx context.Context, in *pb.NVMeControllerCreateRequest) (*pb.NVMeControllerCreateResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeControllerCreateResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeControllerDelete(ctx context.Context, in *pb.NVMeControllerDeleteRequest) (*pb.NVMeControllerDeleteResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeControllerDeleteResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeControllerGet(ctx context.Context, in *pb.NVMeControllerGetRequest) (*pb.NVMeControllerGetResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeControllerGetResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeNamespaceCreate(ctx context.Context, in *pb.NVMeNamespaceCreateRequest) (*pb.NVMeNamespaceCreateResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeNamespaceCreateResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeNamespaceDelete(ctx context.Context, in *pb.NVMeNamespaceDeleteRequest) (*pb.NVMeNamespaceDeleteResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeNamespaceDeleteResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeNamespaceGet(ctx context.Context, in *pb.NVMeNamespaceGetRequest) (*pb.NVMeNamespaceGetResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeNamespaceGetResponse{Message: "Hello " + in.GetName()}, nil
}

func (s *server) NVMfRemoteControllerDisconnect(ctx context.Context, in *pb.NVMfRemoteControllerDisconnectRequest) (*pb.NVMfRemoteControllerDisconnectResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMfRemoteControllerDisconnectResponse{Message: "Hello " + in.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterNVMeSubsystemServer(s, &server{})
	pb.RegisterNVMeControllerServer(s, &server{})
	pb.RegisterNVMeNamespaceServer(s, &server{})
	pb.RegisterNVMfRemoteControllerServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
package main

import (
	"encoding/json"
	"io/ioutil"
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/opiproject/opi-api/storage/proto"
)

var (
	port = flag.Int("port", 50051, "The server port")
	rpc_sock = flag.String("rpc_sock", "/var/tmp/spdk.sock", "Path to SPDK JSON RPC socket")
)

type server struct {
	pb.UnimplementedNVMeSubsystemServer
	pb.UnimplementedNVMeControllerServer
	pb.UnimplementedNVMeNamespaceServer
	pb.UnimplementedNVMfRemoteControllerServer
}

func spdkCommunicate(buf []byte) ([]byte, error) {
	// TODO: use rpc_sock variable
	conn, err := net.Dial("unix", *rpc_sock)
	if err != nil {
		log.Fatal(err)
	}
	_, err = conn.Write(buf)
	if err != nil {
		log.Fatal(err)
	}

	err = conn.(*net.UnixConn).CloseWrite()
	if err != nil {
		log.Fatal(err)
	}

	reply, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(reply))

	return reply, err
}

func (s *server) NVMeSubsystemCreate(ctx context.Context, in *pb.NVMeSubsystemCreateRequest) (*pb.NVMeSubsystemCreateResponse, error) {
	log.Printf("Received: %v", in.GetName())
	// send NVMeSubsystemCreate to SPDK via json rpc unix socket
	values := map[string]string{"jsonrpc": "2.0", "id": "1", "method": "bdev_get_bdevs"}
	jsonValue, _ := json.Marshal(values)
	jsonReply, _ := spdkCommunicate(jsonValue)
	return &pb.NVMeSubsystemCreateResponse{Name: "Hello " + in.GetName() + " got " + string(jsonReply)}, nil
}

func (s *server) NVMeSubsystemDelete(ctx context.Context, in *pb.NVMeSubsystemDeleteRequest) (*pb.NVMeSubsystemDeleteResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeSubsystemDeleteResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeSubsystemGet(ctx context.Context, in *pb.NVMeSubsystemGetRequest) (*pb.NVMeSubsystemGetResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeSubsystemGetResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeControllerCreate(ctx context.Context, in *pb.NVMeControllerCreateRequest) (*pb.NVMeControllerCreateResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeControllerCreateResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeControllerDelete(ctx context.Context, in *pb.NVMeControllerDeleteRequest) (*pb.NVMeControllerDeleteResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeControllerDeleteResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeControllerGet(ctx context.Context, in *pb.NVMeControllerGetRequest) (*pb.NVMeControllerGetResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeControllerGetResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeNamespaceCreate(ctx context.Context, in *pb.NVMeNamespaceCreateRequest) (*pb.NVMeNamespaceCreateResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeNamespaceCreateResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeNamespaceDelete(ctx context.Context, in *pb.NVMeNamespaceDeleteRequest) (*pb.NVMeNamespaceDeleteResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeNamespaceDeleteResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *server) NVMeNamespaceGet(ctx context.Context, in *pb.NVMeNamespaceGetRequest) (*pb.NVMeNamespaceGetResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMeNamespaceGetResponse{Name: "Hello " + in.GetName()}, nil
}

func (s *server) NVMfRemoteControllerDisconnect(ctx context.Context, in *pb.NVMfRemoteControllerDisconnectRequest) (*pb.NVMfRemoteControllerDisconnectResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pb.NVMfRemoteControllerDisconnectResponse{Name: "Hello " + in.GetName()}, nil
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
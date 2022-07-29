package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func do_middleend(conn grpc.ClientConnInterface, ctx context.Context) {
	log.Printf("Test middleend")
}
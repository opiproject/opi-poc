package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
)

func doMiddleend(ctx context.Context, conn grpc.ClientConnInterface) {
	log.Printf("Test middleend")
}

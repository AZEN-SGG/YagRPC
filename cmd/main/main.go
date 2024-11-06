package main

import (
	"context"
	"gRPCOrderService/internal/config"
	"gRPCOrderService/internal/transport/grpc"
)

func main() {
	ctx := context.Background()
	cfg := config.New("./configs/local.env")
	if cfg == nil {
		panic("failed to create config")
	}

	gRPCServer, err := grpc.New(ctx, cfg.GRPCServerPort)
	if err != nil {
		panic(err)
		return
	}

	if err = gRPCServer.Start(ctx); err != nil {
		panic(err)
	}
}

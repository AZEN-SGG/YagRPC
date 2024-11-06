package grpc

import (
	"context"
	"fmt"
	test "gRPCOrderService/pkg/api/3_1"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	GrpcServer *grpc.Server
	listener   net.Listener
}

func New(_ context.Context, port int) (*Server, error) {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption

	grpcServer := grpc.NewServer(opts...)
	test.RegisterOrderServiceServer(grpcServer, NewOrderServiceServerAPI())

	return &Server{
		GrpcServer: grpcServer,
		listener:   lis,
	}, nil
}

func (s *Server) Start(_ context.Context) error {
	return s.GrpcServer.Serve(s.listener)
}

package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	server *grpc.Server
	port   int
}

// NewGRPCServer creates a new GRPCServer
func NewGRPCServer(port int) *GRPCServer {
	gRPCServer := grpc.NewServer()

	return &GRPCServer{
		server: gRPCServer,
		port:   port,
	}
}

// Run starts the grpc server
func (s *GRPCServer) Run() error {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	err = s.server.Serve(listen)
	if err != nil {
		return err
	}

	return nil
}

// MustRun runs gRPC server and panics if any error occurs.
func (s *GRPCServer) MustRun() {
	if err := s.Run(); err != nil {
		panic(err)
	}
}

// Stop gracefuly stops the server
func (s *GRPCServer) Stop() {
	s.server.GracefulStop()
}

package client

import (
	ssov1 "protos/auth"

	"google.golang.org/grpc"
)

type gRPCClientAuthorization struct {
	ssov1.AuthClient
}

func NewGRPCClientAuthorization(dsn string) (*gRPCClientAuthorization, error) {
	conn, err := grpc.NewClient(dsn)
	if err != nil {
		return nil, err
	}

	client := ssov1.NewAuthClient(conn)

	return &gRPCClientAuthorization{client}, nil
}

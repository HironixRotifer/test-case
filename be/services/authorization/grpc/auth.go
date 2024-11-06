package grpc

import (
	ssov1 "protos/auth"

	"google.golang.org/grpc"
)

type Auth struct {
}

type authServer struct {
	ssov1.AuthServer
}

func NewAuth() *Auth {

	gRPCServer := grpc.NewServer()

	ssov1.RegisterAuthServer(gRPCServer, &authServer{})

	return &Auth{}
}

func (a *Auth) Login() {
}

func (a *Auth) Registr() {

}

func (a *Auth) IsAdmin() {

}

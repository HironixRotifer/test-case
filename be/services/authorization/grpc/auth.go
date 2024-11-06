package grpc

import (
	"errors"
	ssov1 "protos/auth"

	"authorization/internal/service"
	"authorization/internal/storage"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(context.Context, string, string) (string, error)
	RegisterNewUser(context.Context, string, string, string, string, string) (userID int64, err error)
	IsAdmin(context.Context, int64) (bool, error)
}

type authServerAPI struct {
	ssov1.AuthServer
	auth Auth
}

func Register(gRPCServer *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPCServer, &authServerAPI{auth: auth})

}

func (a *authServerAPI) Login(ctx context.Context, in *ssov1.LoginRequest) (*ssov1.LoginResponse, error) {
	resp := &ssov1.LoginResponse{}

	if in.Email == "" {
		respError := &ssov1.Error{
			Code:    "3",
			Message: "email is required",
		}
		resp.Error = respError
		return resp, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		respError := &ssov1.Error{
			Code:    "3",
			Message: "password is required",
		}
		resp.Error = respError
		return resp, status.Error(codes.InvalidArgument, "email or password is required")
	}

	token, err := a.auth.Login(ctx, in.Email, in.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid email or password")
		}

		return resp, status.Error(codes.Internal, "login failed")
	}

	resp.Token = token

	return resp, nil
}

func (a *authServerAPI) Registr(ctx context.Context, in *ssov1.RegisterRequest) (*ssov1.RegisterResponse, error) {

	if in.FirstName == "" {
		return nil, status.Error(codes.InvalidArgument, "first_name is required")
	}

	if in.LastName == "" {
		return nil, status.Error(codes.InvalidArgument, "last_name is required")
	}

	if in.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if in.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if in.PhoneNumber == "" {
		return nil, status.Error(codes.InvalidArgument, "phone_number is required")
	}

	userID, err := a.auth.RegisterNewUser(ctx, in.FirstName, in.LastName, in.Email, in.PhoneNumber, in.Password)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &ssov1.RegisterResponse{UserId: userID}, nil
}

func (a *authServerAPI) IsAdmin(ctx context.Context, in *ssov1.IsAdminRequest) (*ssov1.IsAdminResponse, error) {
	resp := &ssov1.IsAdminResponse{}

	if in.UserId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "userID is required")
	}

	isAdmin, err := a.auth.IsAdmin(ctx, in.UserId)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}

		return nil, status.Errorf(codes.Internal, "failed to check admin status")
	}

	resp.IsAdmin = isAdmin

	return resp, nil
}

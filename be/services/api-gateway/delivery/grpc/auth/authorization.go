package authClient

import (
	ssov1 "protos/auth"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
)

type AuthClient struct {
	ssov1.AuthClient
}

func NewAuthGRPClient(destination string) AuthClient {
	conn, err := grpc.NewClient(destination)
	if err != nil {
		log.Error().Err(err).Msg("failed connection to grpc server")
		return AuthClient{}
	}

	client := ssov1.NewAuthClient(conn)

	return AuthClient{client}
}

// NewGRPClient создаёт gRPC подключение к серверу и возвращает клиент
func NewGRPClient(destination string) ssov1.AuthClient {
	conn, err := grpc.NewClient(destination)
	if err != nil {
		log.Error().Err(err).Msg("failed connection to grpc server")
		return nil
	}

	client := ssov1.NewAuthClient(conn)

	return client
}

// LoginUser - вызывает gRPC обработчик авторизации пользователя.
// Успешная авторизация возвращает Json Web Token пользователя.
// В случае провала возвращает и логирует ошибку.
// func (c *AuthClientHandler) HandleLoginUser(ctx context.Context, email string, password string) (Token string, err error) {
// 	resp, err := c.Login(ctx, &ssov1.LoginRequest{
// 		Email:    email,
// 		Password: password,
// 	})

// 	if err != nil {
// 		return "", err
// 	}

// 	return resp.Token, nil
// }

// func (c *AuthClientHandler) HandleRegistrationNewUser(
// 	ctx context.Context,
// 	firstName string,
// 	lastName string,
// 	email string,
// 	phoneNumber string,
// 	password string,
// ) (
// 	userID int64, err error,
// ) {

// 	logger := log.With().Str("service", "authorization").Str("method", "Register").Logger()
// 	resp, err := c.Register(ctx, &ssov1.RegisterRequest{
// 		FirstName:   firstName,
// 		LastName:    lastName,
// 		Email:       email,
// 		Password:    password,
// 		PhoneNumber: phoneNumber,
// 	})

// 	if err != nil {
// 		logger.Err(err).Msg("error to call method to register user")

// 		return 0, err
// 	}

// 	return resp.UserId, nil
// }

// // HandleIsAdmin вызывает gRPC обработчик проверки статуса администратора у пользователя.
// // Успешная проверка возвращает.
// // В случае провала возвращает событие False и логирует ошибку.
// func (c *AuthClientHandler) HandleIsAdmin(ctx context.Context, id int64) (bool, error) {
// 	response, err := c.IsAdmin(ctx, &ssov1.IsAdminRequest{
// 		UserId: id,
// 	})

// 	if err != nil {
// 		return false, err
// 	}

// 	return response.IsAdmin, nil
// }

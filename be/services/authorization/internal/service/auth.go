package service

import (
	"authorization/internal/lib/hashpassword"
	"authorization/internal/lib/jwt"
	"authorization/internal/models"
	"context"
	"errors"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type UserCreator interface {
	CreateUser(ctx context.Context, user models.User) (int64, error)
}

type UserProvider interface {
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	GetUserByID(ctx context.Context, id int64) (models.User, error)
	IsAdmin(ctx context.Context, id int64) (bool, error)
}

type Auth struct {
	userCreator  UserCreator
	userProvider UserProvider
}

// NewAuth creates a new instance of the Auth service.
func NewAuth(userCreator UserCreator, userProvider UserProvider) *Auth {
	return &Auth{
		userCreator:  userCreator,
		userProvider: userProvider,
	}
}

func (a *Auth) Login(ctx context.Context, email string, password string) (token string, err error) {

	user, err := a.userProvider.GetUserByEmail(ctx, email)
	if err != nil {
		// log
		return "", err
	}

	hashedPassword := hashpassword.HashPassword(password, user.Salt)

	token, err = a.Login(ctx, email, hashedPassword)
	if err != nil {
		// log
		return "", err
	}

	return token, nil
}

func (a *Auth) RegisterNewUser(
	ctx context.Context,
	firstName string,
	lastName string,
	email string,
	phoneNumber string,
	password string,
) (userID int64, err error) {

	salt := hashpassword.GenerateRandomSalt()
	hashedPassword := hashpassword.HashPassword(password, salt)
	_, refreshToken, err := jwt.GenerateToken(firstName, lastName, email, phoneNumber)
	if err != nil {
		// log
		return 0, err
	}

	user := models.User{
		FirstName:    firstName,
		LastName:     lastName,
		Email:        email,
		PhoneNumber:  phoneNumber,
		HashPassword: hashedPassword,
		RefreshToken: refreshToken,
	}

	userID, err = a.userCreator.CreateUser(ctx, user)
	if err != nil {
		// log
		return 0, err
	}

	return userID, nil
}

func (a *Auth) IsAdmin(ctx context.Context, id int64) (bool, error) {
	isAdmin, err := a.userProvider.IsAdmin(ctx, id)
	if err != nil {
		// log
		return false, err
	}

	return isAdmin, nil
}

func (a *Auth) UpdateUser(ctx context.Context, id int64) error {
	return nil
}

func (a *Auth) DeleteUser(ctx context.Context, id int64) error {
	return nil
}

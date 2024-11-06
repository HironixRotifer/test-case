package service

import (
	"authorization/internal/models"
	"context"
)

type Auth struct {
}

// NewAuth creates a new instance of the Auth service.
func NewAuth() *Auth {
	return &Auth{}
}

func (a *Auth) Login(ctx context.Context, user models.User) (token string, err error) {
	return
}

func (a *Auth) Registr(ctx context.Context, user models.User) (userID int64, err error) {
	return
}

func (a *Auth) IsAdmin(ctx context.Context, id int64) (bool, error) {
	return false, nil
}

func (a *Auth) UpdateUser(ctx context.Context, id int64) error {
	return nil
}

func (a *Auth) DeleteUser(ctx context.Context, id int64) error {
	return nil
}

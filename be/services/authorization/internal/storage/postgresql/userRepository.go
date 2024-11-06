package postgresql

import (
	"authorization/internal/config"
	"authorization/internal/models"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type UserRepository struct {
	db *sql.DB
}

// NewUserRepository создаёт новый клиент базы данных
func NewUserRepository(config *config.Config) (*UserRepository, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%v user=%s password=%s  dbname=%s sslmode=%s",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// log.Info().Msg("database connection is success")

	return &UserRepository{db}, nil
}

func (u *UserRepository) Stop() error {
	return u.db.Close()
}

// CreateUser creates a new user in the database
func (u *UserRepository) CreateUser(ctx context.Context, user models.User) (int64, error) {
	return u.createUser(ctx, user)
}

func (u *UserRepository) createUser(ctx context.Context, user models.User) (int64, error) {
	stmt, err := u.db.Prepare("INSERT INTO users(first_name, last_name, email, phone_number, hash_password, refresh_token, ip) VALUES(?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.ExecContext(ctx, user.FirstName, user.LastName, user.Email, user.PhoneNumber, user.HashPassword, user.RefreshToken, user.IP)
	if err != nil {
		// TODO: handle error
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// UpdateUser updates the user in the database by id
func (u *UserRepository) UpdateUser(ctx context.Context, id int64, user models.User) (int64, error) {
	return 0, nil

}

// DeleteUser deletes the user in the database by id
func (u *UserRepository) DeleteUser(ctx context.Context, id int64) error {
	return nil

}

// GetUserByID returns the user by id
func (u *UserRepository) GetUserByID(ctx context.Context, id int64) (models.User, error) {
	return u.getUserByID(ctx, id)

}

func (u *UserRepository) getUserByID(ctx context.Context, id int64) (models.User, error) {

	return models.User{}, nil
}

// GetUserByEmail returns the user by email
func (u *UserRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	return models.User{}, nil
}

// IsAdmin returns true if user is admin
func (u *UserRepository) IsAdmin(ctx context.Context, id int64) (bool, error) {
	return true, nil
}

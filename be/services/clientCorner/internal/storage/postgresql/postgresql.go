package postgresql

import (
	"clientCorner/internal/config"
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

func NewUserRepository(cfg *config.Config) (*UserRepository, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &UserRepository{db: db}, nil
}

func (u *UserRepository) Stop() error {
	return u.db.Close()
}

func (u *UserRepository) UpdatePassword(ctx context.Context, password []byte, id int64) (int64, error) {
	stmt, err := u.db.Prepare("UPDATE users SET password = $1 WHERE id = $2")
	if err != nil {
		return 0, err
	}

	result, err := stmt.ExecContext(ctx, password, id)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, nil
	}

	return count, nil
}

func (u *UserRepository) UpdateEmail(ctx context.Context, email string, id int64) (int64, error) {
	stmt, err := u.db.Prepare("UPDATE users SET email = $1 WHERE id = $2")
	if err != nil {
		return 0, err
	}

	result, err := stmt.ExecContext(ctx, email, id)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (u *UserRepository) UpdateNames(ctx context.Context, firstName string, lastName string, id int64) (int64, error) {
	stmt, err := u.db.Prepare("UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3")
	if err != nil {
		return 0, err
	}

	result, err := stmt.ExecContext(ctx, firstName, lastName, id)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, err
}

func (u *UserRepository) UpdatePhoneNumber(ctx context.Context, phoneNumber string, id int64) (int64, error) {
	stmt, err := u.db.Prepare("UPDATE users SET phone_number = $1 WHERE id = $2")
	if err != nil {
		return 0, err
	}

	result, err := stmt.ExecContext(ctx, phoneNumber, id)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, nil
}

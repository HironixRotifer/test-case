package postgresql

import (
	"context"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateEmail(t *testing.T) {
	var id = int64(1)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	u := &UserRepository{db: db}

	mock.ExpectPrepare("UPDATE users SET email = \\$1 WHERE id = \\$2")
	mock.ExpectExec("UPDATE users SET email = \\$1 WHERE id = \\$2").
		WithArgs("test@example.com", id).
		WillReturnResult(sqlmock.NewResult(id, 1))

	count, err := u.UpdateEmail(context.Background(), "test@example.com", id)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

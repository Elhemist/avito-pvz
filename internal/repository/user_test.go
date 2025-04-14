package repository_test

import (
	"database/sql"
	"errors"
	"pvz-test/internal/models"
	"pvz-test/internal/repository"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestUserPostgres_GetUserById(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewUserPostgres(sqlxDB)

	t.Run("User found", func(t *testing.T) {
		userID := uuid.New()
		expectedUser := models.User{
			ID:           userID,
			Email:        "test@example.com",
			Role:         models.RoleEmployee,
			PasswordHash: "hashed_password",
			CreatedAt:    "2025-04-16T18:00:00Z",
		}

		mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1;`).
			WithArgs(userID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role", "password_hash", "created_at"}).
				AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Role, expectedUser.PasswordHash, expectedUser.CreatedAt))

		user, err := repo.GetUserById(userID)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("User not found", func(t *testing.T) {
		userID := uuid.New()

		mock.ExpectQuery(`SELECT \* FROM users WHERE id = \$1;`).
			WithArgs(userID).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserById(userID)
		assert.Error(t, err)
		assert.Equal(t, models.User{}, user)
	})
}

func TestUserPostgres_GetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewUserPostgres(sqlxDB)

	t.Run("User found", func(t *testing.T) {
		email := "test@example.com"
		expectedUser := models.User{
			ID:           uuid.New(),
			Email:        email,
			Role:         models.RoleEmployee,
			PasswordHash: "hashed_password",
			CreatedAt:    "2025-04-16T18:00:00Z",
		}

		mock.ExpectQuery(`SELECT \* FROM users WHERE email = \$1;`).
			WithArgs(email).
			WillReturnRows(sqlmock.NewRows([]string{"id", "email", "role", "password_hash", "created_at"}).
				AddRow(expectedUser.ID, expectedUser.Email, expectedUser.Role, expectedUser.PasswordHash, expectedUser.CreatedAt))

		user, err := repo.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
	})

	t.Run("User not found", func(t *testing.T) {
		email := "notfound@example.com"

		mock.ExpectQuery(`SELECT \* FROM users WHERE email = \$1;`).
			WithArgs(email).
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUserByEmail(email)
		assert.NoError(t, err)
		assert.Equal(t, models.User{}, user)
	})

	t.Run("Database error", func(t *testing.T) {
		email := "error@example.com"

		mock.ExpectQuery(`SELECT \* FROM users WHERE email = \$1;`).
			WithArgs(email).
			WillReturnError(errors.New("database error"))

		user, err := repo.GetUserByEmail(email)
		assert.Error(t, err)
		assert.Equal(t, models.User{}, user)
	})
}

func TestUserPostgres_CreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewUserPostgres(sqlxDB)

	t.Run("Successful creation", func(t *testing.T) {
		request := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "hashed_password",
			Role:     "employee",
		}
		userID := uuid.New()

		mock.ExpectQuery(`INSERT INTO users \(email, password_hash, role\) VALUES \(\$1, \$2, \$3\) RETURNING id;`).
			WithArgs(request.Email, request.Password, request.Role).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(userID))

		createdID, err := repo.CreateUser(request)
		assert.NoError(t, err)
		assert.Equal(t, userID, createdID)
	})

	t.Run("Database error", func(t *testing.T) {
		request := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "hashed_password",
			Role:     "employee",
		}

		mock.ExpectQuery(`INSERT INTO users \(email, password_hash, role\) VALUES \(\$1, \$2, \$3\) RETURNING id;`).
			WithArgs(request.Email, request.Password, request.Role).
			WillReturnError(errors.New("database error"))

		createdID, err := repo.CreateUser(request)
		assert.Error(t, err)
		assert.Equal(t, uuid.Nil, createdID)
	})
}

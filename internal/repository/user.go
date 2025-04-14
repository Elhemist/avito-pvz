package repository

import (
	"fmt"
	"pvz-test/internal/models"

	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) GetUserById(userID uuid.UUID) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1;", userTable)
	err := r.db.Get(&user, query, userID)
	if err != nil {
		return models.User{}, fmt.Errorf("no user with id: %d found: %w", userID, err)
	}
	return user, err
}

func (r *UserPostgres) GetUserByEmail(email string) (models.User, error) {
	var user models.User

	query := fmt.Sprintf("SELECT * FROM %s WHERE email = $1;", userTable)
	err := r.db.Get(&user, query, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, nil
		}
		return models.User{}, fmt.Errorf("no user with name: %s found: %w", email, err)
	}
	return user, err
}

func (r *UserPostgres) CreateUser(user models.RegisterRequest) (uuid.UUID, error) {
	var userID uuid.UUID
	query := fmt.Sprintf(`INSERT INTO %s (email, password_hash, role) VALUES ($1, $2, $3) RETURNING id;`, userTable)
	err := r.db.QueryRow(query, user.Email, user.Password, user.Role).Scan(&userID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("user create error: %w", err)
	}
	return userID, err
}

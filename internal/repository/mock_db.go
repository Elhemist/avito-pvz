package repository

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func NewMockDB() *sqlx.DB {
	db, _, err := sqlmock.New()
	if err != nil {
		panic("failed to create sqlmock: " + err.Error())
	}
	return sqlx.NewDb(db, "sqlmock")
}

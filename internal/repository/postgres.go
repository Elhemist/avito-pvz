package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

const (
	transferTable  = "transactions"
	walletTable    = "wallets"
	merchTable     = "merch"
	inventoryTable = "inventory"
	userTable      = "users"
)

type Config struct {
	Host     string
	Port     string
	Usename  string
	Password string
	DBName   string
	SSLmode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Usename, cfg.DBName, cfg.Password, cfg.SSLmode))
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

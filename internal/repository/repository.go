package repository

import (
	"pvz-test/internal/models"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user models.RegisterRequest) (uuid.UUID, error)
	GetUserByEmail(email string) (models.User, error)
	GetUserById(userID uuid.UUID) (models.User, error)
}

type PvzRepository interface {
	CreatePvz(city string) (models.PVZ, error)
	Exists(pvzID uuid.UUID) (bool, error)
	GetPVZList(limit, offset int) ([]models.PVZ, error)
}

type ReceptionRepository interface {
	AddItem(pvzID uuid.UUID, itemType string) (models.Item, error)
	DeleteItem(pvzID uuid.UUID) error
	CreateReception(pvzID uuid.UUID) (models.Reception, error)
	GetActiveReception(pvzID uuid.UUID) (models.Reception, error)
	CloseReception(receptionID uuid.UUID) error
	GetReceptionsWithProducts(pvzID uuid.UUID, start, end *time.Time) ([]models.Reception, error)
	GetItemsByReceptionID(receptionID uuid.UUID) ([]models.Item, error)
}

type Repository struct {
	UserRepository
	PvzRepository
	ReceptionRepository
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepository:      NewUserPostgres(db),
		PvzRepository:       NewPvzPostgres(db),
		ReceptionRepository: NewReceptionPostgres(db),
	}
}

package service

import (
	"pvz-test/internal/models"
	"pvz-test/internal/repository"
	"time"

	"github.com/google/uuid"
)

type Authorization interface {
	Register(user models.RegisterRequest) (models.UserResponse, error)
	Login(user models.LoginRequest) (string, error)
	DummyLogin(role models.Role) (string, error)
	ParseToken(token string) (uuid.UUID, models.Role, error)
}

type Reception interface {
	CreateReception(pvzID uuid.UUID) (models.Reception, error)
	CloseActiveReception(pvzID uuid.UUID) (models.Reception, error)
	DeleteItem(pvzID uuid.UUID) error
	AddItem(pvzID uuid.UUID, itemType string) (models.Item, error)
}

type Pvz interface {
	CreatePvz(city string) (models.PVZ, error)
	GetFilteredPVZ(start, end *time.Time, limit, offset int) ([]models.PVZResponse, error)
}

type Service struct {
	Authorization
	Reception
	Pvz
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.UserRepository),
		Reception:     NewReceptionService(repos.ReceptionRepository, repos.PvzRepository),
		Pvz:           NewPvzService(repos.PvzRepository, repos.ReceptionRepository),
	}
}

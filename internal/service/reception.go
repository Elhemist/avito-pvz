package service

import (
	"fmt"
	"pvz-test/internal/models"
	"pvz-test/internal/repository"

	"github.com/google/uuid"
)

type ReceptionService struct {
	receptionRepo repository.ReceptionRepository
	pvzRepo       repository.PvzRepository
}

func NewReceptionService(receptionRepo repository.ReceptionRepository, pvzRepo repository.PvzRepository) *ReceptionService {
	return &ReceptionService{
		receptionRepo: receptionRepo,
		pvzRepo:       pvzRepo}
}

func (s *ReceptionService) CreateReception(pvzID uuid.UUID) (models.Reception, error) {
	exists, err := s.pvzRepo.Exists(pvzID)
	if err != nil {
		return models.Reception{}, fmt.Errorf("failed to check PVZ existence: %s", err.Error())
	}
	if !exists {
		return models.Reception{}, fmt.Errorf("PVZ: %s does not exist", pvzID.String())
	}

	activeReception, err := s.receptionRepo.GetActiveReception(pvzID)
	if err != nil {
		return models.Reception{}, fmt.Errorf("reception get error: %s", pvzID.String())
	}
	if (activeReception != models.Reception{}) {
		return models.Reception{}, fmt.Errorf("an active reception already exists for PVZ: %s", pvzID.String())
	}

	reception, err := s.receptionRepo.CreateReception(pvzID)
	if err != nil {
		return models.Reception{}, fmt.Errorf("failed to create reception: %s", err.Error())
	}

	return reception, nil
}

func (s *ReceptionService) CloseActiveReception(pvzID uuid.UUID) (models.Reception, error) {
	reception, err := s.receptionRepo.GetActiveReception(pvzID)
	if err != nil {
		return models.Reception{}, err
	}
	if (reception == models.Reception{}) {
		return models.Reception{}, nil
	}

	err = s.receptionRepo.CloseReception(reception.ID)
	if err != nil {
		return models.Reception{}, err
	}
	reception.Status = "close"
	return reception, nil
}

func (s *ReceptionService) AddItem(pvzID uuid.UUID, itemType string) (models.Item, error) {
	item, err := s.receptionRepo.AddItem(pvzID, itemType)
	return item, err
}

func (s *ReceptionService) DeleteItem(pvzID uuid.UUID) error {
	err := s.receptionRepo.DeleteItem(pvzID)
	return err
}

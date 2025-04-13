package service

import (
	"fmt"
	"pvz-test/internal/models"
	"pvz-test/internal/repository"
	"time"
)

type PvzService struct {
	pvzRepo       repository.PvzRepository
	receptionRepo repository.ReceptionRepository
}

func NewPvzService(pvzRepo repository.PvzRepository, receptionRepo repository.ReceptionRepository) *PvzService {
	return &PvzService{pvzRepo: pvzRepo, receptionRepo: receptionRepo}
}

var allowedCities = map[string]struct{}{
	"Москва":          {},
	"Санкт-Петербург": {},
	"Казань":          {},
}

func (s *PvzService) CreatePvz(city string) (models.PVZ, error) {

	if _, ok := allowedCities[city]; !ok {
		return models.PVZ{}, fmt.Errorf("city %s city is not supported", city)
	}
	pvz, err := s.pvzRepo.CreatePvz(city)
	if err != nil {
		return models.PVZ{}, fmt.Errorf("pvz create error: %s", err.Error())
	}

	return pvz, nil
}

func (s *PvzService) GetFilteredPVZ(start, end *time.Time, limit, offset int) ([]models.PVZResponse, error) {
	pvzs, err := s.pvzRepo.GetPVZList(limit, offset)
	if err != nil {
		return nil, err
	}

	var result []models.PVZResponse
	for _, pvz := range pvzs {
		receptions, err := s.receptionRepo.GetReceptionsWithProducts(pvz.ID, start, end)
		if err != nil {
			return nil, err
		}

		var blocks []models.ReceptionBlock
		for _, r := range receptions {
			items, err := s.receptionRepo.GetItemsByReceptionID(r.ID)
			if err != nil {
				return nil, err
			}
			blocks = append(blocks, models.ReceptionBlock{
				Reception: r,
				Products:  items,
			})
		}

		result = append(result, models.PVZResponse{
			PVZ:        pvz,
			Receptions: blocks,
		})
	}

	return result, nil
}

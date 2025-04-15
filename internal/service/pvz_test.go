package service_test

import (
	"errors"
	"pvz-test/internal/models"
	"pvz-test/internal/service"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPvzRepository struct {
	mock.Mock
}

func (m *MockPvzRepository) CreatePvz(city string) (models.PVZ, error) {
	args := m.Called(city)
	return args.Get(0).(models.PVZ), args.Error(1)
}

func (m *MockPvzRepository) GetPVZList(limit, offset int) ([]models.PVZ, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.PVZ), args.Error(1)
}

func (m *MockReceptionRepository) GetReceptionsWithProducts(pvzID uuid.UUID, start, end *time.Time) ([]models.Reception, error) {
	args := m.Called(pvzID, start, end)
	return args.Get(0).([]models.Reception), args.Error(1)
}

func (m *MockReceptionRepository) GetItemsByReceptionID(receptionID uuid.UUID) ([]models.Item, error) {
	args := m.Called(receptionID)
	return args.Get(0).([]models.Item), args.Error(1)
}

func TestPvzService_CreatePvz(t *testing.T) {
	mockPvzRepo := new(MockPvzRepository)
	mockReceptionRepo := new(MockReceptionRepository)
	service := service.NewPvzService(mockPvzRepo, mockReceptionRepo)

	t.Run("Invalid city", func(t *testing.T) {
		_, err := service.CreatePvz("InvalidCity")
		assert.EqualError(t, err, "city InvalidCity city is not supported")
	})

	t.Run("Valid city", func(t *testing.T) {
		expectedPvz := models.PVZ{
			ID:               uuid.New(),
			RegistrationDate: time.Now(),
			City:             "Москва",
		}
		mockPvzRepo.On("CreatePvz", "Москва").Return(expectedPvz, nil)

		pvz, err := service.CreatePvz("Москва")
		assert.NoError(t, err)
		assert.Equal(t, expectedPvz, pvz)
		mockPvzRepo.AssertExpectations(t)
	})
}

func TestPvzService_GetFilteredPVZ_Bad(t *testing.T) {
	mockPvzRepo := new(MockPvzRepository)
	mockReceptionRepo := new(MockReceptionRepository)
	service := service.NewPvzService(mockPvzRepo, mockReceptionRepo)

	t.Run("Error fetching PVZ list", func(t *testing.T) {
		mockPvzRepo.On("GetPVZList", 10, 0).Return([]models.PVZ{}, errors.New("database error"))

		_, err := service.GetFilteredPVZ(nil, nil, 10, 0)
		assert.EqualError(t, err, "database error")
		mockPvzRepo.AssertExpectations(t)
	})

}

func TestPvzService_GetFilteredPVZ_Good(t *testing.T) {
	mockPvzRepo := new(MockPvzRepository)
	mockReceptionRepo := new(MockReceptionRepository)
	service := service.NewPvzService(mockPvzRepo, mockReceptionRepo)

	t.Run("Successful fetch with receptions and items", func(t *testing.T) {
		pvzID := uuid.New()
		receptionID := uuid.New()
		expectedPVZ := []models.PVZ{
			{ID: pvzID, RegistrationDate: time.Now(), City: "Москва"},
		}
		expectedReceptions := []models.Reception{
			{ID: receptionID, PVZID: pvzID, Status: "closed", CreatedAt: time.Now()},
		}
		expectedItems := []models.Item{
			{ID: uuid.New(), ReceptionID: receptionID, PvzID: pvzID, Type: models.ItemTypeElectronics, AddedAt: time.Now()},
		}

		mockPvzRepo.On("GetPVZList", 10, 0).Return(expectedPVZ, nil).Once()
		mockReceptionRepo.On("GetReceptionsWithProducts", pvzID, (*time.Time)(nil), (*time.Time)(nil)).Return(expectedReceptions, nil).Once()
		mockReceptionRepo.On("GetItemsByReceptionID", receptionID).Return(expectedItems, nil).Once()

		result, err := service.GetFilteredPVZ(nil, nil, 10, 0)

		assert.NoError(t, err, "Expected no error, but got one")
		assert.Len(t, result, 1, "Expected 1 PVZ in the result")
		assert.Equal(t, expectedPVZ[0], result[0].PVZ, "Expected PVZ does not match")
		assert.Len(t, result[0].Receptions, 1, "Expected 1 reception in the result")
		assert.Equal(t, expectedReceptions[0], result[0].Receptions[0].Reception, "Expected reception does not match")
		assert.Equal(t, expectedItems, result[0].Receptions[0].Products, "Expected items do not match")

		mockPvzRepo.AssertExpectations(t)
		mockReceptionRepo.AssertExpectations(t)
	})
}

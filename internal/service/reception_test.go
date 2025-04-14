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

type MockReceptionRepository struct {
	mock.Mock
}

func (m *MockReceptionRepository) CreateReception(pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *MockReceptionRepository) GetActiveReception(pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *MockReceptionRepository) CloseReception(receptionID uuid.UUID) error {
	args := m.Called(receptionID)
	return args.Error(0)
}

func (m *MockReceptionRepository) AddItem(pvzID uuid.UUID, itemType string) (models.Item, error) {
	args := m.Called(pvzID, itemType)
	return args.Get(0).(models.Item), args.Error(1)
}

func (m *MockReceptionRepository) DeleteItem(pvzID uuid.UUID) error {
	args := m.Called(pvzID)
	return args.Error(0)
}

func (m *MockPvzRepository) Exists(pvzID uuid.UUID) (bool, error) {
	args := m.Called(pvzID)
	return args.Bool(0), args.Error(1)
}

func TestReceptionService_CreateReception(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepository)
	mockPvzRepo := new(MockPvzRepository)
	service := service.NewReceptionService(mockReceptionRepo, mockPvzRepo)

	t.Run("Non-existent PVZ", func(t *testing.T) {
		pvzID := uuid.New()
		mockPvzRepo.On("Exists", pvzID).Return(false, nil)

		_, err := service.CreateReception(pvzID)
		assert.EqualError(t, err, "PVZ: "+pvzID.String()+" does not exist")
		mockPvzRepo.AssertExpectations(t)
	})
}

func TestReceptionService_CloseActiveReception(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepository)
	mockPvzRepo := new(MockPvzRepository)
	service := service.NewReceptionService(mockReceptionRepo, mockPvzRepo)

	t.Run("Error fetching active reception", func(t *testing.T) {
		pvzID := uuid.New()
		mockReceptionRepo.On("GetActiveReception", pvzID).Return(models.Reception{}, errors.New("database error"))

		_, err := service.CloseActiveReception(pvzID)
		assert.EqualError(t, err, "database error")
		mockReceptionRepo.AssertExpectations(t)
	})

	t.Run("No active reception", func(t *testing.T) {
		pvzID := uuid.New()
		mockReceptionRepo.On("GetActiveReception", pvzID).Return(models.Reception{}, nil)

		reception, err := service.CloseActiveReception(pvzID)
		assert.NoError(t, err)
		assert.Equal(t, models.Reception{}, reception)
		mockReceptionRepo.AssertExpectations(t)
	})

	t.Run("Successful closure", func(t *testing.T) {
		pvzID := uuid.New()
		receptionID := uuid.New()
		activeReception := models.Reception{ID: receptionID, PVZID: pvzID, Status: "in_progress"}
		mockReceptionRepo.On("GetActiveReception", pvzID).Return(activeReception, nil)
		mockReceptionRepo.On("CloseReception", receptionID).Return(nil)

		reception, err := service.CloseActiveReception(pvzID)
		assert.NoError(t, err)
		assert.Equal(t, "close", reception.Status)
		mockReceptionRepo.AssertExpectations(t)
	})
}

func TestReceptionService_AddItem(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepository)
	mockPvzRepo := new(MockPvzRepository)
	service := service.NewReceptionService(mockReceptionRepo, mockPvzRepo)

	t.Run("Error adding item", func(t *testing.T) {
		pvzID := uuid.New()
		itemType := "electronics"
		mockReceptionRepo.On("AddItem", pvzID, itemType).Return(models.Item{}, errors.New("database error"))

		_, err := service.AddItem(pvzID, itemType)
		assert.EqualError(t, err, "database error")
		mockReceptionRepo.AssertExpectations(t)
	})

	t.Run("Successful item addition", func(t *testing.T) {
		pvzID := uuid.New()
		itemType := "electronics"
		expectedItem := models.Item{ID: 1, ReceptionID: uuid.New(), PvzID: pvzID, Type: models.ItemTypeElectronics, AddedAt: time.Now()}
		mockReceptionRepo.On("AddItem", pvzID, itemType).Return(expectedItem, nil)

		item, err := service.AddItem(pvzID, itemType)
		assert.NoError(t, err)
		assert.Equal(t, expectedItem, item)
		mockReceptionRepo.AssertExpectations(t)
	})
}

func TestReceptionService_DeleteItem(t *testing.T) {
	mockReceptionRepo := new(MockReceptionRepository)
	mockPvzRepo := new(MockPvzRepository)
	service := service.NewReceptionService(mockReceptionRepo, mockPvzRepo)

	t.Run("Error deleting item", func(t *testing.T) {
		pvzID := uuid.New()
		mockReceptionRepo.On("DeleteItem", pvzID).Return(errors.New("database error"))

		err := service.DeleteItem(pvzID)
		assert.EqualError(t, err, "database error")
		mockReceptionRepo.AssertExpectations(t)
	})

	t.Run("Successful item deletion", func(t *testing.T) {
		pvzID := uuid.New()
		mockReceptionRepo.On("DeleteItem", pvzID).Return(nil)

		err := service.DeleteItem(pvzID)
		assert.NoError(t, err)
		mockReceptionRepo.AssertExpectations(t)
	})
}

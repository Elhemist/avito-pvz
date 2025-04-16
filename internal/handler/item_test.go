package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"pvz-test/internal/handler"
	"pvz-test/internal/models"
	"pvz-test/internal/service"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockReceptionService struct {
	mock.Mock
}

func (m *MockReceptionService) DeleteItem(pvzID uuid.UUID) error {
	args := m.Called(pvzID)
	return args.Error(0)
}

func (m *MockReceptionService) AddItem(pvzID uuid.UUID, itemType string) (models.Item, error) {
	args := m.Called(pvzID, itemType)
	return args.Get(0).(models.Item), args.Error(1)
}

func (m *MockReceptionService) CloseActiveReception(pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *MockReceptionService) CreateReception(pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func TestHandler_RemoveLastItem(t *testing.T) {
	mockService := new(MockReceptionService)
	h := handler.NewHandler(&service.Service{Reception: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/pvz/:pvzId/delete_last_product", func(c *gin.Context) {
		c.Set("role", models.RoleEmployee)
		h.RemoveLastItem(c)
	})

	t.Run("Successful removal", func(t *testing.T) {
		pvzID := uuid.New()
		mockService.On("DeleteItem", pvzID).Return(nil)

		req, _ := http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product?pvz_id="+pvzID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid PVZ ID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/pvz/invalid-uuid/delete_last_product", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error during removal", func(t *testing.T) {
		pvzID := uuid.New()
		mockService.On("DeleteItem", pvzID).Return(assert.AnError)

		req, _ := http.NewRequest(http.MethodPost, "/pvz/"+pvzID.String()+"/delete_last_product?pvz_id="+pvzID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestHandler_AddItem(t *testing.T) {
	mockService := new(MockReceptionService)
	h := handler.NewHandler(&service.Service{Reception: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/products", func(c *gin.Context) {
		c.Set("role", models.RoleEmployee)
		h.AddItem(c)
	})

	t.Run("Successful addition", func(t *testing.T) {
		pvzID := uuid.New()
		itemType := "electronics"
		expectedItem := models.Item{ID: uuid.New(), ReceptionID: uuid.New(), Type: models.ItemTypeElectronics}
		mockService.On("AddItem", pvzID, itemType).Return(expectedItem, nil)

		body, _ := json.Marshal(models.AddProductRequest{PvzID: pvzID, Type: itemType})
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid request body", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Error during addition", func(t *testing.T) {
		pvzID := uuid.New()
		itemType := "electronics"
		mockService.On("AddItem", pvzID, itemType).Return(models.Item{}, assert.AnError)

		body, _ := json.Marshal(models.AddProductRequest{PvzID: pvzID, Type: itemType})
		req, _ := http.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})
}

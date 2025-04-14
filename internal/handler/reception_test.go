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

const roleCtx = "role" // Добавлено определение roleCtx

type MockService struct {
	mock.Mock
}

func (m *MockService) CreateReception(pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *MockService) CloseActiveReception(pvzID uuid.UUID) (models.Reception, error) {
	args := m.Called(pvzID)
	return args.Get(0).(models.Reception), args.Error(1)
}

func (m *MockService) AddItem(pvzID uuid.UUID, itemType string) (models.Item, error) {
	args := m.Called(pvzID, itemType)
	return args.Get(0).(models.Item), args.Error(1)
}

func (m *MockService) DeleteItem(pvzID uuid.UUID) error {
	args := m.Called(pvzID)
	return args.Error(0)
}

func TestHandler_CreateReception(t *testing.T) {
	mockService := new(MockService)
	h := handler.NewHandler(&service.Service{Reception: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/reception", func(c *gin.Context) {
		// Устанавливаем роль в контексте
		c.Set(roleCtx, models.RoleEmployee)
		h.CreateReception(c)
	})

	t.Run("Successful creation", func(t *testing.T) {
		pvzID := uuid.New()
		expectedReception := models.Reception{ID: uuid.New(), PVZID: pvzID, Status: "created"}
		mockService.On("CreateReception", pvzID).Return(expectedReception, nil)

		body, _ := json.Marshal(pvzID)
		req, _ := http.NewRequest(http.MethodPost, "/reception", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Binding error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/reception", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_CloseReception(t *testing.T) {
	mockService := new(MockService)
	h := handler.NewHandler(&service.Service{Reception: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.GET("/reception/close", func(c *gin.Context) {
		c.Set(roleCtx, models.RoleEmployee)
		h.CloseReception(c)
	})

	t.Run("Successful closure", func(t *testing.T) {
		pvzID := uuid.New()
		expectedReception := models.Reception{ID: uuid.New(), PVZID: pvzID, Status: "closed"}
		mockService.On("CloseActiveReception", pvzID).Return(expectedReception, nil)

		req, _ := http.NewRequest(http.MethodGet, "/reception/close?pvz_id="+pvzID.String(), nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid UUID", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodGet, "/reception/close?pvz_id=invalid-uuid", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPvzService struct {
	mock.Mock
}

func (m *MockPvzService) CreatePvz(city string) (models.PVZ, error) {
	args := m.Called(city)
	return args.Get(0).(models.PVZ), args.Error(1)
}

func (m *MockPvzService) GetFilteredPVZ(start, end *time.Time, limit, offset int) ([]models.PVZResponse, error) {
	args := m.Called(start, end, limit, offset)
	return args.Get(0).([]models.PVZResponse), args.Error(1)
}

func TestHandler_CreatePVZ(t *testing.T) {
	mockService := new(MockPvzService)
	h := handler.NewHandler(&service.Service{Pvz: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.POST("/pvz", func(c *gin.Context) {
		c.Set("role", models.RoleModerator)
		h.CreatePVZ(c)
	})

	t.Run("Invalid city", func(t *testing.T) {
		body, _ := json.Marshal(models.PVZRequest{City: "InvalidCity"})
		req, _ := http.NewRequest(http.MethodPost, "/pvz", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
func TestHandler_GetPVZList(t *testing.T) {
	mockService := new(MockPvzService)
	h := handler.NewHandler(&service.Service{Pvz: mockService})

	t.Run("Unauthorized user", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/pvz", func(c *gin.Context) {
			h.GetPVZList(c) // Не устанавливаем пользователя
		})

		req, _ := http.NewRequest(http.MethodGet, "/pvz", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
	})

	t.Run("Invalid query parameters", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		router := gin.New()
		router.GET("/pvz", func(c *gin.Context) {
			c.Set("user", models.User{Role: models.RoleEmployee}) // Устанавливаем роль сотрудника
			h.GetPVZList(c)
		})

		req, _ := http.NewRequest(http.MethodGet, "/pvz?startDate=invalid-date", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

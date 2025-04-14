package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
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

type MockAuthorizationService struct {
	mock.Mock
}

func (m *MockAuthorizationService) Login(input models.LoginRequest) (string, error) {
	args := m.Called(input)
	return args.String(0), args.Error(1)
}

func (m *MockAuthorizationService) Register(input models.RegisterRequest) (models.UserResponse, error) {
	args := m.Called(input)
	return args.Get(0).(models.UserResponse), args.Error(1)
}

func (m *MockAuthorizationService) DummyLogin(role models.Role) (string, error) {
	args := m.Called(role)
	return args.String(0), args.Error(1)
}

func (m *MockAuthorizationService) ParseToken(token string) (uuid.UUID, models.Role, error) {
	args := m.Called(token)
	return args.Get(0).(uuid.UUID), args.Get(1).(models.Role), args.Error(2)
}

func TestHandler_DummyLogin(t *testing.T) {
	mockService := new(MockAuthorizationService)
	h := handler.NewHandler(&service.Service{Authorization: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/dummy-login", h.DummyLogin)

	t.Run("Successful dummy login", func(t *testing.T) {
		role := models.RoleModerator
		token := "mockToken"

		mockService.On("DummyLogin", role).Return(token, nil)

		body, _ := json.Marshal(models.DummyLoginRequest{Role: role})
		req, _ := http.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `"mockToken"`, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Binding error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message":"binding error"}`, w.Body.String())
	})

	t.Run("Validation error", func(t *testing.T) {
		body, _ := json.Marshal(models.DummyLoginRequest{Role: ""})
		req, _ := http.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message":"role validate error"}`, w.Body.String())
	})

	t.Run("Service error", func(t *testing.T) {
		role := models.RoleEmployee

		mockService.On("DummyLogin", role).Return("", errors.New("service error"))

		body, _ := json.Marshal(models.DummyLoginRequest{Role: role})
		req, _ := http.NewRequest(http.MethodPost, "/dummy-login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, `{"message":"service error"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})
}
func TestHandler_Register_Bad(t *testing.T) {
	mockService := new(MockAuthorizationService)
	h := handler.NewHandler(&service.Service{Authorization: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/register", h.Register)

	t.Run("Binding error", func(t *testing.T) {
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Исправлено на 400
	})

	t.Run("Validation error", func(t *testing.T) {
		input := models.RegisterRequest{
			Email:    "",
			Password: "",
			Role:     "",
		}

		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code) // Исправлено на 400
	})

	t.Run("Service error", func(t *testing.T) {
		input := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Role:     string(models.RoleEmployee),
		}

		mockService.On("Register", input).Return(models.UserResponse{}, errors.New("service error"))

		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code) // Исправлено на 500
		assert.JSONEq(t, `{"message":"service error"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestHandler_Register_Good(t *testing.T) {
	mockService := new(MockAuthorizationService)
	h := handler.NewHandler(&service.Service{Authorization: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/register", h.Register)

	t.Run("Successful registration", func(t *testing.T) {
		input := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Role:     string(models.RoleEmployee),
		}
		expectedUser := models.UserResponse{
			ID:    uuid.New(),
			Email: input.Email,
			Role:  input.Role,
		}

		mockService.On("Register", input).Return(expectedUser, nil)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		expectedResponse, _ := json.Marshal(expectedUser)
		assert.JSONEq(t, string(expectedResponse), w.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestHandler_Login_Good(t *testing.T) {
	mockService := new(MockAuthorizationService)
	h := handler.NewHandler(&service.Service{Authorization: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", h.Login)

	t.Run("Successful login", func(t *testing.T) {
		input := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		token := "mockToken"

		mockService.On("Login", input).Return(token, nil)

		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, `"mockToken"`, w.Body.String())
		mockService.AssertExpectations(t)
	})

}
func TestHandler_Login_Bad(t *testing.T) {
	mockService := new(MockAuthorizationService)
	h := handler.NewHandler(&service.Service{Authorization: mockService})

	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.POST("/login", h.Login)
	t.Run("Service error", func(t *testing.T) {
		input := models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		}

		mockService.On("Login", input).Return("", errors.New("unauthorized"))

		body, _ := json.Marshal(input)
		req, _ := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		assert.JSONEq(t, `{"message":"Unauthorized"}`, w.Body.String())
		mockService.AssertExpectations(t)
	})
}

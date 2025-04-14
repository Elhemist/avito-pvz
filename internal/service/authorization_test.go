package service_test

import (
	"errors"
	"pvz-test/internal/models"
	"pvz-test/internal/service"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

const (
	salt       = "someSalt"
	signingKey = "podpis"
	tokenTTL   = time.Hour / 2
)

func (m *MockUserRepository) CreateUser(user models.RegisterRequest) (uuid.UUID, error) {
	args := m.Called(user)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (models.User, error) {
	args := m.Called(email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(userID uuid.UUID) (models.User, error) {
	args := m.Called(userID)
	return args.Get(0).(models.User), args.Error(1)
}

func TestAuthorizationService_Login_Bad(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	t.Run("User not found", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "test@example.com").Return(models.User{}, errors.New("Unauthorized"))

		token, err := authService.Login(models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		assert.Empty(t, token)
		assert.EqualError(t, err, "Unauthorized")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Incorrect password", func(t *testing.T) {
		mockRepo.On("GetUserByEmail", "test@example.com").Return(models.User{
			ID:           uuid.New(),
			Email:        "test@example.com",
			PasswordHash: service.GeneratePasswordHash("password"),
		}, errors.New("user not found"))

		token, err := authService.Login(models.LoginRequest{
			Email:    "test@example.com",
			Password: "wrong_password",
		})

		assert.Empty(t, token)
		assert.EqualError(t, err, "Unauthorized")
		mockRepo.AssertExpectations(t)
	})

	mockRepo = new(MockUserRepository)

}

func TestAuthorizationService_Login_Good(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)
	t.Run("Successful login", func(t *testing.T) {
		userID := uuid.New()
		mockRepo.On("GetUserByEmail", "test@example.com").Return(models.User{
			ID:           userID,
			Email:        "test@example.com",
			PasswordHash: service.GeneratePasswordHash("password123"),
			Role:         models.RoleEmployee,
		}, nil)

		token, err := authService.Login(models.LoginRequest{
			Email:    "test@example.com",
			Password: "password123",
		})

		assert.NoError(t, err)
		assert.NotEmpty(t, token)

		// Validate the token
		parsedToken, err := jwt.ParseWithClaims(token, &service.TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})
		assert.NoError(t, err)

		claims, ok := parsedToken.Claims.(*service.TokenClaims)
		assert.True(t, ok)
		assert.Equal(t, userID, claims.UserId)
		assert.Equal(t, models.RoleEmployee, claims.Role)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthorizationService_DummyLogin(t *testing.T) {
	authService := service.NewAuthService(nil)

	role := models.RoleEmployee
	token, err := authService.DummyLogin(role)
	assert.NoError(t, err)

	userID, parsedRole, err := authService.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, uuid.Max, userID)
	assert.Equal(t, role, parsedRole)
}

func TestAuthorizationService_ParseToken(t *testing.T) {
	authService := service.NewAuthService(nil)

	role := models.RoleModerator
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &service.TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		uuid.New(),
		role,
	})
	tokenString, _ := token.SignedString([]byte(signingKey))

	userID, parsedRole, err := authService.ParseToken(tokenString)
	assert.NoError(t, err)
	assert.Equal(t, role, parsedRole)
	assert.NotEqual(t, uuid.Nil, userID)
}

func TestAuthorizationService_ParseToken_InvalidToken(t *testing.T) {
	authService := service.NewAuthService(nil)

	_, _, err := authService.ParseToken("invalid_token")
	assert.Error(t, err)
}

func TestGeneratePasswordHash(t *testing.T) {
	password := "password"
	expectedHash := service.GeneratePasswordHash(password)

	assert.Equal(t, service.GeneratePasswordHash(password), expectedHash)
	assert.NotEqual(t, service.GeneratePasswordHash("different_password"), expectedHash)
}

func TestAuthorizationService_Register_Good(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	t.Run("Successful registration", func(t *testing.T) {
		request := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Role:     "employee",
		}
		hashedPassword := service.GeneratePasswordHash(request.Password)
		expectedRequest := models.RegisterRequest{
			Email:    request.Email,
			Password: hashedPassword,
			Role:     request.Role,
		}
		expectedUser := models.UserResponse{
			ID:    uuid.New().String(),
			Email: request.Email,
			Role:  request.Role,
		}
		mockRepo.On("CreateUser", expectedRequest).Return(uuid.MustParse(expectedUser.ID), nil)

		user, err := authService.Register(request)
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})
}

func TestAuthorizationService_Register_Bad(t *testing.T) {
	mockRepo := new(MockUserRepository)
	authService := service.NewAuthService(mockRepo)

	t.Run("Error during user creation", func(t *testing.T) {
		request := models.RegisterRequest{
			Email:    "test@example.com",
			Password: "password123",
			Role:     "employee",
		}
		hashedPassword := service.GeneratePasswordHash(request.Password)
		expectedRequest := models.RegisterRequest{
			Email:    request.Email,
			Password: hashedPassword,
			Role:     request.Role,
		}
		mockRepo.On("CreateUser", expectedRequest).Return(uuid.Nil, errors.New("database error"))

		_, err := authService.Register(request)
		assert.EqualError(t, err, "database error")
		mockRepo.AssertExpectations(t)
	})
}

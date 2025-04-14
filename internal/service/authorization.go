package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"pvz-test/internal/models"
	"pvz-test/internal/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const NEW_USER_BALANCE = 1000
const (
	salt       = "someSalt"
	signingKey = "podpis"
	tokenTTL   = time.Hour / 2
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID   `json:"user_id"`
	Role   models.Role `json:"role"`
}

type AuthorizationService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) *AuthorizationService {
	return &AuthorizationService{userRepo: userRepo}
}

func (s *AuthorizationService) Register(user models.RegisterRequest) (models.UserResponse, error) {
	user.Password = GeneratePasswordHash(user.Password)
	UserID, err := s.userRepo.CreateUser(user)
	if err != nil {
		logrus.Info(err)
		return models.UserResponse{}, err
	}
	response := models.UserResponse{ID: UserID.String(), Email: user.Email, Role: user.Role}
	return response, err
}

func (s *AuthorizationService) Login(userReq models.LoginRequest) (string, error) {
	user, err := s.userRepo.GetUserByEmail(userReq.Email)
	if err != nil {
		logrus.Info(err)
		return "", err
	}
	if user == (models.User{}) {
		return "", fmt.Errorf("Unauthorized")
	} else if user.PasswordHash != GeneratePasswordHash(userReq.Password) {
		return "", fmt.Errorf("Unauthorized")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.ID,
		user.Role,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthorizationService) DummyLogin(role models.Role) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		uuid.Max,
		role,
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthorizationService) ParseToken(accessToken string) (uuid.UUID, models.Role, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(signingKey), nil
	})
	if err != nil {
		return uuid.Nil, "", err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return uuid.Nil, "", errors.New("invalid token struct")
	}

	return claims.UserId, claims.Role, nil
}

func GeneratePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

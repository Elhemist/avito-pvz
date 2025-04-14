package service_test

import (
	"pvz-test/internal/repository"
	"pvz-test/internal/service"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	mockUserRepo := new(MockUserRepository)
	mockPvzRepo := new(MockPvzRepository)
	mockReceptionRepo := new(MockReceptionRepository)

	repos := &repository.Repository{
		UserRepository:      mockUserRepo,
		PvzRepository:       mockPvzRepo,
		ReceptionRepository: mockReceptionRepo,
	}

	svc := service.NewService(repos)

	assert.NotNil(t, svc.Authorization)
	assert.NotNil(t, svc.Reception)
	assert.NotNil(t, svc.Pvz)
}

package models

import (
	"time"

	"github.com/google/uuid"
)

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
	Role     string `json:"role" validate:"required,oneof=employee moderator client"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID    uuid.UUID `json:"id"`
	Email string    `json:"email"`
	Role  string    `json:"role"`
}

type DummyLoginRequest struct {
	Role Role `json:"role" validate:"required,oneof=employee moderator client"`
}

type PVZRequest struct {
	City string `json:"city" validate:"required,oneof=Москва Санкт-Петербург Казань"`
}

type AddProductRequest struct {
	Type  string    `json:"type" binding:"required"`
	PvzID uuid.UUID `json:"pvzId" binding:"required"`
}

type GetPVZListQuery struct {
	StartDate *time.Time `form:"startDate"`
	EndDate   *time.Time `form:"endDate"`
	Page      int        `form:"page"`
	Limit     int        `form:"limit"`
}

type PVZResponse struct {
	PVZ        PVZ              `json:"pvz"`
	Receptions []ReceptionBlock `json:"receptions"`
}

type CreateReceptionRequest struct {
	PvzID uuid.UUID `json:"pvzId" binding:"required"`
}

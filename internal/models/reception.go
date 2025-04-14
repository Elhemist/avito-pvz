package models

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PVZID     uuid.UUID `json:"pvz_id" db:"pvz_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type ReceptionBlock struct {
	Reception Reception `json:"reception"`
	Products  []Item    `json:"products"`
}

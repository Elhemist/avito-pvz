package models

import (
	"time"

	"github.com/google/uuid"
)

type Reception struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PVZID     uuid.UUID `json:"pvzId" db:"pvz_id"`
	Status    string    `json:"status" db:"status"`
	CreatedAt time.Time `json:"dateTime" db:"created_at"`
}

type ReceptionBlock struct {
	Reception Reception `json:"reception"`
	Products  []Item    `json:"products"`
}

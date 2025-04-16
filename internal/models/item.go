package models

import (
	"time"

	"github.com/google/uuid"
)

type ItemType string

const (
	ItemTypeElectronics ItemType = "electronics"
	ItemTypeClothing    ItemType = "clothing"
	ItemTypeShoes       ItemType = "shoes"
)

type Item struct {
	ID          uuid.UUID `db:"id" json:"id"`
	ReceptionID uuid.UUID `db:"reception_id" json:"receptionId"`
	Type        ItemType  `db:"type" json:"type"`
	AddedAt     time.Time `db:"added_at" json:"dateTime"`
}

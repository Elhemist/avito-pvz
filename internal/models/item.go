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
	ID          int       `db:"id"`
	ReceptionID int       `db:"reception_id"`
	PvzID       uuid.UUID `db:"pvz_id"`
	Type        ItemType  `db:"type"`
	AddedAt     time.Time `db:"added_at"`
}

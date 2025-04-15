package repository

import (
	"fmt"
	"pvz-test/internal/models"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type PvzPostgres struct {
	db *sqlx.DB
}

func NewPvzPostgres(db *sqlx.DB) *PvzPostgres {
	return &PvzPostgres{db: db}
}

func (r *PvzPostgres) CreatePvz(city string) (models.PVZ, error) {
	var pvz models.PVZ
	logrus.Infof("Inserting new PVZ with city: %s", city)
	err := r.db.Get(&pvz, `
        INSERT INTO pvz (city)
        VALUES ($1)
        RETURNING id, city, registration_date
    `, city)
	if err != nil {
		logrus.Errorf("Error inserting PVZ: %v", err)
		return models.PVZ{}, err
	}
	logrus.Infof("Inserted new PVZ: %v", pvz)
	return pvz, nil
}

func (r *PvzPostgres) Exists(pvzID uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM pvz WHERE id = $1
		)
	`, pvzID)
	if err != nil {
		return false, fmt.Errorf("failed to check PVZ existence: %w", err)
	}
	logrus.Infof("PVZ existence check for ID %s: %v", pvzID, exists)
	return exists, nil
}

func (r *PvzPostgres) GetPVZList(limit, offset int) ([]models.PVZ, error) {
	var pvzs []models.PVZ
	err := r.db.Select(&pvzs, `
		SELECT id, registration_date, city
		FROM pvz
		ORDER BY registration_date DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	return pvzs, err
}

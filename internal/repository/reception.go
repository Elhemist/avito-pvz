package repository

import (
	"database/sql"
	"fmt"
	"pvz-test/internal/models"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ReceptionPostgres struct {
	db *sqlx.DB
}

func NewReceptionPostgres(db *sqlx.DB) *ReceptionPostgres {
	return &ReceptionPostgres{db: db}
}

func (r *ReceptionPostgres) CreateReception(pvzID uuid.UUID) (models.Reception, error) {
	tx := r.db.MustBegin()

	var exists bool
	err := tx.Get(&exists, `
		SELECT EXISTS (
			SELECT 1 FROM receptions
			WHERE pvz_id = $1 AND status = 'in_progress'
		)
	`, pvzID)
	if err != nil {
		tx.Rollback()
		return models.Reception{}, err
	}
	if exists {
		tx.Rollback()
		return models.Reception{}, fmt.Errorf("pvz: %s have active reception", pvzID.String())
	}

	var reception models.Reception
	err = tx.Get(&reception, `
		INSERT INTO receptions (pvz_id, status)
		VALUES ($1, 'in_progress')
		RETURNING id, pvz_id, status, created_at
	`, pvzID)
	if err != nil {
		tx.Rollback()
		return models.Reception{}, err
	}

	if err := tx.Commit(); err != nil {
		return models.Reception{}, err
	}

	return reception, nil
}

func (r *ReceptionPostgres) AddItem(pvzID uuid.UUID, itemType string) (models.Item, error) {
	tx := r.db.MustBegin()

	var receptionID uuid.UUID
	err := tx.Get(&receptionID, `
		SELECT id
		FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY created_at DESC
		LIMIT 1
	`, pvzID)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return models.Item{}, fmt.Errorf("no active reception for PVZ %s", pvzID)
		}
		return models.Item{}, err
	}
	var item models.Item
	err = tx.Get(&item, `
		INSERT INTO goods (reception_id, pvz_id, type, added_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING *
	`, receptionID, pvzID, itemType)
	if err != nil {
		tx.Rollback()
		return models.Item{}, fmt.Errorf("failed to insert item: %w", err)
	}

	return item, nil
}

func (r *ReceptionPostgres) DeleteItem(pvzID uuid.UUID) error {
	tx := r.db.MustBegin()

	var receptionID uuid.UUID
	err := tx.Get(&receptionID, `
		SELECT id
		FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY created_at DESC
		LIMIT 1
	`, pvzID)
	if err != nil {
		tx.Rollback()
		if err == sql.ErrNoRows {
			return fmt.Errorf("no active reception for pvz %s", pvzID.String())
		}
		return err
	}

	var item models.Item
	err = tx.Get(&item, `
		SELECT id, reception_id, pvz_id, type, added_at
		FROM goods
		WHERE reception_id = $1
		ORDER BY added_at DESC
		LIMIT 1
	`, receptionID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM goods
		WHERE id = $1
	`, item.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *ReceptionPostgres) GetActiveReception(pvzID uuid.UUID) (models.Reception, error) {
	var reception models.Reception
	err := r.db.Get(&reception, `
		SELECT id, pvz_id, status, created_at
		FROM receptions
		WHERE pvz_id = $1 AND status = 'in_progress'
		ORDER BY created_at DESC
		LIMIT 1
	`, pvzID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Reception{}, nil
		}
		return models.Reception{}, fmt.Errorf("failed to get active reception for PVZ %s: %w", pvzID.String(), err)
	}

	return reception, nil
}

func (r *ReceptionPostgres) CloseReception(receptionID uuid.UUID) error {
	res, err := r.db.Exec(`
		UPDATE receptions
		SET status = 'closed'
		WHERE id = $1 AND status = 'in_progress'
	`, receptionID)
	if err != nil {
		return fmt.Errorf("failed to close reception %s: %w", receptionID.String(), err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not determine result of reception close: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("reception %s is already closed or does not exist", receptionID.String())
	}

	return nil
}

func (r *ReceptionPostgres) GetReceptionsWithProducts(pvzID uuid.UUID, start, end *time.Time) ([]models.Reception, error) {
	query := sq.
		Select("id", "pvz_id", "created_at", "status").
		From("receptions").
		Where(sq.Eq{"pvz_id": pvzID}).
		OrderBy("created_at DESC")

	if start != nil {
		query = query.Where(sq.GtOrEq{"created_at": *start})
	}
	if end != nil {
		query = query.Where(sq.LtOrEq{"created_at": *end})
	}

	sqlQuery, args, err := query.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	var receptions []models.Reception
	err = r.db.Select(&receptions, sqlQuery, args...)
	return receptions, err
}

func (r *ReceptionPostgres) GetItemsByReceptionID(receptionID uuid.UUID) ([]models.Item, error) {
	var items []models.Item
	err := r.db.Select(&items, `
        SELECT id, reception_id, pvz_id, type, added_at
        FROM items
        WHERE reception_id = $1
        ORDER BY added_at ASC
    `, receptionID)
	return items, err
}

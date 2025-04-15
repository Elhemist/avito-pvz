package repository_test

import (
	"database/sql"
	"fmt"
	"pvz-test/internal/models"
	"pvz-test/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestReceptionPostgres_CreateReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB)

	t.Run("Successful creation", func(t *testing.T) {
		pvzID := uuid.New()
		expectedReception := models.Reception{
			ID:        uuid.New(),
			PVZID:     pvzID,
			Status:    "in_progress",
			CreatedAt: time.Now(),
		}

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM receptions WHERE pvz_id = \$1 AND status = 'in_progress' \)`).
			WithArgs(pvzID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))
		mock.ExpectQuery(`INSERT INTO receptions \(pvz_id, status\) VALUES \(\$1, 'in_progress'\) RETURNING id, pvz_id, status, created_at`).
			WithArgs(pvzID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "status", "created_at"}).
				AddRow(expectedReception.ID, expectedReception.PVZID, expectedReception.Status, expectedReception.CreatedAt))
		mock.ExpectCommit()

		reception, err := repo.CreateReception(pvzID)
		assert.NoError(t, err)
		assert.Equal(t, expectedReception, reception)
	})

	t.Run("Active reception exists", func(t *testing.T) {
		pvzID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM receptions WHERE pvz_id = \$1 AND status = 'in_progress' \)`).
			WithArgs(pvzID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
		mock.ExpectRollback()

		_, err := repo.CreateReception(pvzID)
		assert.EqualError(t, err, fmt.Sprintf("pvz: %s have active reception", pvzID.String()))
	})
}

func TestReceptionPostgres_GetActiveReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB)

	t.Run("Reception found", func(t *testing.T) {
		pvzID := uuid.New()
		expectedReception := models.Reception{
			ID:        uuid.New(),
			PVZID:     pvzID,
			Status:    "in_progress",
			CreatedAt: time.Now(),
		}

		mock.ExpectQuery(`SELECT id, pvz_id, status, created_at FROM receptions WHERE pvz_id = \$1 AND status = 'in_progress' ORDER BY created_at DESC LIMIT 1`).
			WithArgs(pvzID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "status", "created_at"}).
				AddRow(expectedReception.ID, expectedReception.PVZID, expectedReception.Status, expectedReception.CreatedAt))

		reception, err := repo.GetActiveReception(pvzID)
		assert.NoError(t, err)
		assert.Equal(t, expectedReception, reception)
	})

	t.Run("No active reception", func(t *testing.T) {
		pvzID := uuid.New()

		mock.ExpectQuery(`SELECT id, pvz_id, status, created_at FROM receptions WHERE pvz_id = \$1 AND status = 'in_progress' ORDER BY created_at DESC LIMIT 1`).
			WithArgs(pvzID).
			WillReturnError(sql.ErrNoRows)

		reception, err := repo.GetActiveReception(pvzID)
		assert.NoError(t, err)
		assert.Equal(t, models.Reception{}, reception)
	})
}

func TestReceptionPostgres_CloseReception(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB)

	t.Run("Successful closure", func(t *testing.T) {
		receptionID := uuid.New()

		mock.ExpectExec(`UPDATE receptions SET status = 'closed' WHERE id = \$1 AND status = 'in_progress'`).
			WithArgs(receptionID).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CloseReception(receptionID)
		assert.NoError(t, err)
	})

	t.Run("Reception already closed", func(t *testing.T) {
		receptionID := uuid.New()

		mock.ExpectExec(`UPDATE receptions SET status = 'closed' WHERE id = \$1 AND status = 'in_progress'`).
			WithArgs(receptionID).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.CloseReception(receptionID)
		assert.EqualError(t, err, "reception "+receptionID.String()+" is already closed or does not exist")
	})
}

func TestReceptionPostgres_AddItem_Good(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB)

	fixedTime := time.Date(2025, 4, 16, 22, 29, 15, 0, time.UTC)

	pvzID := uuid.New()
	receptionID := uuid.New()
	itemID := uuid.New()

	expectedItem := models.Item{
		ID:          itemID,
		ReceptionID: receptionID,
		PvzID:       pvzID,
		Type:        models.ItemTypeElectronics,
		AddedAt:     fixedTime,
	}

	mock.ExpectBegin()

	mock.ExpectQuery(`SELECT id FROM receptions WHERE pvz_id = \$1 AND status = 'in_progress' ORDER BY created_at DESC LIMIT 1`).
		WithArgs(pvzID).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionID))

	mock.ExpectQuery(`INSERT INTO goods \(reception_id, pvz_id, type, added_at\) VALUES \(\$1, \$2, \$3, NOW\(\)\) RETURNING \*`).
		WithArgs(receptionID, pvzID, expectedItem.Type).
		WillReturnRows(sqlmock.NewRows([]string{"id", "reception_id", "pvz_id", "type", "added_at"}).
			AddRow(
				expectedItem.ID,
				expectedItem.ReceptionID,
				expectedItem.PvzID,
				expectedItem.Type,
				expectedItem.AddedAt,
			))

	mock.ExpectCommit()

	item, err := repo.AddItem(pvzID, string(expectedItem.Type))
	assert.NoError(t, err, "unexpected error: %v", err)

	assert.Equal(t, expectedItem.ID, item.ID)
	assert.Equal(t, expectedItem.ReceptionID, item.ReceptionID)
	assert.Equal(t, expectedItem.PvzID, item.PvzID)
	assert.Equal(t, expectedItem.Type, item.Type)
	assert.WithinDuration(t, expectedItem.AddedAt, item.AddedAt, time.Second)
}

func TestReceptionPostgres_AddItem_Bad(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB)

	t.Run("No active reception", func(t *testing.T) {
		pvzID := uuid.New()
		itemType := "electronics"

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT id FROM receptions WHERE pvz_id = \$1 AND status = 'in_progress' ORDER BY created_at DESC LIMIT 1`).
			WithArgs(pvzID).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		_, err := repo.AddItem(pvzID, itemType)
		assert.EqualError(t, err, "no active reception for PVZ "+pvzID.String())
	})
}

func TestReceptionPostgres_DeleteItem(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB)

	t.Run("Successful item deletion", func(t *testing.T) {
		pvzID := uuid.New()
		receptionID := uuid.New()
		itemID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT id FROM receptions WHERE pvz_id = \$1 AND status = 'in_progress' ORDER BY created_at DESC LIMIT 1`).
			WithArgs(pvzID).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(receptionID))
		mock.ExpectQuery(`SELECT id, reception_id, pvz_id, type, added_at FROM goods WHERE reception_id = \$1 ORDER BY added_at DESC LIMIT 1`).
			WithArgs(receptionID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "reception_id", "pvz_id", "type", "added_at"}).
				AddRow(itemID, receptionID, pvzID, "electronics", time.Now()))
		mock.ExpectExec(`DELETE FROM goods WHERE id = \$1`).
			WithArgs(itemID).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.DeleteItem(pvzID)
		assert.NoError(t, err)
	})

	t.Run("No active reception", func(t *testing.T) {
		pvzID := uuid.New()

		mock.ExpectBegin()
		mock.ExpectQuery(`SELECT id FROM receptions WHERE pvz_id = \$1 AND status = 'in_progress' ORDER BY created_at DESC LIMIT 1`).
			WithArgs(pvzID).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectRollback()

		err := repo.DeleteItem(pvzID)
		assert.EqualError(t, err, "no active reception for pvz "+pvzID.String())
	})
}

func TestReceptionPostgres_GetReceptionsWithProducts(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB)

	t.Run("Successful fetch", func(t *testing.T) {
		pvzID := uuid.New()
		expectedReceptions := []models.Reception{
			{ID: uuid.New(), PVZID: pvzID, Status: "closed", CreatedAt: time.Now()},
		}

		mock.ExpectQuery(`SELECT id, pvz_id, created_at, status FROM receptions WHERE pvz_id = \$1 ORDER BY created_at DESC`).
			WithArgs(pvzID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "pvz_id", "created_at", "status"}).
				AddRow(expectedReceptions[0].ID, expectedReceptions[0].PVZID, expectedReceptions[0].CreatedAt, expectedReceptions[0].Status))

		receptions, err := repo.GetReceptionsWithProducts(pvzID, nil, nil)
		assert.NoError(t, err)
		assert.Equal(t, expectedReceptions, receptions)
	})
}

func TestReceptionPostgres_GetItemsByReceptionID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewReceptionPostgres(sqlxDB)

	t.Run("Successful fetch", func(t *testing.T) {
		receptionID := uuid.New()
		pvzID := uuid.New()
		expectedItems := []models.Item{
			{ID: uuid.New(), ReceptionID: receptionID, PvzID: pvzID, Type: models.ItemTypeElectronics, AddedAt: time.Now()},
		}

		mock.ExpectQuery(`SELECT id, reception_id, pvz_id, type, added_at FROM items WHERE reception_id = \$1 ORDER BY added_at ASC`).
			WithArgs(receptionID).
			WillReturnRows(sqlmock.NewRows([]string{"id", "reception_id", "pvz_id", "type", "added_at"}).
				AddRow(expectedItems[0].ID, expectedItems[0].ReceptionID, expectedItems[0].PvzID, expectedItems[0].Type, expectedItems[0].AddedAt))

		items, err := repo.GetItemsByReceptionID(receptionID)
		assert.NoError(t, err)
		assert.Equal(t, expectedItems, items)
	})
}

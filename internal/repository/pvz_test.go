package repository_test

import (
	"errors"
	"pvz-test/internal/models"
	"pvz-test/internal/repository"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestPvzPostgres_CreatePvz(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewPvzPostgres(sqlxDB)

	t.Run("Successful creation", func(t *testing.T) {
		city := "Москва"
		expectedPvz := models.PVZ{
			ID:               uuid.New(),
			RegistrationDate: time.Now(),
			City:             city,
		}

		mock.ExpectQuery(`INSERT INTO pvz \(city\) VALUES \(\$1\) RETURNING id, city, registration_date`).
			WithArgs(city).
			WillReturnRows(sqlmock.NewRows([]string{"id", "city", "registration_date"}).
				AddRow(expectedPvz.ID, expectedPvz.City, expectedPvz.RegistrationDate))

		pvz, err := repo.CreatePvz(city)
		assert.NoError(t, err)
		assert.Equal(t, expectedPvz, pvz)
	})

	t.Run("Database error", func(t *testing.T) {
		city := "Казань"

		mock.ExpectQuery(`INSERT INTO pvz \(city\) VALUES \(\$1\) RETURNING id, city, registration_date`).
			WithArgs(city).
			WillReturnError(errors.New("database error"))

		_, err := repo.CreatePvz(city)
		assert.EqualError(t, err, "database error")
	})
}

func TestPvzPostgres_Exists(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewPvzPostgres(sqlxDB)

	t.Run("PVZ exists", func(t *testing.T) {
		pvzID := uuid.New()

		mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM pvz WHERE id = \$1 \)`).
			WithArgs(pvzID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := repo.Exists(pvzID)
		assert.NoError(t, err)
		assert.True(t, exists)
	})

	t.Run("PVZ does not exist", func(t *testing.T) {
		pvzID := uuid.New()

		mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM pvz WHERE id = \$1 \)`).
			WithArgs(pvzID).
			WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := repo.Exists(pvzID)
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("Database error", func(t *testing.T) {
		pvzID := uuid.New()

		mock.ExpectQuery(`SELECT EXISTS \( SELECT 1 FROM pvz WHERE id = \$1 \)`).
			WithArgs(pvzID).
			WillReturnError(errors.New("database error"))

		_, err := repo.Exists(pvzID)
		assert.EqualError(t, err, "failed to check PVZ existence: database error")
	})
}

func TestPvzPostgres_GetPVZList(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	sqlxDB := sqlx.NewDb(db, "sqlmock")
	repo := repository.NewPvzPostgres(sqlxDB)

	t.Run("Successful fetch", func(t *testing.T) {
		limit, offset := 10, 0
		expectedPVZs := []models.PVZ{
			{ID: uuid.New(), RegistrationDate: time.Now(), City: "Москва"},
			{ID: uuid.New(), RegistrationDate: time.Now(), City: "Казань"},
		}

		mock.ExpectQuery(`SELECT id, registration_date, city FROM pvz ORDER BY registration_date DESC LIMIT \$1 OFFSET \$2`).
			WithArgs(limit, offset).
			WillReturnRows(sqlmock.NewRows([]string{"id", "registration_date", "city"}).
				AddRow(expectedPVZs[0].ID, expectedPVZs[0].RegistrationDate, expectedPVZs[0].City).
				AddRow(expectedPVZs[1].ID, expectedPVZs[1].RegistrationDate, expectedPVZs[1].City))

		pvzs, err := repo.GetPVZList(limit, offset)
		assert.NoError(t, err)
		assert.Equal(t, expectedPVZs, pvzs)
	})

	t.Run("Database error", func(t *testing.T) {
		limit, offset := 10, 0

		mock.ExpectQuery(`SELECT id, registration_date, city FROM pvz ORDER BY registration_date DESC LIMIT \$1 OFFSET \$2`).
			WithArgs(limit, offset).
			WillReturnError(errors.New("database error"))

		_, err := repo.GetPVZList(limit, offset)
		assert.EqualError(t, err, "database error")
	})
}

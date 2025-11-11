// Package repository_test contains unit tests for the BeerRepository
// implementation using SQLMock and testify.
package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/samuskitchen/go-health-checker/beer/model"
	"github.com/samuskitchen/go-health-checker/configs/storage"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const (
	fakeBeerIdUint uint = 1 // Example beer ID for tests
)

// dataBeers returns a sample of beer data to use in tests.
func dataBeers() []model.Beers {
	now := time.Now()

	return []model.Beers{
		{
			ID:        fakeBeerIdUint,
			Name:      "Gulden Draak",
			Brewery:   "Blót",
			Country:   "BE",
			Price:     6.50,
			Currency:  "EUR",
			CreatedAt: now,
			UpdatedAt: now,
		},
		{
			ID:        uint(2),
			Name:      "Club Colombia",
			Brewery:   "Bavaria",
			Country:   "CO",
			Price:     3.483,
			Currency:  "COP",
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

// Test_beerRepository_GetAllBeers validates the functionality to retrieve all beers.
func Test_beerRepository_GetAllBeers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	defer func() {
		mock.ExpectClose()
		if errDB := db.Close(); errDB != nil {
			log.Error().Msgf("Error closing the database connection: %v", errDB)
		}
		assert.NoError(t, mock.ExpectationsWereMet())
	}()

	data := &storage.Data{
		DB: db,
	}

	repo := NewBeerRepository(data)
	ctx := context.Background()
	beersTest := dataBeers()

	t.Run("Success SQL", func(tt *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at", "updated_at"})
		for _, beer := range beersTest {
			rows.AddRow(beer.ID, beer.Name, beer.Brewery, beer.Country, beer.Price, beer.Currency, beer.CreatedAt, beer.UpdatedAt)
		}
		mock.ExpectQuery(regexp.QuoteMeta(selectAllBeers)).WillReturnRows(rows)

		gotBeers, errRepo := repo.GetAllBeers(ctx)
		assert.NoError(t, errRepo)
		assert.Equal(t, beersTest, gotBeers)
	})

	t.Run("Error SQL", func(tt *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(selectAllBeers)).WillReturnError(assert.AnError)
		gotBeers, errRepo := repo.GetAllBeers(ctx)
		assert.Error(t, errRepo)
		assert.Empty(t, gotBeers)
	})

	t.Run("No Results", func(tt *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name", "brewery", "country", "price", "currency", "created_at", "updated_at"})
		mock.ExpectQuery(regexp.QuoteMeta(selectAllBeers)).WillReturnRows(rows)
		gotBeers, errRepo := repo.GetAllBeers(ctx)
		assert.NoError(t, errRepo)
		assert.Empty(t, gotBeers)
	})
}

// Package repository implements data access logic for beers,
// using PostgreSQL as storage. Here the concrete repository is defined
// that satisfies the BeerRepository interface.
package repository

import (
	"context"
	// interfaces defines the contract that the repository must fulfill.
	"github.com/samuskitchen/go-health-checker/beer/interfaces"
	"github.com/samuskitchen/go-health-checker/configs/storage"

	// model contains domain structures (e.g., model.Beers).
	"github.com/samuskitchen/go-health-checker/beer/model"

	// zerolog for structured logging in each method.
	"github.com/rs/zerolog/log"
)

const (
	// selectAllBeers is a query that selects all rows from the beers table
	selectAllBeers = "SELECT id, \"name\", brewery, country_code, price, currency, created_at, updated_at FROM beers;"
)

// beerRepository is the implementation of BeerRepository that uses
// *sql.DB to communicate with PostgreSQL.
type beerRepository struct {
	connection *storage.Data
}

// NewBeerRepository builds an instance of BeerRepository using the given connection.
func NewBeerRepository(db *storage.Data) interfaces.BeerRepository {
	return &beerRepository{
		connection: db,
	}
}

// GetAllBeers retrieves all beers registered in the database.
//
// Parameters:
//   - ctx: context for timeout and cancellation control.
//
// Returns:
//   - []model.Beers: slice with all beers.
//   - error: in case of failure in the query or in row scanning.
func (pb *beerRepository) GetAllBeers(ctx context.Context) ([]model.Beers, error) {
	// Logger with Method field to track log origin.
	subLogger := log.With().Str("Method", "BeerRepository.GetAllBeers").Logger()
	subLogger.Info().Msg("INIT")

	// Execute the parameterized query defined in selectAllBeers.
	rows, err := pb.connection.DB.QueryContext(ctx, selectAllBeers)
	if err != nil {
		subLogger.Error().Msgf("error executing query: %v", err)
		return nil, err
	}

	// Ensure rows are closed when finished.
	defer func() {
		if errClose := rows.Close(); errClose != nil {
			subLogger.Error().Msgf("error closing rows: %v", errClose)
		}
	}()

	// Iterate over each row and map to model.Beers.
	var beers []model.Beers
	for rows.Next() {
		var beerRow model.Beers
		if errScan := rows.Scan(
			&beerRow.ID,
			&beerRow.Name,
			&beerRow.Brewery,
			&beerRow.Country,
			&beerRow.Price,
			&beerRow.Currency,
			&beerRow.CreatedAt,
			&beerRow.UpdatedAt,
		); errScan != nil {
			subLogger.Error().Msgf("error scanning row: %v", errScan)
			return nil, errScan
		}
		beers = append(beers, beerRow)
	}

	subLogger.Info().Msgf("END_OK")
	return beers, nil
}

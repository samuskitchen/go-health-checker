// Package interfaces define the repository contracts for the geography module,
package interfaces

import (
	"context"

	"github.com/samuskitchen/go-health-checker/beer/model"
)

// BeerRepository define the repository contract for the BeerRepository
type BeerRepository interface {
	GetAllBeers(ctx context.Context) ([]model.Beers, error)
}

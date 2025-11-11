// Package interfaces define the service layer contracts for the geography module
package interfaces

import (
	"context"

	"github.com/samuskitchen/go-health-checker/beer/model"
)

// BeerService define the service layer contract for the BeerService
type BeerService interface {
	GetAllBeers(ctx context.Context) ([]model.BeersResponse, error) // now returns ready response
}

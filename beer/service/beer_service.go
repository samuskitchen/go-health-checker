// Package service contains the implementation of use cases for the Beer entity.
// Orchestrates the interaction between the beer repository and the currency conversion
// client to provide high-level functionality.
package service

import (
	"context"

	"github.com/samuskitchen/go-health-checker/beer/interfaces"
	"github.com/samuskitchen/go-health-checker/beer/model"
	"github.com/samuskitchen/go-health-checker/configs/cache"
	events "github.com/samuskitchen/go-health-checker/configs/event"

	"github.com/rs/zerolog/log"
)

// beerService implements interfaces.BeerService.
// Combines a beer repository and a currency conversion client.
type beerService struct {
	beerRepository interfaces.BeerRepository
	hazelcast      *cache.Cache
	rabbit         *events.RabbitEvent
}

// NewBeerService creates a new instance of BeerService.
func NewBeerService(
	beerRepository interfaces.BeerRepository, hazelcast *cache.Cache, rabbit *events.RabbitEvent,
) interfaces.BeerService {
	return &beerService{
		beerRepository: beerRepository,
		hazelcast:      hazelcast,
		rabbit:         rabbit,
	}
}

// GetAllBeers retrieves all beers from the database.
func (b *beerService) GetAllBeers(ctx context.Context) ([]model.BeersResponse, error) {
	subLogger := log.With().Str("Method", "BeerService.GetAllBeers").Logger()
	subLogger.Info().Msg("INIT")

	beers, err := b.beerRepository.GetAllBeers(ctx)
	if err != nil {
		subLogger.Error().Msgf("error GetAllBeers repo: %v", err)
		return nil, err
	}

	resp := make([]model.BeersResponse, 0, len(beers))
	for _, v := range beers {
		resp = append(resp, v.ToBeersResponse())
	}

	subLogger.Info().Msg("END_OK")
	return resp, nil
}

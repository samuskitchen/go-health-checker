// Package service_test contains unit tests for the BeerService service layer.
// Uses testify for assertions and mocks of repository and external client interfaces.
package service

import (
	"context"
	"testing"
	"time"

	_mockInterfaces "github.com/samuskitchen/go-health-checker/beer/mocks/interfaces"
	"github.com/samuskitchen/go-health-checker/configs/cache"
	events "github.com/samuskitchen/go-health-checker/configs/event"
	_mockToolsBroker "github.com/samuskitchen/go-health-checker/pkg/tools/mocks/broker"
	_mockToolsDataStore "github.com/samuskitchen/go-health-checker/pkg/tools/mocks/data_store"

	"github.com/samuskitchen/go-health-checker/beer/model"

	"github.com/stretchr/testify/assert"
)

const (
	fakeBeerIdUint  uint   = 1     // Example beer ID
	fakeQuantityInt int    = 6     // Example quantity for GetOneBoxPrice
	fakeCurrencyStr string = "USD" // Example currency
)

// dataBeers prepares model.Beers data for use in tests.
func dataBeers() []model.Beers {
	now := time.Now()
	return []model.Beers{
		{ID: fakeBeerIdUint, Name: "Gulden Draak", Brewery: "Blót", Country: "BE", Price: 6.50, Currency: "EUR", CreatedAt: now, UpdatedAt: now},
		{ID: 2, Name: "Club Colombia", Brewery: "Bavaria", Country: "CO", Price: 3.483, Currency: "COP", CreatedAt: now, UpdatedAt: now},
	}
}

// Test_beerService_GetAllBeers placeholder for GetAllBeers tests.
func Test_beerService_GetAllBeers(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		beers := dataBeers()
		beersResponse := make([]model.BeersResponse, 0, len(beers))

		for _, b := range beers {
			beersResponse = append(beersResponse, b.ToBeersResponse())
		}

		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockBroker := _mockToolsBroker.NewMockClient(t)
		mockDataStore := _mockToolsDataStore.NewMockIClient(t)

		hazelcast := &cache.Cache{
			Hazelcast: mockDataStore,
		}

		rabbitMq := &events.RabbitEvent{
			RabbitMQClient: mockBroker,
		}

		service := NewBeerService(mockRepository, hazelcast, rabbitMq)

		mockRepository.On("GetAllBeers", ctx).Return(beers, nil)

		gotBeers, errService := service.GetAllBeers(ctx)
		assert.NoError(t, errService)
		assert.Equal(t, beersResponse, gotBeers)
		mockRepository.AssertExpectations(t)
	})

	t.Run("error repository", func(t *testing.T) {
		mockRepository := _mockInterfaces.NewMockBeerRepository(t)
		mockBroker := _mockToolsBroker.NewMockClient(t)
		mockDataStore := _mockToolsDataStore.NewMockIClient(t)

		hazelcast := &cache.Cache{
			Hazelcast: mockDataStore,
		}

		rabbitMq := &events.RabbitEvent{
			RabbitMQClient: mockBroker,
		}

		service := NewBeerService(mockRepository, hazelcast, rabbitMq)

		mockRepository.On("GetAllBeers", ctx).Return(nil, assert.AnError)

		gotBeers, errService := service.GetAllBeers(ctx)
		assert.Error(t, errService)
		assert.Nil(t, gotBeers)
		mockRepository.AssertExpectations(t)
	})
}

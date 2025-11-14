// Package injector dependency injection
package injector

import (
	"fmt"

	"github.com/samuskitchen/go-health-checker/beer/handler"
	"github.com/samuskitchen/go-health-checker/beer/repository"
	"github.com/samuskitchen/go-health-checker/beer/service"
	"github.com/samuskitchen/go-health-checker/configs/cache"
	events "github.com/samuskitchen/go-health-checker/configs/event"
	"github.com/samuskitchen/go-health-checker/configs/generals/router"
	"github.com/samuskitchen/go-health-checker/configs/storage"
	echo "github.com/samuskitchen/go-health-checker/pkg/tools/server"

	"go.uber.org/dig"
)

// Container dig container variable
var Container *dig.Container

// BuildContainer dependency injection wrapper function
func BuildContainer() *dig.Container {
	Container = dig.New()

	// DB / Cache
	checkError(Container.Provide(storage.PostgresConnection))
	checkError(Container.Provide(cache.HazelcastConnection))

	// Broker
	checkError(Container.Provide(events.RabbitConnection))

	// Router / Server
	checkError(Container.Provide(echo.NewServer))
	checkError(Container.Provide(router.NewRouter))

	// Health Check
	checkError(Container.Provide(router.NewHealthHandler))

	// Handlers
	checkError(Container.Provide(handler.NewBeerHandler))

	// Services
	checkError(Container.Provide(service.NewBeerService))

	// Repository
	checkError(Container.Provide(repository.NewBeerRepository))

	return Container
}

func checkError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Error injecting %v", err))
	}
}

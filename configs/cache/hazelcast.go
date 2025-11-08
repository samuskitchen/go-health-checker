// Package cache provides a singleton connection to Hazelcast
// to manage the application cache.
package cache

import (
	"context"
	"os"
	"sync"

	"github.com/samuskitchen/go-health-checker/pkg/kit/enums"
	hazelcast "github.com/samuskitchen/go-health-checker/pkg/tools/data_store"
	modelCache "github.com/samuskitchen/go-health-checker/pkg/tools/models"

	"github.com/rs/zerolog/log"
)

var (
	once      sync.Once
	dataCache *Cache
)

// Cache wraps the Hazelcast singleton client
// and exposes the application's caching functionality.
type Cache struct {
	Hazelcast hazelcast.IClient
}

// HazelcastConnection returns the singleton Cache instance that maintains
// the connection to the Hazelcast cluster. If it isn't already initialized, it creates it.
func HazelcastConnection() *Cache {
	once.Do(getConnection)
	return dataCache
}

func getConnection() {
	cacheConfigs := modelCache.Config{
		ClusterName: enums.HazelClusterName,
		Addresses:   []string{os.Getenv(enums.HazelServer)},
		ClientName:  enums.HazelClientName,
	}

	conn, err := hazelcast.NewClientHazelcast(cacheConfigs)
	if err != nil {
		log.Error().Msgf("Error connecting to database: %v", err)
	}

	dataCache = &Cache{
		Hazelcast: conn,
	}
}

// HazelcastCloseConnection closes the Hazelcast singleton connection if it has been initialized.
// Logs fatal on error closing.
func HazelcastCloseConnection() {
	if dataCache != nil {
		if err := dataCache.Hazelcast.Disconnect(context.Background()); err != nil {
			log.Fatal().Msgf("Error closing the database: %v", err)
		}
	}
}

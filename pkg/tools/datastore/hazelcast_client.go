// Package datastore provides in-memory data storage clients such as Redis and Hazelcast,
// designed for high-performance caching, temporary data persistence, and low-latency access.
package datastore

import (
	"context"
	"errors"
	"fmt"

	"github.com/hazelcast/hazelcast-go-client"
	tools "github.com/samuskitchen/go-health-checker/pkg/tools/models"

	"github.com/rs/zerolog/log"
)

// ClientHazelcast is the concrete implementation with enhanced type operations.
type ClientHazelcast struct {
	Client *hazelcast.Client
}

// NewClientHazelcast initializes and returns a new Hazelcast client instance.
//
// This method sets up a connection to a Hazelcast cluster using the provided configuration.
// It validates the configuration fields, builds the Hazelcast client settings, and establishes
// the client connection to the cluster.
//
// Returns:
//   - A valid IClient implementation if the connection is successful.
//   - An error if required configuration fields are missing or the connection cannot be established.
//
// Expected errors:
//   - "missing required field: addresses" if no cluster addresses are provided.
//   - "missing required field: clusterName" if the cluster name is empty.
//   - "failed to start Hazelcast client" if the client fails to connect.
func NewClientHazelcast(config tools.Config) (IClient, error) {
	if len(config.Addresses) == 0 {
		return nil, errors.New("missing required field: addresses")
	}

	if config.ClusterName == "" {
		return nil, errors.New("missing required field: clusterName")
	}

	hzConfig := config.ToHazelcastConfig()
	hzClient, err := hazelcast.StartNewClientWithConfig(context.Background(), hzConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to start Hazelcast client: %w", err)
	}

	log.Info().Msgf("Successfully connected to Hazelcast cluster: %s", config.ClusterName)
	return &ClientHazelcast{
		Client: hzClient,
	}, nil
}

// Disconnect gracefully shuts down the active Hazelcast client connection.
//
// This method terminates the connection between the client and the Hazelcast cluster,
// ensuring that all resources are properly released.
//
// Returns an error if:
//   - The shutdown process fails due to communication or internal client issues.
//   - The client is already closed or uninitialized.
func (ch *ClientHazelcast) Disconnect(ctx context.Context) error {
	log.Info().Msgf("Disconnecting from Hazelcast cluster: %s", ch.Client.Name())
	return ch.Client.Shutdown(ctx)
}

// Ping verifies that the Hazelcast client connection is active and healthy.
//
// This method checks if the Hazelcast client is running and connected to the cluster.
// It's designed for health checks and monitoring the connection status.
//
// Returns an error if:
//   - The client is nil
//   - The client is not running (shutdown or disconnected)
func (ch *ClientHazelcast) Ping() error {
	if ch.Client == nil {
		return fmt.Errorf("hazelcast client is not initialized")
	}

	if !ch.Client.Running() {
		return fmt.Errorf("hazelcast client is not running")
	}

	return nil
}

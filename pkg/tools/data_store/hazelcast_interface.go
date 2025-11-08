// Package data_store provides a generic interface and helper types
// to connect to and operate against a Hazelcast cluster in a
// typed and configurable way.
package data_store

import (
	"context"
)

// IClient defines the generic interface for a type-safe Hazelcast client.
type IClient interface {
	// Disconnect gracefully shuts down the connection to the Hazelcast cluster.
	Disconnect(ctx context.Context) error

	// Ping verifies that the Hazelcast client connection is active and healthy.
	Ping() error
}

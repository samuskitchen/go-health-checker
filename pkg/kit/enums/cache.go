// Package enums define environment configuration constants
// and general values used for Hazelcast.
package enums

import "time"

// Environment variables to configure the connection to Hazelcast.
const (
	// HazelServer specifies the address (host:port) of the Hazelcast server.
	HazelServer string = "HAZEL_SERVER"
	// HazelClusterName specifies the name of the Hazelcast cluster.
	HazelClusterName string = "hz-cache-cluster"
	// HazelClientName specifies the name the Hazelcast client will use.
	HazelClientName string = "health-checker-cluster"
	// CacheGeneralTTL defines the time to live (24h) for the general cache.
	CacheGeneralTTL time.Duration = 24 * time.Hour
)

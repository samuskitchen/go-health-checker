// Package tools package provides utilities and frameworks for configuring and connecting to
// a Hazelcast cluster generically.
package tools

import "github.com/hazelcast/hazelcast-go-client"

// Config contains the settings for connecting to Hazelcast.
type Config struct {
	Addresses    []string
	ClusterName  string
	ClientName   string
	ClientLabels []string
}

// ToHazelcastConfig converts our configuration structure to that of Hazelcast.
func (c Config) ToHazelcastConfig() hazelcast.Config {
	var cfg hazelcast.Config

	cfg.Cluster.Name = c.ClusterName
	cfg.Cluster.Network.SetAddresses(c.Addresses...)
	cfg.ClientName = c.ClientName
	cfg.Labels = c.ClientLabels

	return cfg
}

// Package broker provides a thread-safe RabbitMQ client
// for local development, with automatic reconnection and channel handling.
package broker

import (
	"fmt"
	"sync"
	"time"

	tools "github.com/samuskitchen/go-health-checker/pkg/tools/models"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

// clientImpl implements the Client interface for RabbitMQ.
//
// This struct manages the AMQP connection, channel, and all concurrency primitives required
// for safe parallel consumption and publishing.
// It is not intended to be used directly; use the NewClient constructor and the Client interface for best results.
type clientImpl struct {
	connection *amqp.Connection // Underlying AMQP connection
	channel    *amqp.Channel    // AMQP channel for operations
	mu         sync.Mutex       // Mutex for thread safety on connection/channel
	params     tools.Params     // Connection parameters
	closeCh    chan struct{}    // Used to close goroutines and signal shutdown
}

// NewClient returns a new concurrent-safe RabbitMQ client.
//
// This constructor creates a new client instance that implements the Client interface.
// The client is initialized with internal channels and mutexes for thread safety.
// Use this function to create client instances; do not create clientImpl directly.
//
// Returns a Client interface that can be used for all RabbitMQ operations.
func NewClient() Client {
	return &clientImpl{
		closeCh: make(chan struct{}),
	}
}

// ConnectLocal establishes a non-secure, thread-safe connection to RabbitMQ for local development.
//
// This method creates a standard, non-TLS connection and is intended for use in
// local or trusted environments where setting up TLS is unnecessary. It sets up
// the internal AMQP channel and starts a background goroutine for auto-reconnection.
//
// Parameters:
//   - host: RabbitMQ server hostname or IP address
//   - port: RabbitMQ server port (typically "5672" for standard AMQP)
//   - user: Username for authentication
//   - password: Password for authentication
//
// Returns an error if the connection cannot be established.
func (c *clientImpl) ConnectLocal(host, port, user, password string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.params = tools.Params{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Vhost:    "/", // Use default virtual host
	}

	return c.establishLocalConnection()
}

// Close gracefully closes the connection and channel, and signals all goroutines to stop.
//
// This method performs a graceful shutdown of the client by:
//   - Signaling all background goroutines to stop via the closeCh channel
//   - Closing the AMQP channel
//   - Closing the AMQP connection
//
// The method is thread-safe and idempotent - it can be called multiple times
// safely. It ensures that all resources are properly cleaned up.
//
// Returns an error if any part of the shutdown process fails.
func (c *clientImpl) Close() error {
	// Lock mutex to ensure thread-safe shutdown
	c.mu.Lock()
	defer c.mu.Unlock()

	// Signal all goroutines to stop (idempotent operation)
	select {
	case <-c.closeCh:
		// Client is already closed, do nothing
	default:
		close(c.closeCh) // Signal shutdown to all goroutines
	}

	// Close the AMQP channel first
	if err := c.closeChannel(); err != nil {
		return err
	}

	// Then close the AMQP connection
	if err := c.closeConnection(); err != nil {
		return err
	}

	log.Info().Msg("RabbitMQ client closed")
	return nil
}

// Ping verifies that the RabbitMQ connection is active and healthy.
//
// This method checks the status of the AMQP connection and channel to ensure
// they are properly initialized and not closed. It's designed for health checks
// and monitoring the connection status in production environments.
//
// The method is thread-safe and can be called concurrently.
//
// Returns an error if:
//   - The connection is nil or closed
//   - The channel is nil or not initialized
func (c *clientImpl) Ping() error {
	// Lock mutex to ensure thread-safe access to connection and channel
	c.mu.Lock()
	defer c.mu.Unlock()

	// Check if the connection exists and is open
	if c.connection == nil || c.connection.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}

	// Check if a channel is initialized
	if c.channel == nil {
		return fmt.Errorf("rabbitmq channel is not initialized")
	}

	return nil
}

// establishLocalConnection creates a standard (non-TLS) AMQP connection and channel.
//
// This internal method uses the "amqp://" protocol. It is specifically
// for local development and does not enforce encryption.
//
// Returns an error if the connection or channel cannot be created.
func (c *clientImpl) establishLocalConnection() error {
	// Build the AMQP URL for a standard, non-TLS connection.
	url := fmt.Sprintf("amqp://%s:%s@%s:%s", c.params.User, c.params.Password, c.params.Host, c.params.Port)

	// Establish a standard connection to RabbitMQ, not TLS.
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Error().Err(err).Msg("Failed to dial RabbitMQ (local)")
		return err
	}

	// Store connection for later use
	c.connection = conn

	// Create an AMQP channel for operations.
	ch, errConn := conn.Channel()
	if errConn != nil {
		_ = conn.Close() // Clean up connection if channel creation fails
		log.Error().Err(errConn).Msg("Failed to create channel (local)")
		return errConn
	}

	// Store channel for later use
	c.channel = ch

	// Start a background goroutine to monitor connection health
	go c.monitorLocalConnection()

	// Channel closure monitor
	closeChan := ch.NotifyClose(make(chan *amqp.Error))
	go func() {
		if errClose := <-closeChan; err != nil {
			log.Warn().Err(errClose).Msg("RabbitMQ channel closed, reconnecting...")
			c.reconnectLocalLoop()
		}
	}()

	log.Info().Msg("RabbitMQ local connection established")
	return nil
}

// monitorLocalConnection supervises the connection (non-TLS) and attempts to reconnect if it drops.
//
// This internal method runs in a background goroutine and monitors the AMQP
// connection for closure events. When the connection is lost, it automatically
// starts the reconnection process.
//
// The method uses the AMQP connection's NotifyClose channel to detect
// connection failures and triggers the reconnection loop.
func (c *clientImpl) monitorLocalConnection() {
	// Create a channel to receive connection close notifications
	closeChan := c.connection.NotifyClose(make(chan *amqp.Error))

	// Monitor the connection for closure events
	for err := range closeChan {
		log.Warn().Err(err).Msg("RabbitMQ connection closed. Reconnecting...")
		// Start reconnection process when connection is lost
		c.reconnectLocalLoop()
	}
}

// reconnectLocalLoop tries to reconnect (non-TLS) every 5 seconds until successful.
//
// This internal method implements an exponential backoff strategy for
// reconnection attempts. It waits 5 seconds between attempts and continues
// until a successful connection is established.
//
// The method is thread-safe and uses the client's mutex to ensure
// exclusive access during reconnection attempts.
func (c *clientImpl) reconnectLocalLoop() {
	for {
		// Wait 5 seconds before attempting reconnection
		time.Sleep(5 * time.Second)

		// Lock mutex to ensure exclusive access during reconnection
		c.mu.Lock()
		err := c.establishLocalConnection()
		c.mu.Unlock()

		// If reconnection succeeds, break out of the loop
		if err == nil {
			log.Info().Msg("RabbitMQ reconnected successfully")
			break
		}

		// Log reconnection failure and continue loop
		log.Error().Err(err).Msg("Failed to reconnect to RabbitMQ")
	}
}

// closeChannel closes the AMQP channel if it exists.
//
// This internal method safely closes the AMQP channel and logs any errors
// that occur during the closure process.
//
// Returns an error if the channel cannot be closed properly.
func (c *clientImpl) closeChannel() error {
	// Check if a channel exists before attempting to close
	if c.channel != nil {
		// Close the AMQP channel
		if err := c.channel.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close channel")
			return err
		}
	}

	return nil
}

// closeConnection closes the AMQP connection if it exists.
//
// This internal method safely closes the AMQP connection and logs any errors
// that occur during the closure process.
//
// Returns an error if the connection cannot be closed properly.
func (c *clientImpl) closeConnection() error {
	// Check if the connection exists before attempting to close
	if c.connection != nil {
		// Close the AMQP connection
		if err := c.connection.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close connection")
			return err
		}
	}

	return nil
}

package broker

// Client defines the interface for the concurrent RabbitMQ client.
//
// The Client interface provides methods for connecting to RabbitMQ,
// publishing messages, consuming messages with configurable parallelism,
// and managing queues and exchanges. All methods are designed to be
// thread-safe and suitable for concurrent use.
type Client interface {
	// ConnectLocal establishes a non-TLS connection for local development.
	//
	// This method is intended for local or trusted network environments
	// where TLS is not required. It uses the standard "amqp://" protocol.
	//
	// Returns an error if the connection cannot be established.
	ConnectLocal(host, port, user, password string) error

	// Close closes the connection and all resources associated with the client.
	//
	// This method gracefully shuts down the client by:
	//   - Stopping all message consumers
	//   - Closing the AMQP channel
	//   - Closing the AMQP connection
	//   - Signaling all goroutines to stop
	//
	// The method is safe to call multiple times and will only perform
	// the shutdown once.
	//
	// Returns an error if any part of the shutdown process fails.
	Close() error

	// Ping verifies that the RabbitMQ connection is active and healthy.
	//
	// This method checks the status of the AMQP connection and channel
	// to ensure they are properly initialized and not closed. It's useful
	// for health checks and monitoring the connection status.
	//
	// Returns an error if the connection or channel is closed or not initialized.
	Ping() error
}

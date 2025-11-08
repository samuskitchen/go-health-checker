// Package enums defines application constants, such as environment variables for services like RabbitMQ.
package enums

// RabbitMQ environment variable keys and logging formats.
const (
	// RabbitHost is the environment variable for the RabbitMQ host address.
	RabbitHost string = "RABBITMQ_HOST"
	// RabbitPort is the environment variable for the RabbitMQ port.
	RabbitPort string = "RABBITMQ_PORT"
	// RabbitUser is the environment variable for the RabbitMQ username.
	RabbitUser string = "RABBITMQ_USERNAME"
	// RabbitPassword is the environment variable for the RabbitMQ password.
	RabbitPassword string = "RABBITMQ_PASSWORD"
)

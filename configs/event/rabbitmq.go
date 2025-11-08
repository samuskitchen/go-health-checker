// Package events contain messaging listeners for the audit module, including a RabbitMQListener
// that consumes and processes events.
package events

import (
	"os"
	"sync"

	"github.com/samuskitchen/go-health-checker/pkg/kit/enums"
	libRabbitmq "github.com/samuskitchen/go-health-checker/pkg/tools/broker"

	"github.com/rs/zerolog/log"
)

var (
	once         sync.Once
	rabbitClient *RabbitEvent
)

// RabbitEvent struct provides a RabbitMQ client instance
type RabbitEvent struct {
	RabbitMQClient libRabbitmq.Client
}

// RabbitConnection provides a singleton instance of RabbitEvent
func RabbitConnection() *RabbitEvent {
	once.Do(func() { getConnectionRabbit() })
	return rabbitClient
}

// NewRabbitEvent is a clean constructor for RabbitEvent, compatible with dig
func getConnectionRabbit() {
	client := libRabbitmq.NewClient()

	host := os.Getenv(enums.RabbitHost)
	port := os.Getenv(enums.RabbitPort)
	user := os.Getenv(enums.RabbitUser)
	password := os.Getenv(enums.RabbitPassword)

	validateParams(host, port, user, password)

	err := client.ConnectLocal(host, port, user, password)
	if err != nil {
		log.Error().Err(err).Msg("failed to connect to RabbitMQ")
	}

	rabbitClient = &RabbitEvent{
		RabbitMQClient: client,
	}
}

func validateParams(host string, port string, user string, password string) {
	var missingVars []string
	if host == "" {
		missingVars = append(missingVars, enums.RabbitHost)
	}
	if port == "" {
		missingVars = append(missingVars, enums.RabbitPort)
	}
	if user == "" {
		missingVars = append(missingVars, enums.RabbitUser)
	}
	if password == "" {
		missingVars = append(missingVars, enums.RabbitPassword)
	}

	if len(missingVars) > 0 {
		log.Fatal().Msgf("RabbitMQ environment variables missing: %v", missingVars)
	}
}

package router

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/samuskitchen/go-health-checker/configs/cache"
	events "github.com/samuskitchen/go-health-checker/configs/event"
	"github.com/samuskitchen/go-health-checker/configs/storage"
	"github.com/samuskitchen/go-health-checker/pkg/tools/healthcheck"
)

type healthHandler struct {
	clientPg        *storage.Data
	clientHazelcast *cache.Cache
	clientRabbit    *events.RabbitEvent
}

// HealthHandler defines the interface for the health check endpoint
type HealthHandler interface {
	HealthChecker(c echo.Context) error
}

// NewHealthHandler builds a new HealthHandler
func NewHealthHandler(clientPg *storage.Data, clientHazelcast *cache.Cache,
	clientRabbit *events.RabbitEvent,
) HealthHandler {
	return &healthHandler{
		clientPg:        clientPg,
		clientHazelcast: clientHazelcast,
		clientRabbit:    clientRabbit,
	}
}

// HealthChecker checks the health of the service
// @Description Check if service is up and healthy
// @Tags Health
// @ID finance
// @Success 200 {object} health.Response
// @Failure 404
// @Router /health [get]
func (hh *healthHandler) HealthChecker(c echo.Context) error {
	ctx := c.Request().Context()

	clients := healthcheck.Clients{
		RabbitClient:    hh.clientRabbit.RabbitMQClient,
		HazelcastClient: hh.clientHazelcast.Hazelcast,
		PgClient:        hh.clientPg.DB,
	}

	return c.JSON(http.StatusOK, clients.CheckerHealth(ctx))
}

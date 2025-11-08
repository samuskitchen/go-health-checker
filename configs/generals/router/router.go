// Package router defines the API's routing and middleware,
// configuring health checks, Swagger documentation, and status endpoints.
package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/samuskitchen/go-health-checker/pkg/kit/enums"
	kitZeroLog "github.com/samuskitchen/go-health-checker/pkg/kit/logger/zerolog"

	"github.com/labstack/echo/v4"
	middlewareEcho "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Router struct for handling routing with echo-go
type Router struct {
	server *echo.Echo
}

type healthCheckResponse struct {
	Status string `json:"status"`
}

// NewRouter constructor for routing with echo-go
func NewRouter(server *echo.Echo) *Router {
	return &Router{
		server,
	}
}

// Init configures the middleware, health check routes, Swagger, and status endpoints  on the application router.
func (r *Router) Init() {
	// Custom zerolog logger instance
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	// Middleware
	logConfig := kitZeroLog.Config{
		Logger: logger,
		FieldMap: map[string]string{
			"uri":    "@uri",
			"host":   "@host",
			"method": "@method",
			"status": "@status",
		},
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, enums.HealthPath)
		},
	}

	r.server.Use(kitZeroLog.LogWithConfig(logConfig))
	r.server.Use(middlewareEcho.Recover())
	r.server.Use(middlewareEcho.RequestID())

	apiGroup := r.server.Group(enums.BasePath)

	apiGroup.GET(enums.HealthPath, healthCheckHandler)

	for _, router := range r.server.Routes() {
		log.Info().Msgf("[%s] %s", router.Method, router.Path)
	}
}

// healthCheckHandler is a handler function that returns the health status of the server.
func healthCheckHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, healthCheckResponse{Status: "ok"})
}

// Package router defines the API's routing and middleware,
// configuring health checks, Swagger documentation, and status endpoints.
package router

import (
	"os"
	"strings"

	"github.com/samuskitchen/go-health-checker/beer/handler"
	"github.com/samuskitchen/go-health-checker/pkg/kit/enums"
	kitZeroLog "github.com/samuskitchen/go-health-checker/pkg/kit/logger/zerolog"

	// Echo es el framework web utilizado para definir rutas y handlers.
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/labstack/echo/v4"
	middlewareEcho "github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Router struct for handling routing with echo-go
type Router struct {
	server        *echo.Echo
	beerHandler   handler.BeerHandler // Handler que delega la lógica de BeerService
	healthHandler HealthHandler
}

// NewRouter constructor for routing with echo-go
func NewRouter(server *echo.Echo, beerHandler handler.BeerHandler, healthHandler HealthHandler) *Router {
	return &Router{
		server:        server,
		beerHandler:   beerHandler,
		healthHandler: healthHandler,
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

	apiGroup.GET(enums.HealthPath, r.healthHandler.HealthChecker)
	apiGroup.GET("/docs/*", echoSwagger.WrapHandler)

	// Endpoints de Beer
	apiGroup.GET("/beers", r.beerHandler.GetAllBeersHandler)

	for _, router := range r.server.Routes() {
		log.Info().Msgf("[%s] %s", router.Method, router.Path)
	}
}

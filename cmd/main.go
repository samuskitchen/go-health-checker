// Package main is the main package of the application.
// Loads environment variables, configures DI, logger, routes, and servers.
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/samuskitchen/go-health-checker/configs/generals/injector"
	"github.com/samuskitchen/go-health-checker/configs/generals/router"
	"github.com/samuskitchen/go-health-checker/configs/storage"
	"github.com/samuskitchen/go-health-checker/pkg/kit/enums"
	kitZeroLog "github.com/samuskitchen/go-health-checker/pkg/kit/logger/zerolog"
	serverEcho "github.com/samuskitchen/go-health-checker/pkg/tools/server"

	// Swagger auto-generated documentation
	_ "github.com/samuskitchen/go-health-checker/docs"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog/log"
)

// main
// @title Swagger Data the Health Checker
// @version 0.1
// @tag Health Checker
// @description This is health checker API
// @termsOfService http://swagger.io/terms/
// @host localhost:8080
// @BasePath /api-health-checker
func main() {
	// Load the dependency injection container.
	container := injector.BuildContainer()

	errEnv := godotenv.Load()
	if errEnv != nil {
		log.Warn().Msgf("Warning: No .env file found: %v", errEnv)
	}

	// Check if it starts in debugger mode.
	boolVal, errBool := strconv.ParseBool(os.Getenv("LOGGER_DEBUG"))
	if errBool != nil {
		log.Warn().Msgf("Warning: LOGGER_DEBUG must be set to true or false: %v", errBool)
	}

	// Init Logger
	debug := flag.Bool("debug", boolVal, "sets log level to debug")
	kitZeroLog.InitLogger(enums.App, *debug)

	// Configure server times
	configureServerTimes()

	err := container.Invoke(func(server *echo.Echo, route *router.Router) {
		address := fmt.Sprintf("%s:%s", os.Getenv(enums.ServerHost), os.Getenv(enums.ServerPort))
		server.Debug = os.Getenv(enums.ServerPostfix) == enums.PostfixDev
		route.Init()
		server.Logger.Fatal(server.Start(address))
	})

	if err != nil {
		panic(err)
	}

	defer func() {
		log.Info().Msg("Closing connections...")

		// Try closing database Postgres and report if there is an error
		storage.PostgresCloseConnection()

		log.Info().Msg("Resource cleanup complete.")
	}()

}

func configureServerTimes() {
	serverEcho.SetServersTimeConfiguration(serverEcho.ServersTimeConfiguration{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  30 * time.Second,
	})
}

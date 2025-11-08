// Package echo provides a thin wrapper around labstack/echo to create a server
// with CORS and configurable timeouts.
package echo

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Default timeouts for the server
const (
	defaultReadTimeout  = 5 * time.Second
	defaultWriteTimeout = 10 * time.Second
	defaultIdleTimeout  = 15 * time.Second
)

// ServersTimeConfiguration contains parameters to configure server timeouts
type ServersTimeConfiguration struct {
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// ServerConfig holds the configuration for the Echo server
type ServerConfig struct {
	acceptedHeaders []string
	acceptedHosts   []string
	timeConfig      ServersTimeConfiguration
}

// Global server configuration instance
var serverConfig = &ServerConfig{
	acceptedHeaders: []string{echo.HeaderContentType, echo.HeaderAuthorization},
	acceptedHosts:   []string{},
	timeConfig: ServersTimeConfiguration{
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
	},
}

// NewServer creates and configures an Echo server with CORS and timeouts
func NewServer() *echo.Echo {
	server := echo.New()

	// Configure server timeouts
	server.Server.ReadTimeout = serverConfig.timeConfig.ReadTimeout
	server.Server.WriteTimeout = serverConfig.timeConfig.WriteTimeout
	server.Server.IdleTimeout = serverConfig.timeConfig.IdleTimeout

	// Configure and apply CORS middleware
	configureCORS(server)

	return server
}

// configureCORS sets up CORS configuration and applies it to the server
func configureCORS(server *echo.Echo) {
	corsConfig := middleware.CORSConfig{
		AllowHeaders: serverConfig.acceptedHeaders,
	}

	if len(serverConfig.acceptedHosts) > 0 {
		corsConfig.AllowOrigins = serverConfig.acceptedHosts
	}

	server.Use(middleware.CORS())
	server.Use(middleware.CORSWithConfig(corsConfig))
}

// AddAcceptedHeader adds a new header to the list of accepted headers
func AddAcceptedHeader(header string) {
	serverConfig.acceptedHeaders = append(serverConfig.acceptedHeaders, header)
}

// AddAcceptedHost adds a new host to the list of accepted hosts
func AddAcceptedHost(host string) {
	serverConfig.acceptedHosts = append(serverConfig.acceptedHosts, host)
}

// SetServersTimeConfiguration sets timeout values for the server
func SetServersTimeConfiguration(stc ServersTimeConfiguration) {
	if stc.ReadTimeout > 0 {
		serverConfig.timeConfig.ReadTimeout = stc.ReadTimeout
	}

	if stc.WriteTimeout > 0 {
		serverConfig.timeConfig.WriteTimeout = stc.WriteTimeout
	}

	if stc.IdleTimeout > 0 {
		serverConfig.timeConfig.IdleTimeout = stc.IdleTimeout
	}
}

// GetServersTimeConfiguration returns the current time configuration
func GetServersTimeConfiguration() ServersTimeConfiguration {
	return serverConfig.timeConfig
}

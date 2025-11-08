// Package enums contains path constants and global configuration.
package enums

const (
	// BasePath is the common prefix for all API paths.
	BasePath string = "api-health-checker"

	// HealthPath is the path to the health check endpoint.
	HealthPath string = "/health"

	// ServerHost is the config key for the server hostname.
	ServerHost string = "SERVER_HOST"

	// ServerPort is the config key for the server port.
	ServerPort string = "SERVER_PORT"

	// ServerPostfix is the config key for differentiating environments (e.g., "dev").
	ServerPostfix string = "SERVER_POSTFIX"

	// DeployCountry is the config key for country.
	DeployCountry string = "DEPLOY_COUNTRY"

	// PostfixDev is the suffix used to indicate the development environment.
	PostfixDev string = "dev"

	// App is the application name used in logs and metrics.
	App string = "go-health-checker"
)

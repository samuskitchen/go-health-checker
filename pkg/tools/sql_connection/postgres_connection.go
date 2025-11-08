// Package sql_connection provides a client to interact with PostgreSQL
package sql_connection

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	tools "github.com/samuskitchen/go-health-checker/pkg/tools/models"

	// This package is used to initialize the PostgresSQL driver
	_ "github.com/lib/pq"
)

// Default pool settings
const (
	maxOpenDefaultValue   = 150
	maxIdleDefaultValue   = 64
	lifeTimeDefaultString = "15m"
	idleTimeDefaultString = "5m"
)

// Pre-parsed default durations
var (
	lifeTimeDefaultValue time.Duration
	idleTimeDefaultValue time.Duration
)

func init() {
	// Parse defaults at package initialization.
	// This will panic if the default strings are invalid, which is desirable
	// as it indicates a compile-time configuration error.
	var err error
	lifeTimeDefaultValue, err = time.ParseDuration(lifeTimeDefaultString)
	if err != nil {
		log.Fatalf("Invalid default duration 'lifeTimeDefaultString': %v", err)
	}

	idleTimeDefaultValue, err = time.ParseDuration(idleTimeDefaultString)
	if err != nil {
		log.Fatalf("Invalid default duration 'idleTimeDefaultString': %v", err)
	}
}

// sqlOpener defines the function signature for opening a database connection.
// This allows for mocking sql.Open in tests.
type sqlOpener func(driverName, dataSourceName string) (*sql.DB, error)

// Connector handles the creation of PostgreSQL database connections.
type Connector struct {
	openDB sqlOpener
}

// NewConnector creates a new instance of a Connector.
// It uses sql.Open as the default database opener.
func NewConnector() *Connector {
	return &Connector{
		openDB: sql.Open,
	}
}

// Connect validates parameters, builds a DSN, and establishes a connection
// to the PostgreSQL database, applying connection pool settings.
func (c *Connector) Connect(params tools.DbParams) (*sql.DB, error) {
	// 1. Validate required parameters
	if err := validateParams(params); err != nil {
		return nil, err
	}

	// 2. Build the Data Source Name (DSN)
	dsn := buildDSN(params)

	// 3. Open the database connection
	db, err := c.openDB("postgres", dsn)
	if err != nil {
		// Error opening the connection (e.g., driver issue)
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	// sql.Open might return a nil db without an error if the driver is unknown
	if db == nil {
		return nil, fmt.Errorf("database driver 'postgres' not found or failed to initialize")
	}

	// 4. Verify the connection is alive
	if err = db.Ping(); err != nil {
		// If ping fails, close the potentially problematic connection pool
		if closeErr := db.Close(); closeErr != nil {
			log.Printf("Failed to close database connection after ping failure: %v", closeErr)
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to the database %s", params.DbName)

	// 5. Apply connection pool settings
	applyPoolSettings(db, params)

	return db, nil
}

// validateParams checks if all required fields in DbParams are present.
func validateParams(params tools.DbParams) error {
	var missingFields []string

	if params.User == "" {
		missingFields = append(missingFields, "User")
	}
	if params.Host == "" {
		missingFields = append(missingFields, "Host")
	}
	if params.Port == "" {
		missingFields = append(missingFields, "Port")
	}
	if params.DbName == "" {
		missingFields = append(missingFields, "DbName")
	}
	if params.SslMode == "" {
		missingFields = append(missingFields, "SslMode")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required fields: %v", missingFields)
	}
	return nil
}

// buildDSN constructs the connection string for PostgreSQL.
func buildDSN(params tools.DbParams) string {
	// Use url.QueryEscape for the password to handle special characters.
	return fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		params.User,
		url.QueryEscape(params.Password),
		params.Host,
		params.Port,
		params.DbName,
		params.SslMode,
	)
}

// applyPoolSettings configures database/sql's built-in pool.
func applyPoolSettings(db *sql.DB, params tools.DbParams) {
	maxOpen := parseIntWithDefault(params.MaxOpenCon, maxOpenDefaultValue)
	maxIdle := parseIntWithDefault(params.MaxIdleCon, maxIdleDefaultValue)
	lifetime := parseDurationWithDefault(params.MaxLifeTimeCon, lifeTimeDefaultValue)
	idleTime := parseDurationWithDefault(params.MaxIdleTimeCon, idleTimeDefaultValue)

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(lifetime)
	db.SetConnMaxIdleTime(idleTime)
}

// parseIntWithDefault parses a string to an int, returning a default value
// if the string is empty or invalid.
func parseIntWithDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	n, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Warning: could not parse int value '%s', using default %d. Error: %v", value, defaultValue, err)
		return defaultValue
	}
	return n
}

// parseDurationWithDefault parses a string to a time.Duration, returning a
// default value if the string is empty or invalid.
func parseDurationWithDefault(value string, defaultValue time.Duration) time.Duration {
	if value == "" {
		return defaultValue
	}
	d, err := time.ParseDuration(value)
	if err != nil {
		log.Printf("Warning: could not parse duration value '%s', using default %v. Error: %v", value, defaultValue, err)
		return defaultValue
	}
	return d
}

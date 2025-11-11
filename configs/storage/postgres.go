// Package storage provides the singleton connection to PostgreSQL,
// offering functions to get and close the database connection.
package storage

import (
	"database/sql"
	"os"
	"sync"

	"github.com/samuskitchen/go-health-checker/pkg/kit/enums"
	modelConnection "github.com/samuskitchen/go-health-checker/pkg/tools/models"
	"github.com/samuskitchen/go-health-checker/pkg/tools/sqlconnection"

	"github.com/rs/zerolog/log"
)

var (
	once sync.Once
	data *Data
)

// Data contains the configuration needed to connect to the Postgres database.
type Data struct {
	DB *sql.DB
}

// PostgresConnection returns the singleton instance of the connection to PostgreSQL.
// Initializes the connection the first time it is invoked.
func PostgresConnection() *Data {
	once.Do(getConnections)
	return data
}

func getConnections() {
	dbParams := modelConnection.DbParams{
		Host:           os.Getenv(enums.HostConnection),
		Port:           os.Getenv(enums.PortConnection),
		User:           os.Getenv(enums.UserConnection),
		Password:       os.Getenv(enums.PasswordConnection),
		DbName:         os.Getenv(enums.PostgresDatabase),
		SslMode:        os.Getenv(enums.SslModeConnection),
		MaxOpenCon:     os.Getenv(enums.PostgresMaxOpenCons),
		MaxIdleCon:     os.Getenv(enums.PostgresMaxIdleCons),
		MaxLifeTimeCon: os.Getenv(enums.PostgresMaxLifetime),
		MaxIdleTimeCon: os.Getenv(enums.PostgresMaxIdleTime),
	}

	conn, err := sqlconnection.NewConnector().Connect(dbParams)
	if err != nil {
		log.Error().Msgf("error connecting to database: %v", err)
	}

	// Optional but recommended: verify we can actually talk to the DB now.
	if errConn := conn.Ping(); errConn != nil {
		log.Error().Msgf("database ping failed: %v", err)
	}

	data = &Data{
		DB: conn,
	}
}

// PostgresCloseConnection closes the PostgreSQL singleton connection if it has been initialized.
// Logs fatal on error closing.
func PostgresCloseConnection() {
	if data != nil {
		if err := data.DB.Close(); err != nil {
			log.Fatal().Msgf("Error closing the database: %v", err)
		}
	}
}

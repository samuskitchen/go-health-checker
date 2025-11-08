// Package enums contains configuration constants for connecting to PostgreSQL.
package enums

const (
	// HostConnection it is the configuration key for the PostgreSQL host.
	HostConnection string = "DB_HOST"
	// PortConnection is the configuration key for the PostgreSQL port.
	PortConnection string = "DB_PORT"
	// UserConnection is the configuration key for the PostgreSQL user.
	UserConnection string = "DB_USER"
	// PasswordConnection is the configuration key for the PostgreSQL password.
	PasswordConnection string = "DB_PASSWORD"
	// SslModeConnection Connection is the configuration key for PostgreSQL SSL mode.
	SslModeConnection string = "DB_SSL_MODE"
	// PostgresDatabase is the configuration key for the database name in PostgreSQL.
	PostgresDatabase string = "DB_NAME"
	// PostgresMaxOpenCons is the configuration key for the database open connections.
	PostgresMaxOpenCons string = "DB_MAX_OPEN"
	// PostgresMaxIdleCons is the configuration key for the maximun connections idle at the same time.
	PostgresMaxIdleCons string = "DB_MAX_IDLE"
	// PostgresMaxLifetime is the configuration key for the connection max timeline.
	PostgresMaxLifetime string = "DB_CONN_MAX_LIFETIME"
	// PostgresMaxIdleTime is the configuration key for the database to config max idle time.
	PostgresMaxIdleTime string = "DB_CONN_MAX_IDLE_TIME"
)

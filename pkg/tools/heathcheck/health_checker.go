// Package heathcheck provides a health checker for various clients
package heathcheck

import (
	"context"
	"database/sql"
	"time"

	"github.com/hellofresh/health-go/v5"
	"github.com/samuskitchen/go-health-checker/pkg/tools/broker"
	"github.com/samuskitchen/go-health-checker/pkg/tools/datastore"
)

// Clients represent the clients to be checked
type Clients struct {
	RabbitClient    broker.Client
	HazelcastClient datastore.IClient
	PgClient        *sql.DB
}

// Response represents the health check response
type Response struct {
	OverallStatus string   `json:"overallStatus"`
	Timestamp     string   `json:"timestamp"`
	Checks        []Health `json:"checks"`
}

// Health represents the health check response
type Health struct {
	Status    string `json:"status"`
	Component string `json:"component"`
	Version   string `json:"version"`
}

// CheckerHealth performs a health check on all clients
func (cl *Clients) CheckerHealth(ctx context.Context) Response {
	var checks []Health

	// Collect Health Checks from all clients
	cl.collectRabbitMQChecks(ctx, &checks)
	cl.collectHazelcastChecks(ctx, &checks)
	cl.collectPostgresSQLChecks(ctx, &checks)

	// Calculate Overall Status based on the number of OK checks
	overallStatus := calculateOverallStatus(checks)

	return Response{
		OverallStatus: overallStatus,
		Timestamp:     time.Now().Format(time.RFC3339),
		Checks:        checks,
	}
}

// collectRabbitMQChecks collects health checks for RabbitMQ client
func (cl *Clients) collectRabbitMQChecks(ctx context.Context, checks *[]Health) {
	if check := cl.checkRabbitMQ(ctx); check != nil {
		*checks = append(*checks, *check)
	}
}

// collectHazelcastChecks collects health checks for a Hazelcast client
func (cl *Clients) collectHazelcastChecks(ctx context.Context, checks *[]Health) {
	if check := cl.checkHazelcast(ctx); check != nil {
		*checks = append(*checks, *check)
	}
}

// collectPostgresSQLChecks collects health checks for PostgresSQL database/sql client
func (cl *Clients) collectPostgresSQLChecks(ctx context.Context, checks *[]Health) {
	if check := cl.checkPostgresSQL(ctx); check != nil {
		*checks = append(*checks, *check)
	}
}

// checkRabbitMQ performs a health check on RabbitMQ
func (cl *Clients) checkRabbitMQ(ctx context.Context) *Health {
	if cl.RabbitClient == nil {
		return nil
	}

	h, _ := health.New(
		health.WithComponent(health.Component{Name: "RabbitMQ", Version: "1.0.0"}),
		health.WithChecks(health.Config{
			Name:      "rabbitmq-connection",
			Timeout:   time.Second * 5,
			SkipOnErr: true,
			Check: func(_ context.Context) error {
				return cl.RabbitClient.Ping()
			},
		}),
	)

	data := h.Measure(ctx)
	return &Health{
		Status:    string(data.Status),
		Component: data.Name,
		Version:   data.Component.Version,
	}
}

// checkHazelcast performs a health check on a Hazelcast client
func (cl *Clients) checkHazelcast(ctx context.Context) *Health {
	if cl.HazelcastClient == nil {
		return nil
	}

	h, _ := health.New(
		health.WithComponent(health.Component{Name: "Hazelcast", Version: "1.0.0"}),
		health.WithChecks(health.Config{
			Name:      "hazelcast-connection",
			Timeout:   time.Second * 5,
			SkipOnErr: true,
			Check: func(_ context.Context) error {
				return cl.HazelcastClient.Ping()
			},
		}),
	)

	data := h.Measure(ctx)
	return &Health{
		Status:    string(data.Status),
		Component: data.Name,
		Version:   data.Component.Version,
	}
}

// checkPostgresSQL performs a health check for PostgresSQL database/sql client
func (cl *Clients) checkPostgresSQL(ctx context.Context) *Health {
	if cl.PgClient == nil {
		return nil
	}

	h, _ := health.New(
		health.WithComponent(health.Component{Name: "postgresql-sql", Version: "1.0.0"}),
		health.WithChecks(health.Config{
			Name:      "postgresql-sql-connection",
			Timeout:   time.Second * 5,
			SkipOnErr: true,
			Check: func(ctx context.Context) error {
				return cl.PgClient.PingContext(ctx)
			},
		}),
	)

	data := h.Measure(ctx)
	return &Health{
		Status:    string(data.Status),
		Component: data.Name,
		Version:   data.Component.Version,
	}
}

// calculateOverallStatus calculates the overall status of the checks based on the number of OK checks
func calculateOverallStatus(checks []Health) string {
	if len(checks) == 0 {
		return "unknown"
	}

	okCount := 0
	totalCount := len(checks)

	for _, check := range checks {
		if check.Status == "OK" {
			okCount++
		}
	}

	// All checks are OK
	if okCount == totalCount {
		return "Available"
	}

	// Not all checks are OK
	if okCount == 0 {
		return "Unavailable"
	}

	// Some checks are OK, some are not
	return "Partially Available"
}
